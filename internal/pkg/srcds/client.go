package srcds

import (
	"errors"
	"regexp"
	"strings"
)

const (
	clientPattern = `^"(.{1,32})<(\d{0,2})><([\w:]*)><{0,1}([a-zA-Z0-9]*?)>{0,1}"`
)

var (
	clientRegex = regexp.MustCompile(clientPattern)
)

// Client represents a client connected to (or simulated by) the srcds
type Client struct {
	Username   string
	SteamID    string
	ServerSlot string
	ServerTeam string
}

//Clients is a collection of type Client
type Clients []Client

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

//ClientDropped handles when a client drops from srcds
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

func (m Clients) clientIndex(client Client) int {

	for i := range m {
		if ClientsAreEquivalent(&m[i], &client) {
			return i
		}
	}

	return -1
}

//ClientJoined handles when a client connects to srcds
func (m *Clients) ClientJoined(client Client) {
	if !m.HasClient(client) {
		*m = append(*m, client)
	}
}

//HasClient determines when a client exists
func (m Clients) HasClient(client Client) bool {
	return m.clientIndex(client) > -1
}

// ParseClient attempts to parse a srcds client
func ParseClient(s string) (Client, error) {
	c := clientRegex.FindStringSubmatch(s)

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
