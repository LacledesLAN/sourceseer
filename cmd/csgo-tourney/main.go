package main

import (
	"flag"
	"os"
	"strings"

	"github.com/lacledeslan/sourceseer/csgo"
)

func main() {
	ctName := flag.String("mp_teamname_1", "", "The name of the team starting on CT")
	tName := flag.String("mp_teamname_2", "", "The name of the team starting on Terrorist")
	flag.Parse()
	maps := flag.Args()

	mpTeamname1 := strings.TrimSpace(*ctName)
	if len(strings.TrimSpace(mpTeamname1)) == 0 {
		panic("mp_teamname_1 must be provided.")
	}

	mpTeamname2 := strings.TrimSpace(*tName)
	if len(strings.TrimSpace(mpTeamname2)) == 0 {
		panic("mp_teamname_2 must be provided.")
	}

	if len(maps) == 0 {
		panic("At least one map must be provided")
	}

	if len(maps)%2 == 0 {
		panic("Must provide an odd number of maps")
	}

	csgo.StartTourney(mpTeamname1, mpTeamname2, maps)

	os.Exit(0)
}
