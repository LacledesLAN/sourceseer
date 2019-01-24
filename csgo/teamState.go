package csgo

import (
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/srcds"
)

type teamState struct {
	name           string
	pausedTimeUsed time.Duration
	roundsWon      uint8
	knownPlayers   []csgoClient
}

// playerIndex find the index of the player in the knownPlayers pool (-1 if not found)
func (m *teamState) playerIndex(player srcds.Client) int {
	for i := range m.knownPlayers {
		if srcds.ClientsAreEquivalent(&m.knownPlayers[i].Client, &player) {
			return i
		}
	}

	return -1
}

// PlayerJoin adds a player to the team
func (m *teamState) PlayerJoin(player srcds.Client) {
	if m.playerIndex(player) < 0 {
		m.knownPlayers = append(m.knownPlayers, csgoClient{Client: player})
	}
}

// PlayerRemove removes a player from the team
func (m *teamState) PlayerRemove(player srcds.Client) {
	i := m.playerIndex(player)

	if i >= 0 {
		l := len(m.knownPlayers)

		if l > 1 {
			m.knownPlayers = append(m.knownPlayers[:i], m.knownPlayers[i+1:]...)
		} else if l == 1 {
			m.knownPlayers = []csgoClient{}
		}
	}
}

// SetName sets the team's name
func (m *teamState) SetName(teamName string) {
	teamName = strings.TrimSpace(teamName)

	if len(teamName) > 0 {
		m.name = strings.Join(strings.Fields(teamName), "_")
	} else {
		m.name = "Unspecified"
	}
}
