package srcds

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

const (
	clientConnectedPattern    = `^(".*") (?:entered the game|connected, address "")$`
	clientDisconnectedPattern = `^(".*") disconnected(?: \(reason \"([\w ]{1,})\"\))?$`
	serverCvarEchoPattern     = `^"([^\s]{3,})" = "([^\s"]{0,})"`
	serverCvarSetPattern      = `^server_cvar: "([^\s]{3,})" "(.*)"$`
)

var (
	clientConnectedRegex    = regexp.MustCompile(clientConnectedPattern)
	clientDisconnectedRegex = regexp.MustCompile(clientDisconnectedPattern)
	serverCvarEchoRegex     = regexp.MustCompile(serverCvarEchoPattern)
	serverCvarSetRegex      = regexp.MustCompile(serverCvarSetPattern)
)

//ClientDisconnected is sent when a client disconnects from srcds
type ClientDisconnected struct {
	Client Client
	Reason string
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

func parseClientConnected(le LogEntry) (Client, error) {
	result := clientConnectedRegex.FindStringSubmatch(le.Message)

	if len(result) != 2 {
		return Client{}, errors.New("Could not client " + le.Message)
	}

	cl, err := ParseClient(result[1])

	if err != nil {
		return Client{}, errors.New("Could not parse client in " + le.Message)
	}

	return cl, nil
}

func parseClientDisconnected(le LogEntry) (ClientDisconnected, error) {
	r := clientDisconnectedRegex.FindStringSubmatch(le.Message)

	if len(r) < 1 {
		return ClientDisconnected{}, errors.New("Could not parse " + le.Message)
	}

	cl, err := ParseClient(r[1])

	if err != nil {
		return ClientDisconnected{}, errors.New("Could not parse player in " + le.Message)
	}

	return ClientDisconnected{
		Client: cl,
		Reason: r[2],
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
