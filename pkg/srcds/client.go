package srcds

import (
	"strings"
)

// Client represents a client connected to (or simulated by) the srcds
type Client struct {
	Username    string
	SteamID     string
	ServerSlot  string
	Affiliation string
	Flags       map[string]string
}

// Clients is a collection of type Client
type Clients []Client

// ClientsAreEquivalent determines if two srcds clients are effectively the same client
func ClientsAreEquivalent(c0, c1 Client) bool {
	if len(c0.Username) == 0 || len(c1.Username) == 0 {
		return false
	}

	if c0.IsBot() && c1.IsBot() {
		return c0.Username == c1.Username && c0.ServerSlot == c1.ServerSlot
	}

	if len(c0.SteamID) > 0 && len(c1.SteamID) > 0 {
		return c0.SteamID == c1.SteamID
	}

	// TODO: does server slot ever change?
	return c0.Username == c1.Username && c0.ServerSlot == c1.ServerSlot
}

// ClientDropped handles when a client drops from srcds
func (m *Clients) ClientDropped(client Client) {
	i := m.clientIndex(client)

	if i >= 0 {
		l := len(*m)

		if l > 1 {
			*m = append((*m)[:i], (*m)[i+1:]...)
		} else if l == 1 {
			*m = Clients{}
		}
	}
}

// ClientJoined handles when a client connects to srcds
func (m *Clients) ClientJoined(client Client) {
	if !m.HasClient(client) {
		*m = append(*m, client)
	}
}

// HasClient determines when a client exists
func (m Clients) HasClient(client Client) bool {
	return m.clientIndex(client) > -1
}

// IsBot determines if the client is a bot
func (m *Client) IsBot() bool {
	if strings.ToUpper(m.SteamID) == "BOT" {
		return true
	}

	return false
}

// clientIndex gets the index of the provided client in the slice of clients
func (m Clients) clientIndex(client Client) int {
	for i := range m {
		if ClientsAreEquivalent(m[i], client) {
			return i
		}
	}

	return -1
}
