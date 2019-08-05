package srcds

import (
	"regexp"
	"strings"
	"time"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var clientRegex = regexp.MustCompile(`^"(.{1,32})<(\d{0,2})><\[?([\w:]*)\]?>(?:<([\w]*?)>)?"`)

// ParseClient attempts to parse a srcds client
func ParseClient(s string) (Client, bool) {
	tokens := clientRegex.FindStringSubmatch(s)

	if len(tokens) != 5 {
		return Client{}, false
	}

	return Client{
		Affiliation: tokens[4],
		SteamID:     tokens[3],
		ServerSlot:  tokens[2],
		Username:    tokens[1],
	}, true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var clientLogEntryRegex = regexp.MustCompile(`^"(.{1,32})<(\d{0,2})><\[?([\w:]*)\]?>(?:<([\w]*?)>)?" (.+)$`)

// ClientLogEntry is sent when a log entry is caused by a client action
type ClientLogEntry struct {
	Client  Client
	Message string
}

// ParseClientLogEntry determines if a client action caused a log entry
func ParseClientLogEntry(le LogEntry) (cle ClientLogEntry, ok bool) {
	tokens := clientLogEntryRegex.FindStringSubmatch(le.Message)

	if len(tokens) != 6 {
		return ClientLogEntry{}, false
	}

	return ClientLogEntry{
		Message: strings.TrimSpace(tokens[5]),
		Client: Client{
			Affiliation: tokens[4],
			SteamID:     tokens[3],
			ServerSlot:  tokens[2],
			Username:    tokens[1],
		},
	}, true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var clientConnectedRegex = regexp.MustCompile(`^connected, address "[\S]*?"$`)

// ParseClientConnected - determine if a client connected
func ParseClientConnected(clientMsg ClientLogEntry) (ok bool) {
	if clientConnectedRegex.FindIndex([]byte(clientMsg.Message)) != nil {
		return true
	}

	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var clientDisconnectedRegex = regexp.MustCompile(`^disconnected(?: \(reason \"([\w ]{1,})\"\))?$`)

// ClientDisconnectedReason is sent when a client disconnects from srcds
type ClientDisconnectedReason string

// ParseClientDisconnected - determine if a client disconnected
func ParseClientDisconnected(clientMsg ClientLogEntry) (ClientDisconnectedReason, bool) {
	tokens := clientDisconnectedRegex.FindStringSubmatch(clientMsg.Message)

	if len(tokens) != 2 {
		return ClientDisconnectedReason(""), false
	}

	return ClientDisconnectedReason(tokens[1]), true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CvarValueSet is sent when srcds outputs a cvar
type CvarValueSet struct {
	Name  string
	Value string
}

var (
	echoCvarFmt0 = regexp.MustCompile(`^"([^\s]{3,})" = "([^\s"]{0,})"`)
	echoCvarFmt1 = regexp.MustCompile(`^([\w]{3,}) -(?:| (-?[\w.;]{0,}))$`)
)

func parsEchoCvar(s string) (CvarValueSet, bool) {
	tokens := echoCvarFmt0.FindStringSubmatch(s)
	if len(tokens) != 3 {
		tokens = echoCvarFmt1.FindStringSubmatch(s)
	}

	if len(tokens) == 3 {
		return CvarValueSet{
			Value: tokens[2],
			Name:  tokens[1],
		}, true
	}

	return CvarValueSet{}, false
}

var (
	echoServerCvar = regexp.MustCompile(`^server_cvar: "([\w]{3,})" "(.*)"$`)
)

func paresEchoServerCvar(s string) (CvarValueSet, bool) {
	tokens := echoServerCvar.FindStringSubmatch(s)

	if len(tokens) != 3 {
		return CvarValueSet{}, false
	}

	return CvarValueSet{
		Value: tokens[2],
		Name:  tokens[1],
	}, true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const srcdsTimeLayout string = "1/2/2006 - 15:04:05"

// LogEntry is sent when SRCDS sends a log entry
type LogEntry struct {
	Message   string
	Timestamp time.Time
}

func parseLogEntry(s string) (le LogEntry, ok bool) {
	s = strings.TrimSpace(s)

	if len(s) > 0 && strings.HasPrefix(s, "L ") && strings.Contains(s, ": ") {
		i := strings.Index(s, ": ")

		if msgTime, err := time.ParseInLocation(srcdsTimeLayout, s[2:i], time.Local); err == nil {
			return LogEntry{
				Timestamp: msgTime,
				Message:   strings.TrimSpace(s[i+2:]),
			}, true
		}
	}

	return LogEntry{}, false
}
