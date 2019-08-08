package srcds

import (
	"strings"
)

// Dev notes: ServerSlot and Affiliation can change at any given moment!

// Client represents a client connected to (or simulated by) the srcds
type Client struct {
	Username    string
	SteamID     string
	ServerSlot  string
	Affiliation string
	flags       ClientFlag
}

// ClientFlag represents a flag set for a client
type ClientFlag uint16

// ClientUnidentifiable occurs when there is not enough information to properly identify a client
func ClientUnidentifiable(c Client) bool {
	if c.IsBot() {
		return len(strings.TrimSpace(c.Username)) == 0
	}

	return len(c.SteamID) == 0
}

// ClientsAreEquivalent determines if two srcds clients are effectively the same client
func ClientsAreEquivalent(c0, c1 Client) bool {
	if ClientUnidentifiable(c0) || ClientUnidentifiable(c1) {
		return false
	}

	if c0.IsBot() && c1.IsBot() {
		return c0.Username == c1.Username
	}

	if c0.IsConsole() {
		return c1.IsConsole()
	}

	//TODO: Condition possibly needed, likely based on SteamID field, if the server/clients are unable
	//to connect to steam to get their SteamID

	return c0.SteamID == c1.SteamID
}

// HasFlag determines if a client has the specified flag enabled
func (c Client) HasFlag(f ClientFlag) bool {
	return (c.flags & f) != 0
}

// EnableFlag enables the specified flag for the client
func (c *Client) EnableFlag(f ClientFlag) {
	c.flags = (c.flags | f)
}

// RemoveAllFlags resets all flags for the client
func (c *Client) RemoveAllFlags() {
	c.flags = ClientFlag(0)
}

// RemoveFlag remove the specified flag for the client
func (c *Client) RemoveFlag(f ClientFlag) {
	c.flags = c.flags &^ f
}

// ToggleFlag enables the flag if disabled; disable if it was enabled
func (c *Client) ToggleFlag(f ClientFlag) {
	c.flags = c.flags ^ f
}

// IsBot determines if the client is a bot
func (c Client) IsBot() bool {
	if strings.ToUpper(c.SteamID) == "BOT" {
		return true
	}

	return false
}

// IsConsole determines if the client is just the server console
func (c Client) IsConsole() bool {
	return strings.ToUpper(c.SteamID) == "CONSOLE" && strings.ToUpper(c.Affiliation) == "CONSOLE"
}

// Clients is a collection of individual clients; useful for teams
type Clients []Client

// clientIndex gets the index of the provided client in the slice of clients
func (m Clients) clientIndex(client Client) int {
	for i := range m {
		if ClientsAreEquivalent(m[i], client) {
			return i
		}
	}

	return -1
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
