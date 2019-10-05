package main

import (
	"fmt"
	"os"

	"github.com/lacledeslan/sourceseer/pkg/srcds/csgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Print("\n=======================================================================\n\n")
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	server := csgo.NewServer()
	server.SetExec(`docker`, `run -i --rm --net=host lltest/gamesvr-csgo-tourney ./srcds_run -game csgo +game_type 0 +game_mode 1 -tickrate 128 -console +map de_lltest +sv_lan 1 +mp_teamname_1 "team1" +mp_teamname_2 "team2"`)

	c, err := server.Listen()
	if err != nil {
		panic(err)
	}

	for t := range c {
		fmt.Println(">>>" + t.Message)
	}

	fmt.Println("fin")
	//server.Start()
	//server.Wait()
}
