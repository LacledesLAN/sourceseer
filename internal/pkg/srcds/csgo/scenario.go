package csgo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

//Scenario represents a set of behaviors and rules to add to a CSGO server.
type Scenario func(*CSGO) *CSGO

func buildSetWonMessage(winningTeamName string, matchesWon, matchesLost int) string {
	return winningTeamName + " won the set (" + strconv.Itoa(matchesWon) + "-" + strconv.Itoa(matchesLost) + ")!"
}

func buildMatchWonMessage(winningTeam teamState, mapName string, matchNumber int) string {
	return winningTeam.name + " won match " + strconv.Itoa(matchNumber) + ` on "` + mapName + `" ` + strconv.Itoa(winningTeam.roundsWon) + "-" + strconv.Itoa(winningTeam.roundsLost)
}

//ClinchableMapCycle takes an odd number of maps, uses them as the CSGO server's map cycle, and ends the server when a team wins 1/2 * map count + 1 maps.
func ClinchableMapCycle(mapCycle []string) Scenario {
	if l := len(mapCycle); l == 0 || l%2 == 0 {
		panic("A positive, odd-number of maps must be provided!")
	}

	fmt.Printf("[SOURCESEER] Will be using clinchable map cycle: %v\n", mapCycle)

	return func(g *CSGO) *CSGO {
		g.addCvarWatch("mp_maxrounds", "mp_match_restart_delay", "mp_overtime_maxrounds", "sv_pausable")
		g.AddLaunchArg("+map " + mapCycle[0] + "")

		matchHistory := []string{}

		setOver := func(g *CSGO, matchHistory []string) {
			for g.cmdIn != nil {
				g.say("GAME OVER - TEAM CAPTAINS REPORT TO TOURNEY ADMIN", false)
				time.Sleep(12 * time.Second)
				for _, m := range matchHistory {
					g.say(m, false)
				}
			}
		}

		g.AddLogProcessor(func(le srcds.LogEntry) (keepProcessing bool) {
			if strings.HasPrefix(le.Message, `World triggered "Round_End"`) {
				mpMaxrounds, _ := g.GetCvarAsInt("mp_maxrounds")
				mpOvertimeMaxrounds, _ := g.GetCvarAsInt("mp_overtime_maxrounds")
				matchWinThreshold := calculateWinThreshold(mpMaxrounds, mpOvertimeMaxrounds, g.currentMap.roundsCompleted)

				if g.currentMap.roundsCompleted >= matchWinThreshold {
					if g.currentMap.mpTeam1.roundsWon >= matchWinThreshold {
						msg := buildMatchWonMessage(g.currentMap.mpTeam1, g.currentMap.name, len(matchHistory)+1)
						matchHistory = append(matchHistory, msg)
					} else if g.currentMap.mpTeam2.roundsWon >= matchWinThreshold {
						msg := buildMatchWonMessage(g.currentMap.mpTeam2, g.currentMap.name, len(matchHistory)+1)
						matchHistory = append(matchHistory, msg)
					} else {
						return true
					}

					g.say("|--------------------------|", false)
					for i, m := range matchHistory {
						g.say(m, i-1 == len(matchHistory))
					}
					g.say("|--------------------------|", false)

					fmt.Printf("[SOURCESEER] Match %v on map %q has ended.\n", len(matchHistory), g.currentMap.name)
					fmt.Printf("[SOURCESEER] mp_team_1 %q - won %v rounds and lost %v rounds.\n", g.currentMap.mpTeam1.name, g.currentMap.mpTeam1.roundsWon, g.currentMap.mpTeam1.roundsLost)
					fmt.Printf("[SOURCESEER] mp_team_2 %q - won %v rounds and lost %v rounds.\n", g.currentMap.mpTeam2.name, g.currentMap.mpTeam2.roundsWon, g.currentMap.mpTeam2.roundsLost)

					if setWinThreshold := (len(mapCycle) / 2) + 1; len(matchHistory) >= setWinThreshold {
						var teamAssignedToCTwins, teamAssignedToTerroristWins int

						for _, m := range g.maps {
							if m.mpTeam1.roundsWon > m.mpTeam2.roundsWon {
								if m.mpTeam1.name == g.teamAssignedToCT {
									teamAssignedToCTwins = teamAssignedToCTwins + 1
								} else {
									teamAssignedToTerroristWins = teamAssignedToTerroristWins + 1
								}
							} else {
								if m.mpTeam2.name == g.teamAssignedToCT {
									teamAssignedToCTwins = teamAssignedToCTwins + 1
								} else {
									teamAssignedToTerroristWins = teamAssignedToTerroristWins + 1
								}
							}
						}

						setWonMsg := ""
						if teamAssignedToCTwins >= setWinThreshold {
							setWonMsg = buildSetWonMessage(g.teamAssignedToCT, teamAssignedToCTwins, teamAssignedToTerroristWins)
						} else if teamAssignedToTerroristWins >= setWinThreshold {
							setWonMsg = buildSetWonMessage(g.teamAssignedToTerrorist, teamAssignedToTerroristWins, teamAssignedToCTwins)
						}

						if len(setWonMsg) > 0 {
							matchHistory = append(matchHistory, setWonMsg)

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
						g.say("NEXT MAP: "+nextMap, true)

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

	hostname := HostnameFromTeamNames(mpTeamname1, mpTeamname2)
	args = append(args, `+hostname "`+hostname+`"`)
	args = append(args, `+tv_name zGO-TV-"`+hostname+`"`)

	return func(g *CSGO) *CSGO {
		g.AddLaunchArg(args...)
		g.teamAssignedToCT = mpTeamname1
		g.teamAssignedToTerrorist = mpTeamname2

		return g
	}
}
