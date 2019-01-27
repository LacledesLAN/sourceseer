package csgo

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/srcds"
)

var (
	regexBetweenQuotes = regexp.MustCompile("(?:\")([0-9A-Za-z_]*)(\")")
)

type gameState struct {
	currentMap  *mapState
	cvars       map[string]string
	maps        []mapState
	mpTeamname1 string
	mpTeamname2 string
	spectators  srcds.Clients
	started     time.Time
}

func (g *gameState) ClientJoinedCT(player srcds.Client) {
	c := playerFromSrcdsClient(player)
	g.currentMap.PlayerJoinedCT(c)
}

func (g *gameState) ClientJoinedSpectator(client srcds.Client) {
	g.spectators.ClientJoined(client)
}

func (g *gameState) ClientJoinedTerrorist(player srcds.Client) {
	c := playerFromSrcdsClient(player)
	g.currentMap.PlayerJoinedTerrorist(c)
}

func (g *gameState) ClientDropped(client srcds.Client) {
	g.spectators.ClientDropped(client)

	p := playerFromSrcdsClient(client)
	g.currentMap.PlayerDropped(p)
}

func (g *gameState) ctWonRound() {
	g.currentMap.CTWonRound()
}

func (g *gameState) RoundNumber() byte {
	return g.currentMap.roundNumber
}

func NewGameState() gameState {
	return gameState{
		started: time.Now(),
	}
}

func (g *gameState) TeamsSwappedSides() {
	g.currentMap.TeamsSwappedSides()
}

func (g *gameState) terroristWonRound() {
	g.currentMap.TerroristWonRound()
}

// mapChanged updates the gameState when a changelevel has occurred
func (g *gameState) mapChanged(mapName string) {
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

// UpdateFromStdErr updates the gameState from the processes's standard error
func UpdateFromStdErr(g *gameState, err string) {
	log.Print("Received output from standard error:", err)
}

// UpdateFromStdIn updates the gameState from the processe's standard in
func UpdateFromStdIn(g *gameState, logEntry srcds.LogEntry) {

	if strings.HasPrefix(logEntry.Message, `"`) {
		originator, target := srcds.ExtractClients(logEntry)

		if originator != nil {
			// a client did something
			if target == nil {
				if strings.Contains(logEntry.Message, `" switched from team <`) {
					if strings.HasSuffix(logEntry.Message, "<CT>") {
						g.ClientJoinedCT(*originator)
					} else if strings.HasSuffix(logEntry.Message, "<TERRORIST>") {
						g.ClientJoinedTerrorist(*originator)
					}

					return
				}

				if strings.Contains(logEntry.Message, `" disconnected (reason "`) {
					g.ClientDropped(*originator)
				}
			}
		} else {
			result := variableEchoRegex.FindStringSubmatch(logEntry.Message)

			if result != nil {
				updatedCvar(g, result[1], result[2])
			}
		}

		return
	}

	if strings.HasPrefix(logEntry.Message, "Team") {
		if strings.HasPrefix(logEntry.Message, `Team "CT" scored "`) {
			g.ctWonRound()
		} else if strings.HasPrefix(logEntry.Message, `Team "TERRORIST" scored "`) {
			g.terroristWonRound()
		} else if strings.HasPrefix(logEntry.Message, `Team playing "`) {
			result := teamSetSideRegex.FindStringSubmatch(logEntry.Message)

			if result[1] == "CT" {
				if g.currentMap.terrorist().name == result[2] {
					g.TeamsSwappedSides()
				}
			} else if result[1] == "TERRORIST" {
				if g.currentMap.ct().name == result[2] {
					g.TeamsSwappedSides()
				}
			}
		}

		return
	}

	if strings.HasPrefix(logEntry.Message, "World triggered") {
		if strings.HasPrefix(logEntry.Message, `World triggered "Match_Start"`) {
			mapName := worldTriggeredMatchStartRegex.FindStringSubmatch(logEntry.Message)[1]
			g.currentMap.mode = ModePlay

			if g.currentMap.name != mapName {
				// log as issue?
			}
		}

		return
	}

	if strings.HasPrefix(logEntry.Message, "Game Over:") {
		//result := gameOverRegex.FindStringSubmatch(logEntry.Message)
		//resultScore1 := result[3]
		//resultScore2 := result[4]

		// hook for change level

		return
	}

	if strings.HasPrefix(logEntry.Message, `Loading map "`) {
		mapName := regexBetweenQuotes.FindString(logEntry.Message)
		mapName = strings.Trim(mapName[1:len(mapName)-1], "")

		g.mapChanged(mapName)

		return
	}

	if strings.HasPrefix(logEntry.Message, `server_cvar: "`) {
		result := serverCvarSetRegex.FindStringSubmatch(logEntry.Message)

		if result != nil {
			updatedCvar(g, result[1], result[2])
		}
	}
}

func updatedCvar(g *gameState, name, value string) {
	if _, found := g.cvars[name]; found {
		g.cvars[name] = value
	}
}

func watchCvar(g *gameState, names ...string) {
	for _, name := range names {
		if _, found := g.cvars[name]; !found {
			g.cvars[name] = ""
		}
	}
}
