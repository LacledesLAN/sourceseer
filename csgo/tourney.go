package csgo

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/csgo/playbooks"
	"github.com/lacledeslan/sourceseer/srcds"
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
	if len(maps) == 0 {
		panic("At least one map must be provided")
	}

	if len(maps)%2 == 0 {
		panic("Must provide an odd number of maps")
	}

	if err := validateStockMapNames(maps); err != nil {
		panic(err)
	}

	srcdsLaunchArgs := []string{"-game csgo", "+game_type 0", "+game_mode 1", "-tickrate 128", "+sv_lan 1"} //TODO: add "-nobots"

	// Process mpTeamname1
	mpTeamname1 = SanitizeTeamName(mpTeamname1)
	if len(mpTeamname1) > 0 {
		srcdsLaunchArgs = append(srcdsLaunchArgs, "+mp_teamname_1 ", mpTeamname1)
	}

	// Process mpTeamname2
	mpTeamname2 = SanitizeTeamName(mpTeamname2)
	if len(mpTeamname2) > 0 {
		srcdsLaunchArgs = append(srcdsLaunchArgs, "+mp_teamname_2 ", mpTeamname2)
	}

	srcdsLaunchArgs = append(srcdsLaunchArgs, `+hostname "`+HostnameFromTeamNames(mpTeamname1, mpTeamname2)+`"`)
	srcdsLaunchArgs = append(srcdsLaunchArgs, "+map "+maps[0])

	logStream := make(chan srcds.LogEntry, 9)
	defer close(logStream)

	csgoState := NewGameState()

	go func(g *gameState) {
		for logEntry := range logStream {
			fmt.Println("std out>", logEntry.Message)
			g.updateFromStdIn(logEntry)
		}
	}(csgoState)

	cmdStream := make(chan string, 6)
	defer close(cmdStream)
	csgoCLI := make(chan string)
	defer close(csgoCLI)
	csgoCMD := make(chan string)
	defer close(csgoCMD)

	// pass cli and playbooks to csgo's stdin
	go func() {
		for {
			select {
			case s := <-csgoCLI:
				csgoStdin <- s
			case s := <-csgoCMD:
				csgoStdin <- s
			}
		}
	}()

	// read cli input
	go func() {
		cliReader := bufio.NewReader(os.Stdin)

		for {
			text, _ := cliReader.ReadString('\n')
			text = strings.Trim(strings.TrimSuffix(text, "\n"), "")

			if len(text) > 0 {
				switch cmd := strings.ToLower(text); cmd {
				case "!lo3":
					playbooks.LiveOnThree(csgoCLI)
				case "!knife":
					csgoCLI <- "exec kniferound"
				case "!reset":
					playbooks.Reset(csgoCLI)
				default:
					csgoCLI <- text
				}
			}
		}
	}()

	srcds.WrapProc(srcdsLaunchArgs, csgoStdin, logStream)
}
