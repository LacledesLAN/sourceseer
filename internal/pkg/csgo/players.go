package csgo

import "github.com/lacledeslan/sourceseer/internal/pkg/srcds"

// Players is a collection of type Player
type Players []Player

// HasPlayer returns true if the player collection contains the specified player
func (m Players) HasPlayer(player Player) bool {
	return m.playerIndex(player) > -1
}

// PlayerDropped from the game; remove them from the player collection
func (m *Players) PlayerDropped(player Player) {
	i := m.playerIndex(player)

	if i >= 0 {
		l := len(*m)

		if l > 1 {
			*m = append((*m)[:i], (*m)[i+1:]...)
		} else if l == 1 {
			*m = Players{}
		}
	}
}

func (m Players) playerIndex(player Player) int {

	for i := range m {
		if srcds.ClientsAreEquivalent(&m[i].Client, &player.Client) {
			return i
		}
	}

	return -1
}

// PlayerJoined the game on this team; add them to the player collection
func (m *Players) PlayerJoined(player Player) {
	if !m.HasPlayer(player) {
		*m = append(*m, player)
	}
}
