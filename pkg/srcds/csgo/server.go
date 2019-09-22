package csgo

import (
	"github.com/lacledeslan/sourceseer/pkg/srcds"
)

// Server represents an interactive CSGO server
type Server interface {
	Observer
}

// Start the CSGO server
func (s *server) Start() {

}

//// WarMod Hacks ¯\_ಠ_ಠ_/¯
//if strings.HasPrefix(le.Message, "[WarMod_BFG]") {
//	// WarMod drops teamnames during the LO3 before knife fights
//	if strings.Contains(le.Message, `", "event": "knife_round_`) {
//		if len(o.game.mpTeamname1) > 0 {
//			g.cmdIn <- "mp_teamname_1 " + g.teamAssignedToCT
//		}
//
//		if len(o.game.mpTeamname1) > 0 {
//			g.cmdIn <- "mp_teamname_2 " + g.teamAssignedToTerrorist
//		}
//	}
//}

//Wrapper for observing and interacting with a SRCDS instance
//func Wrapper(osArgs ...string) Reactor {
//	c := new()
//
//	s := srcds.Wrapper(osArgs...)
//
//	c.s = s
//	c.t = s
//
//	return nil
//}

type server struct {
	srcds.Server
	observer
}

func newServer() *server {
	return &server{}
}

//func (s *server) gameLoop(srcdsMsg <-chan srcds.LogEntry) {
//	for {
//		for le := range srcdsMsg {
//			if le.Message == `World triggered "Game_Commencing"` {
//				fmt.Println("\t==================================================================================================")
//				fmt.Println("\tGame Commencing")
//				fmt.Println("\t==================================================================================================")
//				break
//			}
//
//			if le.Message == `World triggered "Match_Start"` {
//				// First message received after all config files have been processed
//				s.RefreshCvars()
//			}
//
//			if g.currentMap == nil {
//				// This will be logged for the very first map; needed to seed initial map state before players connect
//				if mapName, ok := parseLoadingMap(le); ok {
//					g.mapChanged(mapName)
//				}
//			}
//		}
//
//		for {
//			for le := range srcdsMsg {
//				if le.Message == "Starting Freeze period" {
//					fmt.Println("\t==================================================================================================")
//					fmt.Println("\tFREEZE PERIOD")
//					fmt.Println("\t==================================================================================================")
//					break
//				}
//			}
//
//		FreezePeriod:
//
//			for le := range srcdsMsg {
//				if le.Message == `World triggered "Round_Start"` {
//					if mapName, ok := parseMatchStart(le); ok {
//						g.mapChanged(mapName)
//						// TODO RESET all player and round stats
//						fmt.Println("\t==================================================================================================")
//						fmt.Println("\tMATCH START")
//						fmt.Println("\t==================================================================================================")
//						break
//					}
//				}
//
//				if le.Message == "Starting Freeze period" {
//					goto FreezePeriod
//				}
//			}
//
//			fmt.Println("\t==================================================================================================")
//			fmt.Println("\tROUND START")
//			fmt.Println("\t==================================================================================================")
//
//			for le := range srcdsMsg {
//				if le.Message == "Starting Freeze period" {
//					goto FreezePeriod
//				}
//
//				if _, ok := parseTeamCTScored(le); ok {
//					// process ct score update
//					break
//				}
//			}
//			fmt.Println("\t==================================================================================================")
//			fmt.Println("\tCT SCORED")
//			fmt.Println("\t==================================================================================================")
//
//			if _, ok := parseTeamTScored(<-srcdsMsg); ok {
//				// process t score update
//				break
//			}
//			fmt.Println("\t==================================================================================================")
//			fmt.Println("\tT SCORED")
//			fmt.Println("\t==================================================================================================")
//
//			for le := range srcdsMsg {
//				if le.Message == `World triggered "Round_End"` {
//					//check for end game
//					break
//				}
//			}
//			fmt.Println("\t==================================================================================================")
//			fmt.Println("\tROUND END")
//			fmt.Println("\t==================================================================================================")
//		}
//	}
//}
