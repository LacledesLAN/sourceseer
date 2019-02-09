package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/lacledeslan/sourceseer/internal/pkg/csgo"
	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

func main() {
	ctName := flag.String("mp_teamname_1", "", "The name of the team starting on CT")
	tName := flag.String("mp_teamname_2", "", "The name of the team starting on Terrorist")
	flag.Parse()
	maps := flag.Args()

	mpTeamname1 := strings.TrimSpace(*ctName)
	if len(strings.TrimSpace(mpTeamname1)) == 0 {
		fmt.Print("Argument mp_teamname_1 must be provided!\n\n")
		fmt.Print("\tExample: -mp_teamname_1 red\n\n")
		os.Exit(87)
	}

	mpTeamname2 := strings.TrimSpace(*tName)
	if len(strings.TrimSpace(mpTeamname2)) == 0 {
		fmt.Print("Argument mp_teamname_2 must be provided!\n\n")
		fmt.Print("\tExample: -mp_teamname_2 blu\n\n")
		os.Exit(87)
	}

	if l := len(maps); l == 0 || l%2 == 0 {
		fmt.Print("A positive, odd-number of maps must be provided!\n\n")
		fmt.Print("\tExample: -mp_teamname_1 red -mp_teamname_2 blu de_inferno de_biome de_inferno\n\n")
		os.Exit(87)
	}

	var osArgs []string
	if _, err := os.Stat("/app/srcds_run"); err == nil {
		osArgs = []string{"/app/srcds_run"} // we're inside docker
	} else {
		switch os := runtime.GOOS; os {
		case "windows":
			osArgs = append(osArgs, "powershell.exe", "-NonInteractive", "-Command")
		}

		osArgs = append(osArgs, "docker", "run", "-i", "--rm", "-p 27015:27015", "-p 27015:27015/udp", "lltest/gamesvr-csgo-tourney", "./srcds_run")
	}

	server, err := srcds.New(osArgs)

	if err != nil {
		fmt.Print("Unable to create a Source Dedicated Server!\n\n")
		fmt.Print("\tReason: ", err, "\n\n")
		os.Exit(-1)
	}

	csgoTourney, err := csgo.New(&server, csgo.ClassicCompetitive, csgo.CompetitiveWarmUp(mpTeamname1, mpTeamname2), csgo.ClinchableMapCycle(maps))

	if csgoTourney == nil {
		fmt.Print("Unable to create a CSGO Tournament server!\n\n")
		fmt.Print("\tReason: ", err, "\n\n")
		os.Exit(-1)
	}

	csgoTourney.Start()

	fmt.Print("\n\nfin.\n\n")

	os.Exit(0)
}
