package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
	"github.com/lacledeslan/sourceseer/internal/pkg/srcds/csgo"
)

const (
	alphaNumericChars  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lacledesMaps = "/de_lltest/de_tinyorange/poolday/"
	stockMaps    = "/ar_baggage/ar_dizzy/ar_monastery/ar_shoots/cs_agency/cs_assault/cs_italy/cs_militia/cs_office/de_austria/de_bank/de_biome/de_cache/de_canals/de_cbble/de_dust2/de_inferno/de_lake/de_mirage/de_nuke/de_overpass/de_safehouse/de_shortnuke/de_stmarc/de_subzero/de_sugarcane/de_train/"
)

func main() {
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "+") {
			os.Args[i] = "--" + arg[1:]
		}
	}

	var args []string
	var csgoArgs csgo.Args

	args, err := flags.ParseArgs(&csgoArgs, os.Args)
	if err != nil {
		panic(err)
	}

	if len(csgoArgs.TeamName1) < 1 {
		fmt.Fprint(os.Stderr, "Argument mp_teamname_1 must be provided!\n\n")
		fmt.Fprint(os.Stderr, "\tExample: --mp_teamname_1 red\n\n")
		os.Exit(87)
	}

	if len(csgoArgs.TeamName2) < 1 {
		fmt.Fprint(os.Stderr, "Argument mp_teamname_2 must be provided!\n\n")
		fmt.Fprint(os.Stderr, "\tExample: --mp_teamname_2 blu\n\n")
		os.Exit(87)
	}

	if strings.ToLower(csgoArgs.TeamName1) == strings.ToLower(csgoArgs.TeamName2) {
		fmt.Fprint(os.Stderr, "The values for `--mp_teamname_1` and `--mp_teamname_2` cannot match!\n\n")
		os.Exit(87)
	}

	if csgoArgs.UseRemoteConsole && len(csgoArgs.RConPassword) < 1 {
		fmt.Fprint(os.Stderr, "When remote console is enabled argument rcon_pass must be provided!\n\n")
		fmt.Fprint(os.Stderr, "\tExample: --rcon_password BulldogsFancy957Cupcakes\n\n")
		os.Exit(87)
	}

	if len(csgoArgs.TVPassword) < 1 {
		b := make([]byte, 16)
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}
		csgoArgs.TVPassword = string(b)
	}

	if len(csgoArgs.TVRelayPassword) < 1 {
		b := make([]byte, 16)
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}
		csgoArgs.TVRelayPassword = string(b)
	}

	maps := args[1:]
	if l := len(maps); l == 0 || l%2 == 0 {
		fmt.Fprint(os.Stderr, "A positive, odd-number of maps must be provided!\n\n")
		fmt.Fprint(os.Stderr, "\tExample: --mp_teamname_1 red --mp_teamname_2 blu de_inferno de_biome de_inferno\n\n")
		os.Exit(87)
	}

	csgoTourney, err := csgo.New(csgo.ClassicCompetitive, csgo.UseTeamNames(csgoArgs.TeamName1, csgoArgs.TeamName2), csgo.ClinchableMapCycle(maps))
	if err != nil {
		fmt.Fprint(os.Stderr, "Was unable to create a CSGO instance!\n\n")
		os.Exit(87)
	}

	if csgoArgs.Bots {
		fmt.Println("Allowing bots!")
	}

	var osArgs []string
	if _, err := os.Stat("/app/srcds_run"); err == nil {
		for _, bspFile := range maps {
			if _, err := os.Stat("/app/csgo/maps/" + bspFile + ".bsp"); os.IsNotExist(err) || err != nil {
				fmt.Fprint(os.Stderr, "Could not find file for map `", bspFile, "`!\n\n")
				os.Exit(87)
			}
		}

		osArgs = []string{"/app/srcds_run"} // we're inside docker
	} else {
		for i := 3; i >= 0; i-- {
			time.Sleep(4 * time.Second)
			fmt.Println("RUNNING LOCAL IN", i)
		}

		for i, m := range maps {
			if m != strings.ToLower(m) {
				args[i] = strings.ToLower(m)
			}

			if err := validateStockMapNames(maps[i]); err != nil {
				fmt.Fprint(os.Stderr, err, "\n\n")
				os.Exit(87)
			}
		}

		switch os := runtime.GOOS; os {
		case "windows":
			osArgs = append(osArgs, "powershell.exe", "-NonInteractive", "-Command")
		}

		osArgs = append(osArgs, "docker", "run", "-i", "--rm", "--net=host", `--entrypoint "/bin/bash"`, "lacledeslan/gamesvr-csgo", "/app/srcds_run")
	}

	server, err := srcds.New(csgoTourney)

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

	linkStandardIn(server.CmdIn)

	server.Start(osArgs)

	fmt.Print("\n\nfin.\n\n")

	os.Exit(0)
}

func linkStandardIn(cmdIn chan string) {
	if _, err := os.Stdin.Stat(); err == nil {
		go func() {
			time.Sleep(5 * time.Second)

			s := bufio.NewScanner(os.Stdin)
			defer os.Stdin.Close()

			fmt.Println("<<<< Now accepting input from command line >>>>")

			for s.Scan() {
				text := s.Text()
				cmdIn <- text
			}
		}()
	}
}

// validateStockMapNames tests if the provide map names are all valid stock maps
func validateStockMapNames(mapName string) error {
	if strings.Index(stockMaps, "/"+mapName+"/") == -1 && strings.Index(lacledesMaps, "/"+mapName+"/") == -1 {
		return errors.New("\"" + mapName + "\" is not a valid map")
	}

	return nil
}
