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

type ClientDisconnected struct {
	Client Client
	Reason string
}

type CvarValueSet struct {
	name  string
	value string
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

func parseCvarValueSet(le LogEntry) (CvarValueSet, error) {
	if strings.HasPrefix(le.Message, `"`) {
		r := serverCvarEchoRegex.FindStringSubmatch(le.Message)

		if len(r) == 3 {
			return CvarValueSet{
				name:  r[1],
				value: r[2],
			}, nil
		}
	} else if strings.HasPrefix(le.Message, `server_cvar: "`) {
		r := serverCvarSetRegex.FindStringSubmatch(le.Message)

		if len(r) == 3 {
			return CvarValueSet{
				name:  r[1],
				value: r[2],
			}, nil
		}
	}

	return CvarValueSet{}, errors.New("Could not parse cvar out of :" + le.Message)
}

// parseLogEntry extracts the log entry from a srcds log message
func parseLogEntry(rawLogEntry string) LogEntry {
	rawLogEntry = strings.TrimSpace(rawLogEntry)

	if len(rawLogEntry) > 0 && strings.HasPrefix(rawLogEntry, "L ") && strings.Contains(rawLogEntry, ": ") {
		i := strings.Index(rawLogEntry, ": ")

		msgTime, err := time.ParseInLocation(srcdsTimeLayout, rawLogEntry[2:i], time.Local)

		if err == nil {
			return LogEntry{
				Raw:       rawLogEntry,
				Timestamp: msgTime,
				Message:   rawLogEntry[i+2:],
			}
		}
	}

	return LogEntry{Raw: rawLogEntry}
}
