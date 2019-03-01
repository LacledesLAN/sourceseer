package csgo

import (
	"time"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

type mapState struct {
	name       string
	mpTeam1    teamState // Start as CT
	mpTeam2    teamState // Start as Terrorists
	started    time.Time
	ended      time.Time
	mapStarted time.Time
	spectators srcds.Clients
}

func (m *mapState) PlayerDropped(player Player) {
	m.mpTeam1.PlayerDropped(player)
	m.mpTeam2.PlayerDropped(player)
}

func (m *mapState) RoundsCompleted() int {
	return m.mpTeam1.roundsWon + m.mpTeam2.roundsWon
}
