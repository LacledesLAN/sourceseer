package csgo

import (
	"regexp"
	"strings"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

var (
	regexBetweenQuotes = regexp.MustCompile("(?:\")([0-9A-Za-z_]*)(\")")
)

// UpdateFromStdIn updates the gameState from the processe's standard in
func (g *CSGO) updateFromStdIn(logEntry srcds.LogEntry) {

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

		if logEntry.Message == `World triggered "Game_Commencing"` {

		}

		// set up

		if strings.HasPrefix(logEntry.Message, `World triggered "Match_Start"`) {
			mapName := worldTriggeredMatchStartRegex.FindStringSubmatch(logEntry.Message)[1]

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
}
