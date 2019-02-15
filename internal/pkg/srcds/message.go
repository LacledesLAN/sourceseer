package srcds

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

const (
	serverCvarEchoPattern = `^"([^\s]{3,})" = "([^\s"]{0,})"`
	serverCvarSetPattern  = `^server_cvar: "([^\s]{3,})" "(.*)"$`
)

var (
	serverCvarEchoRegex = regexp.MustCompile(serverCvarEchoPattern)
	serverCvarSetRegex  = regexp.MustCompile(serverCvarSetPattern)
)

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
