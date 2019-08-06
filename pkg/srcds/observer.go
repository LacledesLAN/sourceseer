package srcds

import (
	"bufio"
	"runtime"
	"strings"
	"time"
)

// Observer for watching SRCDS log streams
type Observer interface {
	AddCvarWatcher(names ...string)
	AddCvarWatcherDefault(name string, defaultValue string)
	Start() <-chan LogEntry
	TryCvarAsInt(name string, fallback int) (value int, nonFallback bool)
	TryCvarAsString(name, fallback string) (value string, nonFallback bool)
}

// AddCvarWatcher instructs the system to keep track of the specified cvar name
func (o *observer) AddCvarWatcher(names ...string) {
	o.cvars.addWatcher(names...)
}

// AddCvarWatcherDefault instructs the system to keep track of the specified cvar,  providing a default value
func (o *observer) AddCvarWatcherDefault(name string, defaultValue string) {
	o.cvars.seedWatcher(name, defaultValue)
}

// NewScanner for observing streaming SRCDS data
func NewScanner(scanner bufio.Scanner) Observer {
	o := newObserver()

	o.start = func() <-chan LogEntry {
		runtime.GC()
		logStream := make(chan LogEntry, 6)

		go func(c chan<- LogEntry) {
			defer close(c)
			for scanner.Scan() {
				o.processMessage(scanner.Text(), c)
			}
		}(logStream)

		return logStream
	}

	return o
}

// Start the SRCDS process
func (o *observer) Start() <-chan LogEntry {
	return o.start()
}

// TryCvarAsInt attempts to return a cvar as an integer, returning a bool indicating if the provided fallback value was returned
func (o *observer) TryCvarAsInt(name string, fallback int) (value int, nonFallback bool) {
	return o.cvars.tryInt(name, fallback)
}

// TryCvarAsString attempts to return a cvar as an integer, returning a bool indicating if the provided fallback value was returned
func (o *observer) TryCvarAsString(name, fallback string) (value string, nonFallback bool) {
	return o.cvars.tryString(name, fallback)
}

type observer struct {
	cvars        Cvars
	endOfLine    string
	start        func() <-chan LogEntry
	started      time.Time
	testingFlags struct {
	}
}

func newObserver() *observer {
	r := &observer{
		endOfLine: "\n",
		start: func() <-chan LogEntry {
			panic("srcds > observer > start function was not instantiated.")
		},
	}

	//TODO: VERIFY EOL IN CASE RUNNING LINUX DOCKER IMAGE ON WINDOWS
	//if runtime.GOOS == "windows" {
	//	r.endOfLine = "\r\n"
	//}

	return r
}

func (o *observer) processMessage(line string, outEntries chan<- LogEntry) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return
	}

	if le, ok := parseLogEntry(line); ok {
		if cvarSet, ok := parsEchoCvar(le.Message); ok {
			o.cvars.setIfWatched(cvarSet.Name, cvarSet.Value, le.Timestamp)
			return
		}

		if cvarSet, ok := paresEchoServerCvar(le.Message); ok {
			o.cvars.setIfWatched(cvarSet.Name, cvarSet.Value, le.Timestamp)
			return
		}

		if outEntries != nil {
			outEntries <- le
		}

		Log(SRCDSLog, le.Message)
		return
	}

	if cvarSet, ok := parsEchoCvar(line); ok {
		o.cvars.setIfWatched(cvarSet.Name, cvarSet.Value, time.Now())
		Log(SRCDSCvar, line)
		return
	}

	Log(SRCDSOther, line)
}
