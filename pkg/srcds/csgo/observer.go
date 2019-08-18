package csgo

import (
	"fmt"
	"io"
	"strings"

	"github.com/LacledesLAN/sourceseer/pkg/srcds"
)

//Observer for watching CSGO log streams
type Observer interface {
	DebugDump()
	Start()
}

// NewReader for observing streaming CSGO data
func NewReader(r io.Reader, mpHalftime, mpMaxRounds, mpMaxOvertimeRounds int) Observer {
	o := &observer{
		srcdsObserver: srcds.NewReader(r),
	}

	o.srcdsObserver.AddCvarWatcherDefault("mp_halftime", string(mpHalftime))
	o.srcdsObserver.AddCvarWatcherDefault("mp_maxrounds", string(mpMaxRounds))
	o.srcdsObserver.AddCvarWatcherDefault("mp_overtime_maxrounds", string(mpMaxOvertimeRounds))

	return o
}

func (o observer) DebugDump() {
	fmt.Print("=======================================================================")
	for i, m := range o.game.matches {
		fmt.Printf("\nMatch #%02d on %q started %q; %q (mp_team1) vs %q (mp_team2).\n", i+1, m.mapName, m.started, o.game.mpTeamname1, o.game.mpTeamname2)

		fmt.Printf("\tTotal rounds played: %d.\n\n", len(m.rounds))

		for j, r := range m.rounds {
			fmt.Printf("\t%02d - Won by %q as %q via trigger %q.\n", j+1, r.winningTeam, r.winningAffiliation, r.winningTrigger)
		}
	}

	fmt.Println("\nUNASSIGNED")
	for _, p := range o.players.unassigned {
		fmt.Printf("\t%+v\n", p)
	}

	fmt.Println("\nMP TEAM 1")
	for _, p := range o.players.mpTeam1 {
		fmt.Printf("\t%+v\n", p)
	}

	fmt.Println("\nMP TEAM 2")
	for _, p := range o.players.mpTeam2 {
		fmt.Printf("\t%+v\n", p)
	}

	fmt.Println("\n=======================================================================")
}

// Starts the observer
func (o *observer) Start() {
	for le := range o.srcdsObserver.Start() {
		o.processLogEntry(le)
	}
}

type observer struct {
	players struct {
		mpTeam1    srcds.Clients
		mpTeam2    srcds.Clients
		unassigned srcds.Clients
	}
	game          gameInfo
	srcdsObserver srcds.Observer
}

func (o *observer) processLogEntry(log srcds.LogEntry) {
	if clientLog, ok := srcds.ParseClientLogEntry(log); ok {
		if _, ok := parseClientSay(clientLog); ok {
			// TODO: process the client saying something
			return
		}

		if ok := srcds.ParseClientConnected(clientLog); ok {
			o.playerJoined(unassigned, clientLog.Client)
			return
		}

		if m, ok := parseClientSetAffiliation(clientLog); ok {
			o.playerJoined(m.to, clientLog.Client)
			// TODO: Add ability to swap players
			return
		}

		if _, ok := srcds.ParseClientDisconnected(clientLog); ok {
			//o.playerDropped(clientLog.Client)
			//TODO: Do we even need to drop players?
		}

		return
	}

	if worldLog, ok := parseWorldTrigger(log); ok {
		if mapName, ok := parseWorldTriggerMatchStart(worldLog); ok {
			o.game.nextMatch(mapName)
		}

		if parseWorldTriggerRoundStart(worldLog) {
		}

		if parseWorldTriggerRoundEnd(worldLog) {
		}

		if parseWorldTriggerGameCommencing(worldLog) {
		}
	}

	if strings.HasPrefix(log.Message, "Team") {
		if msg, ok := parseTeamTriggered(log); ok {
			team := o.getTeam(msg.affiliation)
			o.game.setRoundWinner(msg.affiliation, team, msg.trigger)
		}

		if msg, ok := parseTeamSetName(log); ok {
			switch teamName := o.getTeam(msg.affiliation); teamName {
			case mpTeam1:
				o.game.mpTeamname1 = msg.teamName
			case mpTeam2:
				o.game.mpTeamname2 = msg.teamName
			}

			return
		}

		//return
	}

	if parseStartingFreezePeriod(log) {
		return
	}
}

func (o *observer) getTeam(a affiliation) team {
	if a == unassigned {
		return ""
	}

	mpHalftime, _ := o.srcdsObserver.TryCvarAsInt("mp_halftime", defaultMpHalftime)
	mpMaxrounds, _ := o.srcdsObserver.TryCvarAsInt("mp_maxrounds", defaultMpMaxrounds)
	mpOvertimeMaxrounds, _ := o.srcdsObserver.TryCvarAsInt("mp_overtime_maxrounds", defaultMpOvertimeMaxrounds)
	completedRounds := o.game.currentMatchLastCompletedRound()

	if calculateSidesAreCurrentlySwitched(mpHalftime, mpMaxrounds, mpOvertimeMaxrounds, completedRounds) {
		if a == counterterrorist {
			return mpTeam2
		}

		return mpTeam1
	}

	if a == counterterrorist {
		return mpTeam1
	}

	return mpTeam2
}

func (o *observer) playerDropped(c srcds.Client) {
	o.players.mpTeam1.ClientDropped(c)
	o.players.mpTeam2.ClientDropped(c)
	o.players.unassigned.ClientDropped(c)
}

func (o *observer) playerJoined(a affiliation, c srcds.Client) {
	switch team := o.getTeam(a); team {
	case mpTeam1:
		o.players.mpTeam1.ClientJoined(c)
		o.players.mpTeam2.ClientDropped(c)
		o.players.unassigned.ClientDropped(c)
	case mpTeam2:
		o.players.mpTeam1.ClientDropped(c)
		o.players.mpTeam2.ClientJoined(c)
		o.players.unassigned.ClientDropped(c)
	default:
		o.players.mpTeam1.ClientDropped(c)
		o.players.mpTeam2.ClientDropped(c)
		o.players.unassigned.ClientJoined(c)
	}
}
