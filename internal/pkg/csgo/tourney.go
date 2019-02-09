package csgo

import (
	"time"
)

const (
	// Max amount of time for players to readyup for the very first map (seconds)
	maxReadyupConnectTime = time.Minute * 7
	// Max amount of seconds for players to readyup for knife rounds 2+
	maxReadyUpKnifeTime = time.Minute * 3
	// If all players not ready by this time period auto start (seconds)
	maxReadyUpPlayTime = time.Minute * 3
)

var (
	csgoStdin chan string
)

func StartTourney(mpTeamname1 string, mpTeamname2 string, maps []string) {

	//csgoState := NewGameState()
	//
	//go func(g *gameState) {
	//	for logEntry := range logStream {
	//		fmt.Println("std out>", logEntry.Message)
	//		g.updateFromStdIn(logEntry)
	//	}
	//}(csgoState)

	//// pass cli and playbooks to csgo's stdin
	//go func() {
	//	for {
	//		select {
	//		case s := <-csgoCLI:
	//			csgoStdin <- s
	//		case s := <-csgoCMD:
	//			csgoStdin <- s
	//		}
	//	}
	//}()
	//
	//// read cli input
	//go func() {
	//	cliReader := bufio.NewReader(os.Stdin)
	//
	//	for {
	//		text, _ := cliReader.ReadString('\n')
	//		text = strings.Trim(strings.TrimSuffix(text, "\n"), "")
	//
	//		if len(text) > 0 {
	//			switch cmd := strings.ToLower(text); cmd {
	//			case "!lo3":
	//				playbook.LiveOnThree(csgoCLI)
	//			case "!knife":
	//				csgoCLI <- "exec kniferound"
	//			case "!reset":
	//				playbook.Reset(csgoCLI)
	//			default:
	//				csgoCLI <- text
	//			}
	//		}
	//	}
	//}()
	//
	////srcds.WrapProc(srcdsLaunchArgs, csgoStdin, logStream)
}
