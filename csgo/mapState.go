package csgo

import (
	"time"
)

type MapMode uint8

const (
	ModeUnknown MapMode = 0
	ModeWarmUp  MapMode = iota
	ModePlay
	ModeOvertime
)

type mapState struct {
	name           string
	mpTeam1        teamState // Start as CT
	mpTeam2        teamState // Start as Terrorists
	started        time.Time
	ended          time.Time
	isSwappedSides bool
	roundNumber    byte
	mapStarted     time.Time
	mode           MapMode
}

func (m *mapState) ct() *teamState {
	if m.isSwappedSides {
		return &m.mpTeam2
	}

	return &m.mpTeam1
}

func (m *mapState) CTWonRound() {
	m.roundNumber = m.roundNumber + 1

	ct := m.ct()
	ct.roundsWon = ct.roundsWon + 1

	t := m.terrorist()
	t.roundsLost = t.roundsLost + 1
}

func (m *mapState) PlayerDropped(player Player) {
	m.mpTeam1.PlayerDropped(player)
	m.mpTeam2.PlayerDropped(player)
}

func (m *mapState) PlayerJoinedCT(player Player) {
	t := m.terrorist()
	t.PlayerDropped(player)

	ct := m.ct()
	ct.PlayerJoined(player)
}

func (m *mapState) PlayerJoinedTerrorist(player Player) {
	ct := m.ct()
	ct.PlayerDropped(player)

	t := m.terrorist()
	t.PlayerJoined(player)
}

func (m *mapState) TeamsSwappedSides() {
	m.isSwappedSides = !m.isSwappedSides
}

func (m *mapState) terrorist() *teamState {
	if m.isSwappedSides {
		return &m.mpTeam1
	}

	return &m.mpTeam2
}

func (m *mapState) TerroristWonRound() {
	m.roundNumber = m.roundNumber + 1

	ct := m.ct()
	ct.roundsLost = ct.roundsLost + 1

	t := m.terrorist()
	t.roundsWon = t.roundsWon + 1
}
