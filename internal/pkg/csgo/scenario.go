package csgo

import (
	"fmt"
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

//Scenario represents a set of behaviors and rules to add to a CSGO server.
type Scenario func(*CSGO) *CSGO

type clinchableMapCycle struct {
	mapsPending  []string
	mapsFinished []string
}

//ClinchableMapCycle takes an odd number of maps, uses them as the CSGO server's map cycle, and ends the server when a team wins 1/2 * map count + 1 maps.
func ClinchableMapCycle(maps []string) Scenario {
	if l := len(maps); l == 0 || l%2 == 0 {
		panic("A positive, odd-number of maps must be provided!")
	}

	if err := validateStockMapNames(maps); err != nil {
		panic(err)
	}

	mapCycle := clinchableMapCycle{
		mapsPending:  maps,
		mapsFinished: make([]string, len(maps)),
	}

	return func(gs *CSGO) *CSGO {
		gs.srcds.AddCvarWatch("mp_match_restart_delay", "sv_pausable")
		gs.srcds.AddLaunchArg("+map " + maps[0])
		gs.srcds.AddLogProcessor(func(le srcds.LogEntry) (keepProcessing bool) {
			if strings.HasPrefix(le.Message, `World triggered "Round_End"`) {

				mapOver := false

				if gs.isOvertime() {
					// TODO account for mp_overtime_maxrounds even values
					//// OVERTIME WIN CONDITION?
					//
					//delta = mpMaxRounds	//30
					//while totalRoundsPlayed - delta > mpOvertimeMaxRounds {
					//	delta = delta + mpOvertimeMaxRounds
					//}
					//
					//if ct.RoundsWon - delta >= (mpOvertimeMaxRounds/2) + 1 {
					//	// ct won
					//} else if t.RoundsWon >= (mpOvertimeMaxRounds/2) + 1 {
					//	// t won
					//}
				} else {
					if gs.maxRounds() == 1 {
						if gs.currentMap.ct().roundsWon > 0 {
							mapOver = true //ct won
						} else if gs.currentMap.terrorist().roundsWon > 0 {
							mapOver = true //t won
						}
					} else {
						if gs.currentMap.ct().roundsWon > (gs.maxRounds()/2)+1 {
							mapOver = true //ct won
						} else if gs.currentMap.terrorist().roundsWon > (gs.maxRounds()/2)+1 {
							mapOver = true //t won
						}
					}
				}

				if mapOver /* && len(mapCycle.mapsPending) == 0 */ {
					if value, found := gs.srcds.GetCvar("sv_pausable"); found && value == "1" {
						gs.cmdIn <- "pause"
					} else {
						gs.cmdIn <- "mp_warmup_start"
					}

					go func() {
						for {
							gs.cmdIn <- "say GAME OVER; TEAM CAPTAINS REPORT TO TOURNEY ADMIN"
							time.Sleep(5 * time.Second)
						}
					}()
				}
			} else {
				fmt.Println(">>>>", le.Message)
			}

			return true
		})

		return gs
	}
}

//MapPreliminaries executes ready up, knife mode, side selection, and live on three.
func MapPreliminaries(mpTeamname1, mpTeamname2 string) Scenario {
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

	return func(gs *CSGO) *CSGO {
		gs.srcds.AddLaunchArg(args...)

		return gs
	}
}
