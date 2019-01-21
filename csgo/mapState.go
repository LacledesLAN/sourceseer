package csgo

import (
	"time"

	"github.com/lacledeslan/sourceseer/srcds"
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
}

func (m *mapState) ClientDrop(client srcds.Client) {
	ct := m.getCT()
	t := m.getT()

	ct.PlayerRemove(client)
	t.PlayerRemove(client)
}

func (m *mapState) CTPlayerJoin(client srcds.Client) {
	t := m.getT()
	t.PlayerRemove(client)

	ct := m.getCT()
	ct.PlayerJoin(client)
}

func (m *mapState) CTWinRound() {
	ct := m.getCT()
	ct.roundsWon = ct.roundsWon + 1
	m.roundNumber = m.roundNumber + 1
}

func (m *mapState) SwapSides() {
	m.isSwappedSides = !m.isSwappedSides
}

func (m *mapState) TerroristPlayerJoin(player srcds.Client) {
	ct := m.getCT()
	ct.PlayerRemove(player)

	t := m.getT()
	t.PlayerJoin(player)
}

func (m *mapState) TerroristWinRound() {
	t := m.getT()
	t.roundsWon = t.roundsWon + 1
	m.roundNumber = m.roundNumber + 1
}

func (m *mapState) getCT() *teamState {
	if m.isSwappedSides {
		return &m.mpTeam2
	}

	return &m.mpTeam1
}

func (m *mapState) getT() *teamState {
	if m.isSwappedSides {
		return &m.mpTeam1
	}

	return &m.mpTeam2
}
