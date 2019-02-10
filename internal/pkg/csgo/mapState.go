package csgo

import (
	"strconv"
	"time"
)

type mapState struct {
	name           string
	mpTeam1        teamState // Start as CT
	mpTeam2        teamState // Start as Terrorists
	started        time.Time
	ended          time.Time
	isSwappedSides bool
	mapStarted     time.Time
}

func (m *mapState) ct() *teamState {
	if m.isSwappedSides {
		return &m.mpTeam2
	}

	return &m.mpTeam1
}

func (m *mapState) CTSetScore(score string) {
	s, _ := strconv.Atoi(score)
	ct := m.ct()
	ct.roundsWon = s

	t := m.terrorist()
	t.roundsLost = s
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

func (m *mapState) RoundsCompleted() int {
	return m.ct().roundsWon + m.terrorist().roundsWon
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

func (m *mapState) TerroristSetScore(score string) {
	s, _ := strconv.Atoi(score)
	ct := m.ct()
	ct.roundsLost = s

	t := m.terrorist()
	t.roundsWon = s
}
