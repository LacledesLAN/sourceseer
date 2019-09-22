package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lacledeslan/sourceseer/pkg/srcds/csgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Print("\n=======================================================================\n\n")
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	//file, err := os.Open(`C:\Workspace\sourceseer\pkg\srcds\csgo\testdata\tourney_3map_clinch.log`)
	file, err := os.Open(`..\..\pkg\srcds\csgo\testdata\tourney_3map_clinch.log`)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := bufio.NewReader(file)

	c := csgo.NewReader(r, 1, 30, 7)

	c.Start()
}
