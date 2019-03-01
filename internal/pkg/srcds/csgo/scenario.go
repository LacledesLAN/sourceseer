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
		g.AddCvarWatch("mp_maxrounds", "mp_match_restart_delay", "mp_overtime_maxrounds", "sv_pausable")
		g.AddLaunchArg("+map " + mapCycle[0])

		var mpTeam1MatchWins, mpTeam2MatchWins int

		g.AddLogProcessor(func(le srcds.LogEntry) (keepProcessing bool) {
			if strings.HasPrefix(le.Message, `World triggered "Round_End"`) {
				mpMaxrounds, _ := g.GetCvarAsInt("mp_maxrounds")
				mpOvertimeMaxrounds, _ := g.GetCvarAsInt("mp_overtime_maxrounds")
				matchWinThreshold := calculateWinThreshold(mpMaxrounds, mpOvertimeMaxrounds, g.currentMap.RoundsCompleted())

				if g.currentMap.RoundsCompleted() >= matchWinThreshold {
					if g.currentMap.mpTeam1.roundsWon >= matchWinThreshold {
						mpTeam1MatchWins = mpTeam1MatchWins + 1
						msg := g.currentMap.mpTeam1.name + " wins the match with " + strconv.Itoa(g.currentMap.mpTeam1.roundsWon)
						g.cmdIn <- "say " + msg
						g.cmdIn <- "sm_csay " + msg
					} else if g.currentMap.mpTeam2.roundsWon >= matchWinThreshold {
						mpTeam2MatchWins = mpTeam2MatchWins + 1
						msg := g.currentMap.mpTeam2.name + " wins the match with " + strconv.Itoa(g.currentMap.mpTeam2.roundsWon)
						g.cmdIn <- "say " + msg
						g.cmdIn <- "sm_csay " + msg
					} else {
						return true
					}

					if setWinThreshold := (len(mapCycle) / 2) + 1; len(g.maps) >= setWinThreshold {
						var winningTeamName string
						if mpTeam1MatchWins >= setWinThreshold {
							winningTeamName = g.currentMap.mpTeam1.name
						} else if mpTeam2MatchWins >= setWinThreshold {
							winningTeamName = g.currentMap.mpTeam2.name
						}

						if len(winningTeamName) > 0 {
							g.cmdIn <- "say " + winningTeamName + " wins the set!"
							g.cmdIn <- "sm_csay " + winningTeamName + " wins the set!"

							if svPausable, err := g.GetCvarAsInt("sv_pausable"); err == nil && svPausable == 1 {
								time.Sleep(6 * time.Second)
								g.cmdIn <- "pause"
								g.cmdIn <- "say GAME OVER - TEAM CAPTAINS REPORT TO TOURNEY ADMIN"
								g.cmdIn <- "sm_csay GAME OVER - TEAM CAPTAINS REPORT TO TOURNEY ADMIN"
							} else {
								mpMatchRestartDelay, err := g.GetCvarAsInt("mp_match_restart_delay")
								if err != nil {
									mpMatchRestartDelay = defaultMpMatchRestartDelay
								}
								go func() {
									for {
										g.cmdIn <- "say GAME OVER - TEAM CAPTAINS REPORT TO TOURNEY ADMIN"
										g.cmdIn <- "sm_csay GAME OVER - TEAM CAPTAINS REPORT TO TOURNEY ADMIN"
										time.Sleep(8 * time.Second)
									}
								}()

								go func() {
									time.Sleep(time.Duration(mpMatchRestartDelay-2) * time.Second)
									g.cmdIn <- "mp_warmup_pausetimer 1"
									g.cmdIn <- "mp_warmup_start"
								}()
							}

							return false
						}
					}

					mpMatchRestartDelay, err := g.GetCvarAsInt("mp_match_restart_delay")
					if err != nil {
						mpMatchRestartDelay = defaultMpMatchRestartDelay
					}

					go func(g *CSGO, nextLevel string) {
						g.cmdIn <- "say NEXT LEVEL: " + nextLevel
						g.cmdIn <- "sm_csay NEXT LEVEL: " + nextLevel

						mpMatchRestartDelay = mpMatchRestartDelay - 2
						if mpMatchRestartDelay < 0 {
							mpMatchRestartDelay = 0
						}

						time.Sleep(time.Duration(mpMatchRestartDelay) * time.Second)
						g.cmdIn <- "changelevel " + nextLevel
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
		args = append(args, "+mp_teamname_1", mpTeamname1)
	}

	// Process mpTeamname2
	mpTeamname2 = SanitizeTeamName(mpTeamname2)
	if len(mpTeamname2) > 0 {
		args = append(args, "+mp_teamname_2", mpTeamname2)
	}

	args = append(args, `+hostname "`+HostnameFromTeamNames(mpTeamname1, mpTeamname2)+`"`)

	return func(g *CSGO) *CSGO {
		g.AddLaunchArg(args...)

		return g
	}
}
