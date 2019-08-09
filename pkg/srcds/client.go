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

// ClientUnidentifiable occurs when there is not enough information to properly identify a client
func ClientUnidentifiable(c Client) bool {
	if c.IsBot() {
		return len(strings.TrimSpace(c.Username)) == 0
	}

	return len(c.SteamID) == 0
}

func RefreshClient(c Client, cs ...[]Clients) {
	// TODO
}

// EnableFlag enables the specified flag for the client
func (c *Client) EnableFlag(f ClientFlag) {
	c.flags = (c.flags | f)
}

// HasFlag determines if a client has the specified flag enabled
func (c Client) HasFlag(f ClientFlag) bool {
	return (c.flags & f) != 0
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
func (cs Clients) clientIndex(client Client) int {
	for i := range cs {
		if ClientsAreEquivalent(cs[i], client) {
			return i
		}
	}

	return -1
}

// ClientDropped handles when a client drops from srcds
func (cs *Clients) ClientDropped(client Client) {
	i := cs.clientIndex(client)

	if i >= 0 {
		l := len(*cs)

		if l > 1 {
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
		} else if l == 1 {
			*cs = Clients{}
		}
	}
}

// ClientJoined handles when a client connects to srcds
func (cs *Clients) ClientJoined(client Client) {
	if !cs.HasClient(client) {
		*cs = append(*cs, client)
	}
}

// HasClient determines when a client exists
func (cs Clients) HasClient(client Client) bool {
	return cs.clientIndex(client) > -1
}

// RemoveFlags removes the specified flags from all Clients
func (cs *Clients) RemoveFlags(f ClientFlag, fs ...ClientFlag) {
	for _, c := range *cs {
		c.RemoveFlag(f)
		for _, ff := range fs {
			c.RemoveFlag(ff)
		}
	}
}

// WithFlags returns all Clients that have the specified flags
func (cs Clients) WithFlags(f ClientFlag, fs ...ClientFlag) []Client {
	r := []Client{}

iterateClients:
	for _, c := range cs {
		if c.HasFlag(f) {
			for _, ff := range fs {
				if !c.HasFlag(ff) {
					continue iterateClients
				}
			}
			r = append(r, c)
		}
	}

	return r
}

// WithoutFlags returns all Clients that do not have the specified flags
func (cs Clients) WithoutFlags(f ClientFlag, fs ...ClientFlag) []Client {
	r := []Client{}

iterateClients:
	for _, c := range cs {
		if !c.HasFlag(f) {
			for _, ff := range fs {
				if c.HasFlag(ff) {
					continue iterateClients
				}
			}
			r = append(r, c)
		}
	}

	return r
}
