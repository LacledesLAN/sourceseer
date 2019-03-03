package csgo

import (
	"strconv"
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

//Scenario represents a set of behaviors and rules to add to a CSGO server.
type Scenario func(*CSGO) *CSGO

//ClinchableMapCycle takes an odd number of maps, uses them as the CSGO server's map cycle, and ends the server when a team wins 1/2 * map count + 1 maps.
func ClinchableMapCycle(mapCycle []string) Scenario {
	if l := len(mapCycle); l == 0 || l%2 == 0 {
		panic("A positive, odd-number of maps must be provided!")
	}

	if err := validateStockMapNames(mapCycle); err != nil {
		panic(err)
	}

	return func(g *CSGO) *CSGO {
		g.addCvarWatch("mp_maxrounds", "mp_match_restart_delay", "mp_overtime_maxrounds", "sv_pausable")
		g.AddLaunchArg("+map " + mapCycle[0] + "")

		matchHistory := make([]string, len(mapCycle)+1)

		gameSay := func(g *CSGO, msg string) {
			g.cmdIn <- "say " + msg
			g.cmdIn <- "sm_csay " + msg
		}

		setOver := func(g *CSGO, matchHistory []string) {
			for g.cmdIn != nil {
				for _, m := range matchHistory {
					gameSay(g, m)
				}

				gameSay(g, "GAME OVER - TEAM CAPTAINS REPORT TO TOURNEY ADMIN")
				time.Sleep(12 * time.Second)
			}
		}

		g.AddLogProcessor(func(le srcds.LogEntry) (keepProcessing bool) {
			if strings.HasPrefix(le.Message, `World triggered "Round_End"`) {
				mpMaxrounds, _ := g.GetCvarAsInt("mp_maxrounds")
				mpOvertimeMaxrounds, _ := g.GetCvarAsInt("mp_overtime_maxrounds")
				matchWinThreshold := calculateWinThreshold(mpMaxrounds, mpOvertimeMaxrounds, g.currentMap.roundsCompleted)

				if g.currentMap.roundsCompleted >= matchWinThreshold {
					if g.currentMap.mpTeam1.roundsWon >= matchWinThreshold {
						msg := g.currentMap.mpTeam1.name + ` won "` + g.currentMap.name + `" ` + strconv.Itoa(g.currentMap.mpTeam1.roundsWon) + "-" + strconv.Itoa(g.currentMap.mpTeam2.roundsWon)
						matchHistory = append(matchHistory, msg)
						gameSay(g, "|--------------------------|")
						gameSay(g, msg)
						gameSay(g, "|--------------------------|")
					} else if g.currentMap.mpTeam2.roundsWon >= matchWinThreshold {
						msg := g.currentMap.mpTeam2.name + ` won "` + g.currentMap.name + `" ` + strconv.Itoa(g.currentMap.mpTeam2.roundsWon) + "-" + strconv.Itoa(g.currentMap.mpTeam1.roundsWon)
						matchHistory = append(matchHistory, msg)
						gameSay(g, "|--------------------------|")
						gameSay(g, msg)
						gameSay(g, "|--------------------------|")
					} else {
						return true
					}

					if setWinThreshold := (len(mapCycle) / 2) + 1; len(g.maps) >= setWinThreshold {
						var mpTeam1MatchWins, mpTeam2MatchWins int

						for _, m := range g.maps {
							// TODO - need a better tracking mechanism
							if m.mpTeam1.roundsWon > m.mpTeam2.roundsWon && m.mpTeam1.name == g.defaultMpTeamname1 {
								mpTeam1MatchWins = mpTeam1MatchWins + 1
							} else {
								mpTeam2MatchWins = mpTeam2MatchWins + 1
							}
						}

						winningTeamName := ""
						if mpTeam1MatchWins >= setWinThreshold {
							winningTeamName = g.currentMap.mpTeam1.name
						} else if mpTeam2MatchWins >= setWinThreshold {
							winningTeamName = g.currentMap.mpTeam2.name
						}

						if len(winningTeamName) > 0 {
							msg := winningTeamName + " wins the set!"
							matchHistory = append(matchHistory, msg)

							if svPausable, err := g.GetCvarAsInt("sv_pausable"); err == nil && svPausable == 1 {
								go setOver(g, matchHistory)
								time.Sleep(7 * time.Second)
								g.cmdIn <- "pause"
							} else {
								mpMatchRestartDelay, err := g.GetCvarAsInt("mp_match_restart_delay")
								if err != nil {
									mpMatchRestartDelay = defaultMpMatchRestartDelay
								}

								go setOver(g, matchHistory)

								time.Sleep(time.Duration(mpMatchRestartDelay-2) * time.Second)
								g.cmdIn <- "mp_warmup_pausetimer 1"
								g.cmdIn <- "mp_warmup_start"
							}

							return false
						}
					}

					mpMatchRestartDelay, err := g.GetCvarAsInt("mp_match_restart_delay")
					if err != nil {
						mpMatchRestartDelay = defaultMpMatchRestartDelay
					}

					go func(g *CSGO, nextMap string) {
						g.cmdIn <- "say NEXT MAP: " + nextMap
						g.cmdIn <- "sm_csay NEXT MAP: " + nextMap

						mpMatchRestartDelay = mpMatchRestartDelay - 2
						if mpMatchRestartDelay < 0 {
							mpMatchRestartDelay = 0
						}

						time.Sleep(time.Duration(mpMatchRestartDelay) * time.Second)
						g.cmdIn <- "changelevel " + nextMap
					}(g, mapCycle[len(g.maps)])
				}
			}

			return true
		})

		return g
	}
}

//UseTeamNames sets up CSGO to use specified team names
func UseTeamNames(mpTeamname1, mpTeamname2 string) Scenario {
	var args []string

	// Process mpTeamname1
	mpTeamname1 = SanitizeTeamName(mpTeamname1)
	if len(mpTeamname1) > 0 {
		args = append(args, "+mp_teamname_1", `"`+mpTeamname1+`"`)
	}

	// Process mpTeamname2
	mpTeamname2 = SanitizeTeamName(mpTeamname2)
	if len(mpTeamname2) > 0 {
		args = append(args, "+mp_teamname_2", `"`+mpTeamname2+`"`)
	}

	if strings.ToLower(mpTeamname1) == strings.ToLower(mpTeamname2) {
		panic("team names cannot be the same")
	}

	args = append(args, `+hostname "`+HostnameFromTeamNames(mpTeamname1, mpTeamname2)+`"`)

	return func(g *CSGO) *CSGO {
		g.AddLaunchArg(args...)
		g.defaultMpTeamname1 = mpTeamname1
		g.defaultMpTeamname2 = mpTeamname2

		return g
	}
}
