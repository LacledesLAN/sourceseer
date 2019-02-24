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

	//mapCycle := clinchableMapCycle{
	//	mapsPending:  maps,
	//	mapsFinished: make([]string, len(maps)),
	//}

	return func(g *CSGO) *CSGO {
		g.AddCvarWatch("mp_maxrounds", "mp_overtime_maxrounds", "mp_match_restart_delay", "sv_pausable")
		g.AddLaunchArg("+map " + maps[0])

		g.AddLogProcessor(func(le srcds.LogEntry) (keepProcessing bool) {
			if strings.HasPrefix(le.Message, `World triggered "Round_End"`) {

				// determine win threshold
				mpMaxrounds, _ := g.GetCvarAsInt("mp_maxrounds")
				mpOvertimeMaxrounds, _ := g.GetCvarAsInt("mp_overtime_maxrounds")

				fmt.Println("mp_maxrounds", mpMaxrounds)
				fmt.Println("mp_overtime_maxrounds", mpOvertimeMaxrounds)

				winThreshold := calculateWinThreshold(mpMaxrounds, mpOvertimeMaxrounds, g.currentMap.RoundsCompleted())

				if g.currentMap.RoundsCompleted() >= winThreshold {

				}

				//g.currentMap.RoundsCompleted

				// determine if any team won

				mapOver := false

				if mapOver /* && len(mapCycle.mapsPending) == 0 */ {
					if value, found := g.GetCvar("sv_pausable"); found && value == "1" {
						g.cmdIn <- "pause"
					} else {
						g.cmdIn <- "mp_warmup_start"
					}

					go func() {
						for {
							g.cmdIn <- "say GAME OVER; TEAM CAPTAINS REPORT TO TOURNEY ADMIN"
							time.Sleep(5 * time.Second)
						}
					}()
				}
			} else {
				//fmt.Println("?? ", le.Message)
			}

			return true
		})

		return g
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

	return func(g *CSGO) *CSGO {
		g.AddLaunchArg(args...)

		return g
	}
}
