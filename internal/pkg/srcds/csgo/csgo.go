package csgo

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

const (
	// ArmsRace is a gun-progression mode where players gain new weapons after registering a kill and work their way through each weapon in the game. Get a kill with the final weapon, the golden knife, and win the match!
	ArmsRace GameMode = iota + 1

	// ClassicCasual features a simplified economy, no team damage, and all players automatically receive armor and defuse kits.
	ClassicCasual

	// ClassicCompetitive is the "original" mode where two teams compete in a best-of match using standard competitive Counter-Strike rules.
	ClassicCompetitive

	// Deathmatch is a fast-paced casual mode where every player is for themselves, respawn instantly, and have a unlimited amount of time to buy weapons.
	Deathmatch

	// Demolition is a round-based mode where players take turns attacking and defending a single bombsite in a series of maps designed for fast-paced, casual gameplay.
	Demolition
)

const (
	UnknownTrigger WorldTrigger = iota
	GameCommencing
	MatchStart
	RoundEnd
	RoundRestarting
	RoundStart
)

// CSGO represents the state of a CSGO server
type CSGO struct {
	cmdIn              chan string
	currentMap         *mapState
	cvars              map[string]srcds.Cvar
	gameMode           GameMode
	launchArgs         []string
	logProcessorStack  LogEntryProcessor
	maps               []mapState
	defaultMpTeamname1 string // needed hack for warmod :/
	defaultMpTeamname2 string // needed hack for warmod :/
}

// GameMode determines the rulesets used by a CSGO server.
type GameMode byte

// LogEntryProcessor represents a function that can parse log entires; returning false when the log entry has been consumed or its effects undone.
type LogEntryProcessor func(srcds.LogEntry) (keepProcessing bool)

type WorldTrigger byte

func (g *CSGO) AddCvarWatch(names ...string) {
	for _, name := range names {
		name = strings.Trim(name, "")

		if len(name) > 0 {
			cvarNameIsUnique := func() bool {
				for key := range g.cvars {
					if key == name {
						return false
					}
				}

				return true
			}()

			if cvarNameIsUnique {
				g.cvars[name] = srcds.Cvar{}
			}
		}
	}
}

// AddLaunchArg to be used when initializing the SRCDS instance.
func (g *CSGO) AddLaunchArg(args ...string) {
	for _, arg := range args {
		arg = strings.Trim(arg, "")
		if len(arg) > 0 {
			g.launchArgs = append(g.launchArgs, arg)
		}
	}
}

// AddLogProcessor to top of the log processor stack.
func (g *CSGO) AddLogProcessor(p LogEntryProcessor) {
	if p != nil {
		prev := g.logProcessorStack

		if prev == nil {
			g.logProcessorStack = p
		} else {
			g.logProcessorStack = func(le srcds.LogEntry) (keepProcessing bool) {
				if !prev(le) {
					return false
				}

				return p(le)
			}
		}
	}
}

// GetCvar value and a boolean as to if the value was found or not.
func (g *CSGO) GetCvar(name string) (value string, found bool) {
	cvar, found := g.cvars[name]

	if found {
		if !cvar.LastUpdate.IsZero() {
			return cvar.Value, found
		}
	}

	return "", found
}

// GetCvarAsInt attempts to return a cvar as an integer
func (g *CSGO) GetCvarAsInt(name string) (value int, err error) {
	v, found := g.GetCvar(name)

	if !found {
		return 0, errors.New("cvar '" + name + "' was not found.")
	}

	return strconv.Atoi(v)
}

// New creates a CSGO server
func New(gameMode GameMode, scenarios ...Scenario) (srcds.Game, error) {
	game := CSGO{
		cmdIn:    make(chan string, 12),
		cvars:    make(map[string]srcds.Cvar),
		gameMode: gameMode,
	}

	switch gameMode {
	case ClassicCasual:
		game.AddLaunchArg("-game csgo", "+game_type 0", "+game_mode 0")
	case ArmsRace:
		game.AddLaunchArg("-game csgo", "+game_type 1", "+game_mode 0")
	case Demolition:
		game.AddLaunchArg("-game csgo", "+game_type 1", "+game_mode 1")
	case Deathmatch:
		game.AddLaunchArg("-game csgo", "+game_type 1", "+game_mode 2")
	default:
		fallthrough
	case ClassicCompetitive:
		game.AddLaunchArg("-game csgo", "+game_type 0", "+game_mode 1")
	}

	game.AddCvarWatch("hostname", "mp_halftime")
	game.AddLaunchArg("-tickrate 128", "+sv_lan 1", "-norestart") //TODO: "-nobots"

	for _, scenario := range scenarios {
		game = *scenario(&game)
	}

	return &game, nil
}

// RefreshCvars triggers SRCDS to echo all watched cvars to the log stream.
func (g *CSGO) RefreshCvars() {
	go func(g *CSGO) {
		for name := range g.cvars {
			g.cmdIn <- name
		}
	}(g)
}

func (g *CSGO) ClientConnected(client srcds.Client) {
	g.clientJoinedSpectator(client)
}

func (g *CSGO) ClientDisconnected(c srcds.ClientDisconnected) {
	g.currentMap.spectators.ClientDropped(c.Client)

	p := playerFromSrcdsClient(c.Client)
	g.currentMap.PlayerDropped(p)
}

func (g *CSGO) clientJoinedCT(player srcds.Client) {
	c := playerFromSrcdsClient(player)

	g.ct().PlayerJoined(c)
	g.terrorist().PlayerDropped(c)
}

func (g *CSGO) clientJoinedSpectator(client srcds.Client) {
	g.currentMap.spectators.ClientJoined(client)
}

func (g *CSGO) clientJoinedTerrorist(player srcds.Client) {
	c := playerFromSrcdsClient(player)

	g.ct().PlayerDropped(c)
	g.terrorist().PlayerJoined(c)
}

func (g *CSGO) CmdSender() chan string {
	return g.cmdIn
}

func (g *CSGO) CvarSet(name, value string) {
	if _, found := g.cvars[name]; found {
		g.cvars[name] = srcds.Cvar{LastUpdate: time.Now(), Value: value}
	}
}

func (g *CSGO) LaunchArgs() []string {
	return g.launchArgs
}

func (g *CSGO) LogReceiver(le srcds.LogEntry) {
	r := g.processLogEntry(le)

	if g.logProcessorStack != nil {
		if r {
			g.logProcessorStack(le)
		}
	}
}

func (g *CSGO) ct() *teamState {
	mpHalftime, _ := g.GetCvarAsInt("mp_halftime")
	mpMaxrounds, _ := g.GetCvarAsInt("mp_maxrounds")
	mpOvertimeMaxrounds, _ := g.GetCvarAsInt("mp_overtime_maxrounds")

	if calculateSidesAreSwitched(mpHalftime, mpMaxrounds, mpOvertimeMaxrounds, g.currentMap.roundsCompleted) {
		return &g.currentMap.mpTeam2
	}

	return &g.currentMap.mpTeam1
}

func (g *CSGO) mapChanged(mapName string) {
	i := len(g.maps)

	if i > 0 {
		g.maps[i-1].ended = time.Now()
	}

	g.maps = append(g.maps, mapState{
		mpTeam1: teamState{
			name:    g.defaultMpTeamname1,
			players: Players{},
		},
		mpTeam2: teamState{
			name:    g.defaultMpTeamname2,
			players: Players{},
		},
		name:    mapName,
		started: time.Now()},
	)

	g.currentMap = &g.maps[i]
}

func (g *CSGO) processLogEntry(le srcds.LogEntry) (keepProcessing bool) {
	// see if a cvar was set
	cvarSet, err := srcds.ParseCvarValueSet(le.Message)
	if err == nil {
		g.CvarSet(cvarSet.Name, cvarSet.Value)
		return
	}

	// client did something
	if strings.HasPrefix(le.Message, `"`) {
		_, err := parsePlayerSay(le)
		if err != nil {
			// process player said
			return true
		}

		clientSwitchedTeam, err := parseClientSwitchedAffiliation(le)
		if err != nil {
			switch strings.ToUpper(clientSwitchedTeam.to) {
			case "CT":
				g.clientJoinedCT(clientSwitchedTeam.client)
			case "TERRORIST":
				g.clientJoinedTerrorist(clientSwitchedTeam.client)
			default:
				g.clientJoinedSpectator(clientSwitchedTeam.client)
			}

			return true
		}

		return true
	}

	// team did something
	if strings.HasPrefix(le.Message, "Team") {
		var err error

		teamScored, err := parseTeamScored(le)
		if err == nil {
			g.teamScored(teamScored)
			return true
		}

		teamUpdateSides, err := parseTeamSetSide(le)
		if err == nil {
			g.teamSetSide(teamUpdateSides)
			return true
		}

		_, err = parseTeamTriggered(le)
		if err == nil {
			return true
		}

		return true
	}

	// The world got triggered
	if strings.HasPrefix(le.Message, "World triggered") {
		worldTriggered, err := parseWorldTriggered(le)

		if err == nil {
			switch worldTriggered.trigger {
			case MatchStart:
				g.RefreshCvars()
				g.currentMap.ResetStats()
			case RoundEnd:
				g.currentMap.roundsCompleted = g.currentMap.roundsCompleted + 1
			}
		}

		return true
	}

	// WarMod Hacks ¯\_ಠ_ಠ_/¯
	if strings.HasPrefix(le.Message, "[WarMod_BFG]") {
		// WarMod drops teamnames during the LO3 before knife fights
		if strings.Contains(le.Message, `, "event": "knife_round_start",`) {
			if len(g.defaultMpTeamname1) > 0 {
				g.cmdIn <- "mp_teamname_1 " + g.defaultMpTeamname1
			}

			if len(g.defaultMpTeamname2) > 0 {
				g.cmdIn <- "mp_teamname_2 " + g.defaultMpTeamname2
			}
		}
	}

	// Map changed
	mapName, err := parseLoadingMap(le)
	if err == nil {
		g.mapChanged(mapName)
		g.cmdIn <- "mp_teamname_1 " + g.defaultMpTeamname1
		g.cmdIn <- "mp_teamname_2 " + g.defaultMpTeamname2
		return true
	}

	return true
}

func (g *CSGO) teamScored(m TeamScored) {
	switch m.teamAffiliation {
	case "CT":
		g.ct().roundsWon = m.teamScore
		g.terrorist().roundsLost = m.teamScore
	case "TERRORIST":
		g.terrorist().roundsWon = m.teamScore
		g.ct().roundsLost = m.teamScore
	}
}

func (g *CSGO) teamSetSide(m TeamSideSet) {
	mpSwapTeams := func() {
		// swap teams
		t := g.currentMap.mpTeam1
		g.currentMap.mpTeam1 = g.currentMap.mpTeam2
		g.currentMap.mpTeam2 = t
	}

	if m.teamAffiliation == "CT" {
		if len(m.teamName) > 0 && m.teamName == g.terrorist().name {
			mpSwapTeams()
		} else {
			g.ct().SetName(m.teamName)
		}
	} else if m.teamAffiliation == "TERRORIST" {
		if len(m.teamName) > 0 && m.teamName == g.ct().name {
			mpSwapTeams()
		} else {
			g.terrorist().SetName(m.teamName)
		}
	}
}

func (g *CSGO) terrorist() *teamState {
	mpHalftime, _ := g.GetCvarAsInt("mp_halftime")
	mpMaxrounds, _ := g.GetCvarAsInt("mp_maxrounds")
	mpOvertimeMaxrounds, _ := g.GetCvarAsInt("mp_overtime_maxrounds")

	if calculateSidesAreSwitched(mpHalftime, mpMaxrounds, mpOvertimeMaxrounds, g.currentMap.roundsCompleted) {
		return &g.currentMap.mpTeam1
	}

	return &g.currentMap.mpTeam2
}