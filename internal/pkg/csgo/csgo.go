package csgo

import (
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

// GameMode determines the rulesets used by a CSGO server.
type GameMode byte

const (
	// ClassicCasual features a simplified economy, no team damage, and all players automatically receive armor and defuse kits.
	ClassicCasual GameMode = iota + 1

	// ClassicCompetitive is a classic mode where two teams compete in a best-of match using standard competitive Counter-Strike rules.
	ClassicCompetitive

	// ArmsRace is a gun-progression mode where players gain new weapons after registering a kill and work their way through each weapon in the game. Get a kill with the final weapon, the golden knife, and win the match!
	ArmsRace

	// Demolition is a round-based mode where players take turns attacking and defending a single bombsite in a series of maps designed for fast-paced gameplay.
	Demolition

	// Deathmatch is a fast-paced casual mode where every player is for themselves, respawn instantly, and have a unlimited amount of time to buy weapons.
	Deathmatch
)

// CSGO represents the state of a CSGO server
type CSGO struct {
	currentMap  *mapState
	gameMode    GameMode
	maps        []mapState
	mpTeamname1 string
	mpTeamname2 string
	spectators  srcds.Clients
	srcds       *srcds.SRCDS
	cmdIn       chan string
}

// New creates a CSGO server
func New(server *srcds.SRCDS, gameMode GameMode, scenarios ...Scenario) (*CSGO, error) {
	game := CSGO{
		cmdIn:    make(chan string, 6),
		gameMode: gameMode,
		srcds:    server,
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

func (g *CSGO) processLogEntry(le srcds.LogEntry) (keepProcessing bool) {
	// A player did something
	if strings.HasPrefix(le.Message, `"`) {
		originator, target := srcds.ExtractClients(le)

		if originator != nil {
			if target == nil {
				if strings.Contains(le.Message, `" switched from team <`) {
					if strings.HasSuffix(le.Message, "<CT>") {
						g.clientJoinedCT(*originator)
					} else if strings.HasSuffix(le.Message, "<TERRORIST>") {
						g.clientJoinedTerrorist(*originator)
					}

					return
				}

				if strings.Contains(le.Message, `" disconnected (reason "`) {
					g.clientDropped(*originator)
				}
			}
		}

		return
	}

	// A team did something
	if strings.HasPrefix(le.Message, "Team") {
		if strings.HasPrefix(le.Message, `Team "CT" scored "`) {
			g.ctWonRound()
		} else if strings.HasPrefix(le.Message, `Team "TERRORIST" scored "`) {
			g.terroristWonRound()
		} else if strings.HasPrefix(le.Message, `Team playing "`) {
			result := teamSetSideRegex.FindStringSubmatch(le.Message)

			if result[1] == "CT" {
				if g.currentMap.terrorist().name == result[2] {
					g.teamsSwappedSides()
				}
			} else if result[1] == "TERRORIST" {
				if g.currentMap.ct().name == result[2] {
					g.teamsSwappedSides()
				}
			}
		}

		return
	}

	if strings.HasPrefix(le.Message, `Started map`) {
		g.srcds.RefreshCvars()
	}

	if strings.HasPrefix(le.Message, "World triggered") {

		if le.Message == `World triggered "Game_Commencing"` {

		}

		// set up
		if strings.HasPrefix(le.Message, `World triggered "Match_Start"`) {
			mapName := worldTriggeredMatchStartRegex.FindStringSubmatch(le.Message)[1]

			if g.currentMap.name != mapName {
				// log as issue?
			}
		}

		return
	}

	if strings.HasPrefix(le.Message, "Game Over:") {
		//result := gameOverRegex.FindStringSubmatch(logEntry.Message)
		//resultScore1 := result[3]
		//resultScore2 := result[4]

		// hook for change level

		return
	}

	if strings.HasPrefix(le.Message, `Loading map "`) {
		mapName := regexBetweenQuotes.FindString(le.Message)
		mapName = strings.Trim(mapName[1:len(mapName)-1], "")

		g.mapChanged(mapName)
		return
	}

	// update game state
	return true
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

func (g *CSGO) ctWonRound() {
	g.currentMap.CTWonRound()
}

func (g *CSGO) RoundNumber() byte {
	return g.currentMap.roundNumber
}

func (g *CSGO) teamsSwappedSides() {
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
