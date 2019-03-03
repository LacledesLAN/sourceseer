package csgo

import (
	"math"
	"strings"
)

type teamState struct {
	players    Players
	name       string
	roundsLost int
	roundsWon  int
}

func (m *teamState) HasPlayer(player Player) bool {
	return m.players.HasPlayer(player)
}

// ClientCount returns the number of known clients
func (m *teamState) PlayerCount() uint8 {
	c := len(m.players)

	if c > math.MaxUint8 {
		return math.MaxUint8
	}

	return uint8(c)
}

func (m *teamState) PlayerDropped(player Player) {
	m.players.PlayerDropped(player)
}

func (m *teamState) PlayerJoined(player Player) {
	m.players.PlayerJoined(player)
}

func (m *teamState) ResetStats() {
	m.roundsLost = 0
	m.roundsWon = 0
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

func (m *teamState) String() string {
	return m.name
}
