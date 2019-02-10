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
}

// GameMode determines the rulesets used by a CSGO server.
type GameMode byte

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

func (g *CSGO) SetScore(teamAffiliation, score string) {
	switch affiliation := strings.ToUpper(teamAffiliation); affiliation {
	case "CT":
		g.currentMap.CTSetScore(score)
	case "TERRORIST":
		g.currentMap.TerroristSetScore(score)
	default:
		log.Println("UNABLE TO SETSCORE()")
	}
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

	// A player did something
	if strings.HasPrefix(le.Message, `"`) {
		player, playersTarget := srcds.ExtractClients(le)

		if player != nil {
			if playersTarget == nil {
				if strings.Contains(le.Message, `>" connected, address "`) {
					g.clientJoinedSpectator(*player)
				} else if strings.Contains(le.Message, `>" switched from team <`) {
					if strings.HasSuffix(le.Message, "<CT>") {
						g.clientJoinedCT(*player)
					} else if strings.HasSuffix(le.Message, "<TERRORIST>") {
						g.clientJoinedTerrorist(*player)
					}
				}

				if strings.Contains(le.Message, `" disconnected (reason "`) {
					g.clientDropped(*player)
				}
			}
		}

		return true
	}

	// A team did something
	if strings.HasPrefix(le.Message, "Team") {

		teamScoreUpdate := teamScoredRegex.FindStringSubmatch(le.Message)
		if len(teamScoreUpdate) >= 2 {
			g.SetScore(teamScoreUpdate[1], teamScoreUpdate[2])
			return true
		}

		teamUpdateSide := teamSetSideRegex.FindStringSubmatch(le.Message)
		if len(teamUpdateSide) >= 2 {
			resultAffiliation := strings.ToUpper(teamUpdateSide[1])
			resultTeamName := teamUpdateSide[2]

			if resultAffiliation == "CT" {
				if g.currentMap.terrorist().name == resultTeamName {
					g.teamsSwappedSides()
				}
			} else if resultAffiliation == "TERRORIST" {
				if g.currentMap.ct().name == resultTeamName {
					g.teamsSwappedSides()
				}
			}
		}

		return true
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

	if strings.HasPrefix(le.Message, `Loading map "`) {
		mapName := regexBetweenQuotes.FindString(le.Message)
		mapName = strings.Trim(mapName[1:len(mapName)-1], "")

		g.mapChanged(mapName)
		return
	}

	// update game state
	return true
}

func (g *CSGO) teamsSwappedSides() {
	g.currentMap.TeamsSwappedSides()
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
