package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lacledeslan/sourceseer/pkg/srcds/csgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Print("\n=======================================================================\n\n")
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	server, _ := csgo.NewServer(`docker`, `run -i --rm --net=host lltest/gamesvr-csgo-tourney ./srcds_run -game csgo +game_type 0 +game_mode 1 -tickrate 128 -console +map de_lltest +sv_lan 1 +mp_teamname_1 "team1" +mp_teamname_2 "team2"`)

	server.Start()

	server.Wait()

	time.Sleep(1 * time.Second)
	fmt.Println("fin")
}
