package srcds

import "strings"

// Client represents a client connected to (or simulated by) the srcds
type Client struct {
	Username   string
	SteamID    string
	ServerSlot string
	ServerTeam string
}

// ClientsAreEquivalent determines if two srcds clients are effectively the same client
func ClientsAreEquivalent(client0, client1 *Client) bool {
	if client0 == nil || client1 == nil {
		return false
	}

	if len(client0.SteamID) > 0 && len(client1.SteamID) > 0 {
		return client0.Username == client1.Username && client0.SteamID == client1.SteamID
	}

	if len(client0.Username) == 0 || len(client1.Username) == 0 {
		return false
	}

	return client0.Username == client1.Username
}

// IsBot determines if the client is a bot
func IsBot(m *Client) bool {
	// TODO: UNIT TEST
	if strings.ToUpper(m.SteamID) == "BOT" {
		return true
	}

	return false
}
