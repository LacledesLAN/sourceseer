package srcds

import (
	"regexp"
	"strconv"
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

	r := Client{
		Affiliation: tokens[4],
		SteamID:     tokens[3],
		Username:    tokens[1],
	}

	if i, err := strconv.ParseInt(tokens[2], 10, 16); err != nil {
		r.ServerSlot = -1
	} else {
		r.ServerSlot = int16(i)
	}

	return r, true
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

	r := ClientLogEntry{
		Message: strings.TrimSpace(tokens[5]),
		Client: Client{
			Affiliation: tokens[4],
			SteamID:     tokens[3],
			Username:    tokens[1],
		},
	}

	if i, err := strconv.ParseInt(tokens[2], 10, 16); err != nil {
		r.Client.ServerSlot = -1
	} else {
		r.Client.ServerSlot = int16(i)
	}

	return r, true
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
var clientDisconnectedRegex = regexp.MustCompile(`^disconnected(?: \(reason \"([\S ]{1,})\"\))?$`)

// ClientDisconnectedReason is sent when a client disconnects from srcds
type ClientDisconnectedReason string

// ParseClientDisconnected - determine if a client disconnected
func ParseClientDisconnected(clientMsg ClientLogEntry) (ClientDisconnectedReason, bool) {
	tokens := clientDisconnectedRegex.FindStringSubmatch(clientMsg.Message)

	if len(tokens) != 2 {
		return ClientDisconnectedReason(""), false
	}

	reason := tokens[1]

	if strings.HasSuffix(reason, "timed out") {
		return ClientDisconnectedReason("timed out"), true
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
	cvarDescription  = regexp.MustCompile(`^([\w]{3,}) -(?:| (-?[\w.;]{0,}))$`)
	cvarQuotedValue  = regexp.MustCompile(`^"([^\s]{3,})" = "([^\s"]{0,})"`)
	cvarServerFormat = regexp.MustCompile(`^server_cvar: "([\w]{3,})" "(.*)"$`)
)

// parseCvar from logging output
func parseCvar(le LogEntry) (cvar CvarValueSet, ok bool) {
	tokens := cvarServerFormat.FindStringSubmatch(le.Message)
	if len(tokens) != 3 {
		tokens = cvarQuotedValue.FindStringSubmatch(le.Message)
	}

	if len(tokens) == 3 {
		return CvarValueSet{
			Value: tokens[2],
			Name:  tokens[1],
		}, true
	}

	return CvarValueSet{}, false
}

// parseCvarResponse is sent via standard out when a cvar name is passed to SRCDS without an associated value
func parseCvarResponse(s string) (cvar CvarValueSet, ok bool) {
	tokens := cvarDescription.FindStringSubmatch(s)
	if len(tokens) != 3 {
		tokens = cvarQuotedValue.FindStringSubmatch(s)
	}

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
