package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/LacledesLAN/sourceseer/pkg/srcds/csgo"
)

func main() {
	fmt.Print("STARTING TESTS\n\n")

	//file, err := os.Open(`C:\Workspace\sourceseer\pkg\srcds\csgo\testdata\tourney_3map_clinch.log`)
	file, err := os.Open(`..\..\pkg\srcds\csgo\testdata\tourney_3map_clinch.log`)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := bufio.NewReader(file)

	scanner := bufio.NewScanner(r)

	c := csgo.NewScanner(*scanner, 1, 30, 7)

	//c := srcds.NewScanner(*scanner)

	c.Start()

	c.DebugDump()
}
