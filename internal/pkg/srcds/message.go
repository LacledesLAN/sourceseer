package srcds

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

const (
	clientConnectedPattern    = `^(?:entered the game|connected, address "[\S]*?")$`
	clientDisconnectedPattern = `^disconnected(?: \(reason \"([\w ]{1,})\"\))?$`
	clientMessagePattern      = `^"(.{1,32})<(\d{0,2})><\[?([\w:]*)\]?>(?:<([\w]*?)>)?" (.+)$`
	serverCvarEchoPattern     = `^"([^\s]{3,})" = "([^\s"]{0,})"`
	serverCvarSetPattern      = `^server_cvar: "([^\s]{3,})" "(.*)"$`
)

var (
	clientConnectedRegex    = regexp.MustCompile(clientConnectedPattern)
	clientDisconnectedRegex = regexp.MustCompile(clientDisconnectedPattern)
	clientMessageRegex      = regexp.MustCompile(clientMessagePattern)
	serverCvarEchoRegex     = regexp.MustCompile(serverCvarEchoPattern)
	serverCvarSetRegex      = regexp.MustCompile(serverCvarSetPattern)
)

//ClientDisconnected is sent when a client disconnects from srcds
type ClientDisconnected struct {
	Client Client
	Reason string
}

type ClientMessage struct {
	Client  Client
	Message string
}

//Cvar represents a SRCDS cvar
type Cvar struct {
	LastUpdate time.Time
	Value      string
}

//CvarValueSet is sent when srcds outputs a cvar
type CvarValueSet struct {
	Name  string
	Value string
}

// LogEntry represents a log message from srcds
type LogEntry struct {
	Message   string
	Raw       string
	Timestamp time.Time
}

func parseClientConnected(clientMsg ClientMessage) (Client, error) {
	if clientConnectedRegex.MatchString(clientMsg.Message) {
		return clientMsg.Client, nil
	}

	return Client{}, errors.New("Could not parse client connect string from " + clientMsg.Message)
}

func parseClientDisconnected(clientMsg ClientMessage) (ClientDisconnected, error) {
	r := clientDisconnectedRegex.FindStringSubmatch(clientMsg.Message)

	if len(r) < 1 {
		return ClientDisconnected{}, errors.New("Could not parse client disconnected string from " + clientMsg.Message)
	}

	return ClientDisconnected{
		Client: clientMsg.Client,
		Reason: r[1],
	}, nil
}

func parseClientMessage(le LogEntry) (ClientMessage, error) {
	r := clientMessageRegex.FindStringSubmatch(le.Message)

	if len(r) < 1 {
		return ClientMessage{}, errors.New("Could not parse client message from string: " + le.Message)
	}

	return ClientMessage{
		Client: Client{
			Username:    r[1],
			ServerSlot:  r[2],
			SteamID:     r[3],
			Affiliation: r[4],
		},
		Message: r[5],
	}, nil
}

//ParseCvarValueSet parse when srcds outputs a cvar
func ParseCvarValueSet(s string) (CvarValueSet, error) {
	if strings.HasPrefix(s, `"`) {
		r := serverCvarEchoRegex.FindStringSubmatch(s)

		if len(r) == 3 {
			return CvarValueSet{
				Name:  r[1],
				Value: r[2],
			}, nil
		}
	} else if strings.HasPrefix(s, `server_cvar: "`) {
		r := serverCvarSetRegex.FindStringSubmatch(s)

		if len(r) == 3 {
			return CvarValueSet{
				Name:  r[1],
				Value: r[2],
			}, nil
		}
	}

	return CvarValueSet{}, errors.New("Could not parse cvar out of string '" + s + "'")
}

// parseLogEntry extracts the log entry from a srcds log message
func parseLogEntry(s string) (LogEntry, error) {
	s = strings.TrimSpace(s)

	if len(s) > 0 && strings.HasPrefix(s, "L ") && strings.Contains(s, ": ") {
		i := strings.Index(s, ": ")

		msgTime, err := time.ParseInLocation(srcdsTimeLayout, s[2:i], time.Local)

		if err == nil {
			return LogEntry{
				Raw:       s,
				Timestamp: msgTime,
				Message:   s[i+2:],
			}, nil
		}
	}

	return LogEntry{}, errors.New("Unable to parse log entry from: " + s)
}
