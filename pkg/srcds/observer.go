package srcds

import (
	"bufio"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	eolUnix    = "\n"
	eolWindows = "\r\n"
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

func newReader(byteStream io.Reader) *observer {
	o := newObserver()

	o.start = func() <-chan LogEntry {
		logStream := make(chan LogEntry, 6)
		br := bufio.NewReader(byteStream)
		runtime.GC()

		go func(c chan<- LogEntry) {
			defer close(c)
			line, err := br.ReadString(0x0A)

			if err == io.EOF || len(line) == 0 {
				return
			}

			// Determine EOL delimiter as it may not match operating system's EOL
			if strings.HasSuffix(line, eolWindows) {
				o.endOfLine = eolWindows
				line = line[:len(line)-2]
				log.Debug().Msg("Windows EOL delimiter detected.")
			} else {
				o.endOfLine = eolUnix

				if strings.HasSuffix(line, eolUnix) {
					log.Debug().Msg("Unix EOL delimiter detected.")
					line = line[:len(line)-1]
				} else {
					log.Warn().Msg("Couldn't detect EOL delimiter; defaulting to Unix EOL.")
					line = strings.TrimSpace(line)
				}
			}

			o.processMessage(line, c)

			scanner := bufio.NewScanner(br)
			for scanner.Scan() {
				o.processMessage(scanner.Text(), c)
			}
		}(logStream)

		return logStream
	}

	return o
}

// NewReader for processing streaming SRCDS log data
func NewReader(byteStream io.Reader) Observer {
	return newReader(byteStream)
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
	cvars      Cvars
	endOfLine  string
	start      func() <-chan LogEntry
	started    time.Time
	statistics observerStatistics
}

type observerStatistics struct {
	totalLines uint32
	blankLines uint32
	logLines   uint32
}

func newObserver() *observer {
	r := &observer{
		start: func() <-chan LogEntry {
			panic("srcds > observer > start function was not instantiated.")
		},
	}

	if strings.ToLower(runtime.GOOS) == "windows" {
		r.endOfLine = eolWindows
	} else {
		r.endOfLine = eolUnix
	}

	return r
}

func (o *observer) processMessage(line string, outEntries chan<- LogEntry) {
	o.statistics.totalLines++

	line = strings.TrimSpace(line)
	if len(line) == 0 {
		o.statistics.blankLines++
		return
	}

	if le, ok := parseLogEntry(line); ok {
		o.statistics.logLines++

		if cvarSet, ok := parseCvar(le); ok {
			o.cvars.setIfWatched(cvarSet.Name, cvarSet.Value, le.Timestamp)
			return
		}

		log.Debug().Str("source", "log entry").Msg(le.Message)

		if outEntries != nil {
			outEntries <- le
		}

		return
	}

	if cvarSet, ok := parseCvarResponse(line); ok {
		o.cvars.setIfWatched(cvarSet.Name, cvarSet.Value, time.Now())
		log.Debug().Str("source", "stdout").Msg(line)
		return
	}

	log.Debug().Str("source", "stdout").Msg(line)
}
