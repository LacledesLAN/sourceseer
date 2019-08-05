package srcds

import (
	"bufio"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Cvar represents a watched SRCDS console variable
type Cvar struct {
	LastUpdated time.Time
	Value       string
}

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
	for _, name := range names {
		name = strings.TrimSpace(name)

		if len(name) == 0 {
			continue
		}

		if _, found := o.cvars[name]; !found {
			o.cvars[name] = Cvar{}
		}
	}
}

// AddCvarWatcherDefault instructs the system to keep track of the specified cvar,  providing a default value
func (o *observer) AddCvarWatcherDefault(name string, defaultValue string) {
	name = strings.TrimSpace(name)

	if len(name) == 0 {
		return
	}

	defaultValue = strings.TrimSpace(defaultValue)

	if cvar, found := o.cvars[name]; !found {
		o.cvars[name] = Cvar{Value: defaultValue}
	} else if cvar.LastUpdated.IsZero() {
		o.cvars[name] = Cvar{Value: defaultValue}
	}
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
	cvar, found := o.cvars[name]

	if !found {
		return fallback, false
	}

	i, err := strconv.Atoi(cvar.Value)

	if err != nil {
		return fallback, false
	}

	return i, true
}

// TryCvarAsString attempts to return a cvar as an integer, returning a bool indicating if the provided fallback value was returned
func (o *observer) TryCvarAsString(name, fallback string) (value string, nonFallback bool) {
	cvar, found := o.cvars[name]

	if !found {
		return fallback, false
	}

	return cvar.Value, true
}

type observer struct {
	cvars        map[string]Cvar
	endOfLine    string
	start        func() <-chan LogEntry
	started      time.Time
	testingFlags struct {
		watchAllCvars bool
	}
}

func newObserver() *observer {
	r := &observer{
		cvars:     make(map[string]Cvar),
		endOfLine: "\n",
		start: func() <-chan LogEntry {
			panic("srcds > observer > start function was not instantiated.")
		},
	}

	//TODO: VERIFY EOL IN CASE RUNING LINUX DOCKER IMAGE ON WINDOWS
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
			o.setCvarValue(cvarSet.Name, cvarSet.Value)
			Log(SRCDSCvar, le.Message)
			return
		}

		if cvarSet, ok := paresEchoServerCvar(le.Message); ok {
			o.setCvarValue(cvarSet.Name, cvarSet.Value)
			Log(SRCDSCvar, le.Message)
			return
		}

		if outEntries != nil {
			outEntries <- le
		}

		Log(SRCDSLog, le.Message)
		return
	}

	if cvarSet, ok := parsEchoCvar(line); ok {
		o.setCvarValue(cvarSet.Name, cvarSet.Value)
		Log(SRCDSCvar, line)
		return
	}

	Log(SRCDSOther, line)
}

func (o *observer) setCvarValue(name, value string) {
	if _, found := o.cvars[name]; found || o.testingFlags.watchAllCvars {
		o.cvars[name] = Cvar{
			LastUpdated: time.Now(),
			Value:       strings.TrimSpace(value),
		}
	}
}
