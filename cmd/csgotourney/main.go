package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
	"github.com/lacledeslan/sourceseer/internal/pkg/srcds/csgo"
)

var (
	bracket = flag.String("bracket", "", "The tournament bracket this server is for")
	ctName  = flag.String("mp_teamname_1", "", "The name of the team that will select CT on connection")
	tName   = flag.String("mp_teamname_2", "", "The name of the team that will select Terrorist on connection")
	pass    = flag.String("pass", "", "The server's password")
	rcon    = flag.String("rcon_pass", "", "The server's rcon password")
	tvPass  = flag.String("tv_pass", "", "The server's tv password")
)

func main() {
	flag.Parse()
	maps := flag.Args()

	tourneyBracket := strings.TrimSpace(*bracket)
	if len(tourneyBracket) == 0 {
		fmt.Fprint(os.Stderr, "Argument bracket must be provided!\n\n")
		fmt.Fprint(os.Stderr, "\tExample: -bracket 12B\n\n")
		os.Exit(87)
	}

	mpTeamname1 := strings.TrimSpace(*ctName)
	if len(mpTeamname1) == 0 {
		fmt.Fprint(os.Stderr, "Argument mp_teamname_1 must be provided!\n\n")
		fmt.Fprint(os.Stderr, "\tExample: -mp_teamname_1 red\n\n")
		os.Exit(87)
	}

	mpTeamname2 := strings.TrimSpace(*tName)
	if len(mpTeamname2) == 0 {
		fmt.Fprint(os.Stderr, "Argument mp_teamname_2 must be provided!\n\n")
		fmt.Fprint(os.Stderr, "\tExample: -mp_teamname_2 blu\n\n")
		os.Exit(87)
	}

	if strings.ToLower(mpTeamname1) == strings.ToLower(mpTeamname2) {
		fmt.Fprint(os.Stderr, "Both team names cannot be the same!\n\n")
		os.Exit(87)
	}

	svPassword := strings.TrimSpace(*pass)
	if len(svPassword) == 0 {
		fmt.Fprint(os.Stderr, "Argument pass must be provided!\n\n")
		fmt.Fprint(os.Stderr, "\tExample: -pass 72rivers\n\n")
		os.Exit(87)
	}
	svPassword = `+sv_password "` + svPassword + `"`

	rconPassword := strings.TrimSpace(*rcon)
	if len(svPassword) == 0 {
		fmt.Fprint(os.Stderr, "Argument rcon_pass must be provided!\n\n")
		fmt.Fprint(os.Stderr, "\tExample: -rcon_pass BulldogsFancy957Cupcakes\n\n")
		os.Exit(87)
	}
	rconPassword = `+rcon_password  "` + rconPassword + `"`

	tvPassword := strings.TrimSpace(*tvPass)
	if len(tvPassword) == 0 {
		fmt.Fprint(os.Stderr, "Argument tv_pass must be provided!\n\n")
		fmt.Fprint(os.Stderr, "\tExample: -tv_pass CowsH4teTennis\n\n")
		os.Exit(87)
	}
	tvPassword = `+tv_password "` + tvPassword + ` " +tv_relaypassword "` + tvPassword + `"`

	if l := len(maps); l == 0 || l%2 == 0 {
		fmt.Fprint(os.Stderr, "A positive, odd-number of maps must be provided!\n\n")
		fmt.Fprint(os.Stderr, "\tExample: -mp_teamname_1 red -mp_teamname_2 blu de_inferno de_biome de_inferno\n\n")
		os.Exit(87)
	}

	csgoTourney, err := csgo.New(csgo.ClassicCompetitive, csgo.UseTeamNames(mpTeamname1, mpTeamname2), csgo.ClinchableMapCycle(maps))
	if err != nil {
		fmt.Fprint(os.Stderr, "Was unable to create a CSGO instance!\n\n")
		os.Exit(87)
	}

	csgoTourney.AddLaunchArg(svPassword, rconPassword, tvPassword)

	var osArgs []string
	if _, err := os.Stat("/app/srcds_run"); err == nil {
		osArgs = []string{"/app/srcds_run"} // we're inside docker
	} else {
		for i := 5; i >= 0; i-- {
			time.Sleep(4 * time.Second)
			fmt.Println("RUNNING LOCAL IN", i)
		}

		switch os := runtime.GOOS; os {
		case "windows":
			osArgs = append(osArgs, "powershell.exe", "-NonInteractive", "-Command")
		}

		osArgs = append(osArgs, "docker", "run", "-i", "--rm", "--net=host", `--entrypoint "/bin/bash"`, "lacledeslan/gamesvr-csgo-tourney:hasty", "/app/srcds_run")
	}

	server, err := srcds.New(csgoTourney, osArgs)

	if err != nil {
		fmt.Fprint(os.Stderr, "Unable to create a Source Dedicated Server!\n\n")
		fmt.Fprint(os.Stderr, "\tReason: ", err, "\n\n")
		os.Exit(-1)
	}

	if csgoTourney == nil {
		fmt.Fprint(os.Stderr, "Unable to create a CSGO Tournament server!\n\n")
		fmt.Fprint(os.Stderr, "\tReason: ", err, "\n\n")
		os.Exit(-1)
	}

	server.Start()

	fmt.Print("\n\nfin.\n\n")

	os.Exit(0)
}
