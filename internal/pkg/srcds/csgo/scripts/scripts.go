package scripts

import (
	"strconv"
	"time"
)

var (
	lo3Messages = [64]string{"|--------LIVE--------|", "|-------LIVE---------|", "|------LIVE----------|", "|-----LIVE-----------|", "|-----LIV-E----------|", "|-----LIV--E---------|",
		"|-----LIV---E--------|", "|-----LIV----E-------|", "|-----LIV-----E------|", "|-----LI-V-----E-----|", "|-----LI--V----E-----|", "|-----LI---V---E-----|", "|-----LI----V--E-----|",
		"|-----LI-----V-E-----|", "|-----L-I-----VE-----|", "|-----L--I----VE-----|", "|-----L---I---VE-----|", "|-----L----I--VE-----|", "|-----L------IVE-----|", "|------L-----IVE-----|",
		"|-------L----IVE-----|", "|--------L---IVE-----|", "|---------L--IVE-----|", "|----------L-IVE-----|", "|-----------LIVE-----|", "|----------L-IVE-----|", "|---------L--IVE-----|",
		"|--------L---IVE-----|", "|-------L----IVE-----|", "|------L-----IVE-----|", "|-----L------IVE-----|", "|-----L-----I-VE-----|", "|-----L----I--VE-----|", "|-----L---I---VE-----|",
		"|-----L--I----VE-----|", "|-----L-I-----VE-----|", "|-----LI------VE-----|", "|-----LI-----V-E-----|", "|-----LI----V--E-----|", "|-----LI---V---E-----|", "|-----LI--V----E-----|",
		"|-----LI-V-----E-----|", "|-----LIV------E-----|", "|-----LIV-----E------|", "|-----LIV----E-------|", "|-----LIV---E--------|", "|-----LIV--E---------|", "|-----LIVE-----------|",
		"|-----LIV-E----------|", "|-----LIV--E---------|", "|-----LIV---E--------|", "|-----LIV----E-------|", "|-----LI-V---E-------|", "|-----LI--V--E-------|", "|-----LI---V-E-------|",
		"|-----LI----VE-------|", "|-----L-I---VE-------|", "|-----L--I--VE-------|", "|-----L---I-VE-------|", "|-----L----IVE-------|", "|------L---IVE-------|", "|-------L--IVE-------|",
		"|--------L-IVE-------|", "|---------LIVE-------|"}
)

// liveOnThree is used to ensure players are given time to be ready before life play
func liveOnThree(stdin chan string) {
	for i := 3; i > 0; i-- {
		stdin <- "say [LIVE ON 3 in..." + strconv.Itoa(i) + "!]"
		time.Sleep(1250 * time.Millisecond)
	}

	for i := 1; i < 4; i++ {
		stdin <- "say [Restart " + strconv.Itoa(i) + " of 3!]\n"
		stdin <- "mp_restartgame 2"
		time.Sleep(3250 * time.Millisecond)
	}

	for _, msg := range lo3Messages {
		stdin <- "say " + msg
		time.Sleep(65 * time.Millisecond)
	}

	stdin <- "say GLHF!"
}

func reset(stdin chan string) {
	stdin <- "exec gamemode_competitive"
	stdin <- "exec gamemode_competitive_server"
	stdin <- "mp_restartgame 1"
}
