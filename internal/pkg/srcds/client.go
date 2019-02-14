package srcds

import (
	"errors"
	"strings"
)

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

// ExtractClient extract the srcds client from a string
func ExtractClient(s string) (Client, error) {
	c := extractClientRegex.FindStringSubmatch(s)

	if len(c) != 5 {
		return Client{}, errors.New("Unable to parse: " + s)
	}

	return Client{
		Username:   c[1],
		ServerSlot: c[2],
		SteamID:    c[3],
		ServerTeam: c[4],
	}, nil
}

// ExtractClients extracts the players from a srcds log message
func ExtractClients(logEntry LogEntry) (originator, target *Client) {
	originator = nil
	target = nil

	players := extractClientRegex.FindAllStringSubmatch(logEntry.Message, -1)

	if len(players) >= 1 {
		originatorRaw := players[0]
		originator = &Client{Username: originatorRaw[1], ServerSlot: originatorRaw[2], ServerTeam: originatorRaw[4], SteamID: originatorRaw[3]}
	}

	if len(players) >= 2 {
		targetRaw := players[1]
		target = &Client{Username: targetRaw[1], ServerSlot: targetRaw[2], ServerTeam: targetRaw[4], SteamID: targetRaw[3]}
	}

	return
}

// IsBot determines if the client is a bot
func IsBot(m *Client) bool {
	// TODO: UNIT TEST
	if strings.ToUpper(m.SteamID) == "BOT" {
		return true
	}

	return false
}
