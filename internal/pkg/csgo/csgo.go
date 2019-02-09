package csgo

import (
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

type GameMode int

const (
	CasualCompetitive GameMode = iota
	ClassicCompetitive
	ArmsRace
	Demolition
	Deathmatch
)

type CSGO struct {
	currentMap     *mapState
	maps           []mapState
	mpTeamname1    string
	mpTeamname2    string
	spectators     srcds.Clients
	gameMode       GameMode
	srcds          *srcds.SRCDS
	cmdIn          chan string
	downProcessors []func(*CSGO, srcds.LogEntry)
}

func New(server *srcds.SRCDS, mode GameMode, scenarios ...Scenario) (*CSGO, error) {
	game := CSGO{
		cmdIn:          make(chan string, 6),
		downProcessors: make([]func(*CSGO, srcds.LogEntry), len(scenarios)),
		srcds:          server,
	}

	game.srcds.AddCvarWatch("mp_do_warmup_period", "mp_maxrounds", "mp_overtime_enable", "mp_overtime_maxrounds", "mp_warmup_pausetimer")
	game.srcds.AddLaunchArg("-game csgo", argsFromGameMode(mode), "-tickrate 128", "+sv_lan 1", "-norestart") //TODO: add "-nobots"
	game.srcds.AddLogProcessor(game.processLogEntry)

	for _, scenario := range scenarios {
		game = *scenario(&game)
	}

	return &game, nil
}

func (g *CSGO) Start() {
	g.srcds.Start(g.cmdIn)
}

func argsFromGameMode(mode GameMode) string {
	switch mode {
	case CasualCompetitive:
		return "+game_type 0 +game_mode 0"
	case ArmsRace:
		return "+game_type 1 +game_mode 0"
	case Demolition:
		return "+game_type 1 +game_mode 1"
	case Deathmatch:
		return "+game_type 1 +game_mode 2"
	default:
		fallthrough
	case ClassicCompetitive:
		return "+game_type 0 +game_mode 1"
	}
}

func (g *CSGO) processLogEntry(le srcds.LogEntry) (keepProcessing bool) {
	if strings.HasPrefix(le.Message, `Started map`) {
		g.srcds.ReconcileCvars()
	}

	// update game state
	return true
}

func (g *CSGO) ClientJoinedCT(player srcds.Client) {
	c := playerFromSrcdsClient(player)
	g.currentMap.PlayerJoinedCT(c)
}

func (g *CSGO) ClientJoinedSpectator(client srcds.Client) {
	g.spectators.ClientJoined(client)
}

func (g *CSGO) ClientJoinedTerrorist(player srcds.Client) {
	c := playerFromSrcdsClient(player)
	g.currentMap.PlayerJoinedTerrorist(c)
}

func (g *CSGO) ClientDropped(client srcds.Client) {
	g.spectators.ClientDropped(client)

	p := playerFromSrcdsClient(client)
	g.currentMap.PlayerDropped(p)
}

func (g *CSGO) ctWonRound() {
	g.currentMap.CTWonRound()
}

func (g *CSGO) RoundNumber() byte {
	return g.currentMap.roundNumber
}

func (g *CSGO) TeamsSwappedSides() {
	g.currentMap.TeamsSwappedSides()
}

func (g *CSGO) terroristWonRound() {
	g.currentMap.TerroristWonRound()
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
