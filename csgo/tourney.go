package csgo

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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
	csgoStderr     chan string
	csgoStdin      chan string
	csgoStdout     chan string
	srcdsSafeChars = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
)

// StartTourney starts a csgo tournament
func StartTourney(mpTeamname1 string, mpTeamname2 string, maps []string) {

	if err := validateStockMapNames(maps); err != nil {
		panic(err)
	}
	mapCycle := maps

	gameState := NewGameState()

	launchArgs := []string{"-game csgo", "+game_type 0", "+game_mode 1", "-tickrate 128", "+sv_lan 1"} //TODO: add "-nobots"
	launchArgs = append(launchArgs, "+map "+mapCycle[0])

	mpTeamname1 = SanitizeTeamName(mpTeamname1)
	if len(mpTeamname1) > 0 {
		launchArgs = append(launchArgs, "+mp_teamname_1 ", mpTeamname1)
	}

	mpTeamname2 = SanitizeTeamName(mpTeamname2)
	if len(mpTeamname2) > 0 {
		launchArgs = append(launchArgs, "+mp_teamname_2 ", mpTeamname2)
	}

	launchArgs = append(launchArgs, `+hostname "`+HostnameFromTeamNames(mpTeamname1, mpTeamname2)+`"`)

	csgoStdin = make(chan string, 12)
	defer close(csgoStdin)
	csgoStderr = make(chan string, 2)
	defer close(csgoStderr)
	csgoStdout = make(chan string, 32)
	defer close(csgoStdout)

	go func() {
		for s := range csgoStderr {
			fmt.Println("std err>", s)
			UpdateFromStdErr(&gameState, s)
		}
	}()

	go func() {
		for s := range csgoStdout {
			fmt.Println("std out>", s)
		}
	}()

	csgoCLI := make(chan string)
	//csgoCMD := make(chan string)

	// pass cli and playbooks to csgo's stdin
	go func() {
		for {
			select {
			case s := <-csgoCLI:
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
				case "lo3":
					playbooks.LiveOnThree(csgoCLI)
				case "knife":
					csgoCLI <- "exec kniferound"
				case "reset":
					csgoCLI <- "exec gamemode_competitive"
					csgoCLI <- "exec gamemode_competitive_server"
					csgoCLI <- "mp_restartgame 1"
				default:
					csgoCLI <- text
				}
			}
		}
	}()

	srcds.WrapProc(launchArgs, csgoStdin, csgoStdout, csgoStderr)
}
