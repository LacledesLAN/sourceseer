package srcds

import (
	"strings"
	"time"
)

// LogEntry represents a log message from srcds
type LogEntry struct {
	Message   string
	Raw       string
	Timestamp time.Time
}

// ExtractLogEntry extracts the log entry from a srcds log message
func ExtractLogEntry(rawLogEntry string) LogEntry {
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
