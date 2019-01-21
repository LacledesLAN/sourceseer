package csgo

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/srcds"
)

type csgoState struct {
	started     time.Time
	maps        []mapState
	currentMap  *mapState
	mpTeamname1 string
	mpTeamname2 string
}

func (m *csgoState) CTPlayerJoin(player srcds.Client) {
	m.currentMap.CTPlayerJoin(player)
}

func (m *csgoState) CTWinRound() {
	m.currentMap.CTWinRound()
}

func (m *csgoState) GetRoundNumber() byte {
	return m.currentMap.roundNumber
}

func (m *csgoState) SwapSides() {
	m.currentMap.SwapSides()
}

func (m *csgoState) TerroristPlayerJoin(player srcds.Client) {
	m.currentMap.TerroristPlayerJoin(player)
}

func (m *csgoState) TerroristWinRound() {
	m.currentMap.TerroristWinRound()
}

func (m *csgoState) ChangeLevel(mapName string) {
	i := len(m.maps)

	if i > 0 {
		m.maps[i-1].ended = time.Now()
	}

	m.maps = append(m.maps, mapState{
		name:    mapName,
		started: time.Now()},
	)

	m.currentMap = &m.maps[i]

	if (len(m.mpTeamname1)) == 0 {
		m.mpTeamname1 = "mp_team_1"
	}
	m.currentMap.mpTeam1.SetName(m.mpTeamname1)

	if (len(m.mpTeamname2)) == 0 {
		m.mpTeamname2 = "mp_team_2"
	}
	m.currentMap.mpTeam2.SetName(m.mpTeamname2)
}

func ProcessLogEntry(m *csgoState, logEntry srcds.LogEntry) {

	if strings.HasPrefix(logEntry.Message, "\"") {
		originator, _ := srcds.ExtractClients(logEntry)

		if originator != nil {
			// a player did something

			if strings.Contains(logEntry.Message, `" switched from team <`) {
				if strings.HasSuffix(logEntry.Message, "<CT>") {
					m.CTPlayerJoin(*originator)
				} else if strings.HasSuffix(logEntry.Message, "<TERRORIST>") {
					m.TerroristPlayerJoin(*originator)
				}
			}
		} else {
			// a variable was assigned
		}

		return
	}

	if strings.HasPrefix(logEntry.Message, "Team") {
		if strings.HasPrefix(logEntry.Message, `Team "CT" scored "`) {
			m.CTWinRound()
		} else if strings.HasPrefix(logEntry.Message, `Team "TERRORIST" scored "`) {
			m.TerroristWinRound()
		} else if strings.HasPrefix(logEntry.Message, `Team playing "CT":`) {

		} else if strings.HasPrefix(logEntry.Message, `Team playing "TERRORIST":`) {

		} else {

		}

		return
	}

	if strings.HasPrefix(logEntry.Message, "World triggered") {
		return
	}

	if strings.HasPrefix(logEntry.Message, "Game Over:") {
		return
	}

	if strings.HasPrefix(logEntry.Message, "Loading map") {
		regexBetweenQuotes := regexp.MustCompile("(?:\")([0-9A-Za-z_]*)(\")")
		mapName := regexBetweenQuotes.FindString(logEntry.Message)
		mapName = mapName[1 : len(mapName)-1]

		return
	}

	fmt.Println("??? - ", logEntry.Message)
}
