package csgo

import (
	"fmt"
	"log"
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

type listeners struct {
	teamScored  func(TeamScored) bool
	teamSideSet func(TeamSideSet) bool
}

// CSGO represents the state of a CSGO server
type CSGO struct {
	cmdIn       chan string
	currentMap  *mapState
	gameMode    GameMode
	maps        []mapState
	mpTeamname1 string
	mpTeamname2 string
	spectators  srcds.Clients
	srcds       *srcds.SRCDS
	listeners   listeners
}

// GameMode determines the rulesets used by a CSGO server.
type GameMode byte

type WorldTrigger byte

// New creates a CSGO server
func New(gameMode GameMode, scenarios ...Scenario) (*CSGO, error) {
	game := CSGO{
		cmdIn:    make(chan string, 6),
		gameMode: gameMode,
	}

	game.srcds.AddCvarWatch("mp_do_warmup_period", "mp_maxrounds", "mp_overtime_enable", "mp_overtime_maxrounds", "mp_warmup_pausetimer")
	game.srcds.AddLaunchArg(gameMode.launchArgs()...)
	game.srcds.AddLaunchArg("-tickrate 128", "+sv_lan 1", "-norestart") //TODO: add "-nobots"
	game.srcds.AddLogProcessor(game.processLogEntry)

	for _, scenario := range scenarios {
		game = *scenario(&game)
	}

	return &game, nil
}

// Start begins a CSGO server
func (g *CSGO) Start() {
	g.srcds.Start(g.cmdIn)
}

func (g *CSGO) clientJoinedCT(player srcds.Client) {
	c := playerFromSrcdsClient(player)
	g.currentMap.PlayerJoinedCT(c)
}

func (g *CSGO) clientJoinedSpectator(client srcds.Client) {
	g.spectators.ClientJoined(client)
}

func (g *CSGO) clientJoinedTerrorist(player srcds.Client) {
	c := playerFromSrcdsClient(player)
	g.currentMap.PlayerJoinedTerrorist(c)
}

func (g *CSGO) clientDropped(client srcds.Client) {
	g.spectators.ClientDropped(client)

	p := playerFromSrcdsClient(client)
	g.currentMap.PlayerDropped(p)
}

func (m GameMode) launchArgs() []string {
	switch m {
	case ClassicCasual:
		return []string{"-game csgo", "+game_type 0", "+game_mode 0"}
	case ArmsRace:
		return []string{"-game csgo", "+game_type 1", "+game_mode 0"}
	case Demolition:
		return []string{"-game csgo", "+game_type 1", "+game_mode 1"}
	case Deathmatch:
		return []string{"-game csgo", "+game_type 1", "+game_mode 2"}
	default:
		fallthrough
	case ClassicCompetitive:
		return []string{"-game csgo", "+game_type 0", "+game_mode 1"}
	}
}

func (g *CSGO) mapChanged(mapName string) {
	i := len(g.maps)

	if i > 0 {
		g.maps[i-1].ended = time.Now()
	}

	g.maps = append(g.maps, mapState{
		name:    mapName,
		started: time.Now()},
	)

	g.currentMap = &g.maps[i]

	if (len(g.mpTeamname1)) == 0 {
		g.mpTeamname1 = "mp_team_1"
	}
	g.currentMap.mpTeam1.SetName(g.mpTeamname1)

	if (len(g.mpTeamname2)) == 0 {
		g.mpTeamname2 = "mp_team_2"
	}
	g.currentMap.mpTeam2.SetName(g.mpTeamname2)
}

func debugPrint(t, s string) {
	fmt.Println("============================================")
	fmt.Println("\t\t", t)
	fmt.Println(s)
	fmt.Println("============================================")
}

func (g *CSGO) processLogEntry(le srcds.LogEntry) (keepProcessing bool) {

	// client did something
	if strings.HasPrefix(le.Message, `"`) {
		_, err := parsePlayerSay(le)
		if err != nil {
			// process player said
			return true
		}

		client, err := parseClientConnected(le)
		if err != nil {
			g.clientJoinedSpectator(client)
			return true
		}

		clientDisconnected, err := parseClientDisconnected(le)
		if err != nil {
			g.clientDropped(clientDisconnected.client)
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

		if err != nil {
			switch worldTriggered.trigger {
			case MatchStart:
				g.srcds.RefreshCvars()
			}
		}

		return true
	}

	mapName, err := parseLoadingMap(le)
	if err == nil {
		g.mapChanged(mapName)
		return true
	}

	return true
}

func (g *CSGO) teamScored(m TeamScored) {
	switch m.teamAffiliation {
	case "CT":
		g.currentMap.CTSetScore(m.teamScore)
	case "TERRORIST":
		g.currentMap.TerroristSetScore(m.teamScore)
	default:
		log.Println("UNABLE TO teamScored() for affiliation '" + m.teamAffiliation + "'")
	}

	if g.listeners.teamScored != nil {
		g.listeners.teamScored(m)
	}
}

func (g *CSGO) teamSetSide(m TeamSideSet) {
	if m.teamAffiliation == "CT" {
		if g.currentMap.terrorist().name == m.teamName {
			g.currentMap.TeamsSwappedSides()
		}
	} else if m.teamAffiliation == "TERRORIST" {
		if g.currentMap.ct().name == m.teamName {
			g.currentMap.TeamsSwappedSides()
		}
	}

	if g.listeners.teamSideSet != nil {
		g.listeners.teamSideSet(m)
	}
}

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

func (g *CSGO) maxOvertimeRounds() int {
	if s, found := g.srcds.GetCvar("mp_overtime_maxrounds"); found {
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}

	return 6
}

func (g *CSGO) maxRounds() int {
	if s, found := g.srcds.GetCvar("mp_maxrounds"); found {
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}

	return 30
}
