package srcds

import (
	"bufio"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	eolUnix    = "\n"
	eolWindows = "\r\n"
)

// AddCvarWatcher instructs the system to keep track of the specified cvar name
func (o *Observer) AddCvarWatcher(names ...string) {
	o.cvars.addWatcher(names...)
}

// AddCvarWatcherDefault instructs the system to keep track of the specified cvar,  providing a default value
func (o *Observer) AddCvarWatcherDefault(name string, defaultValue string) {
	o.cvars.seedWatcher(name, defaultValue)
}

// NewObserver for SRCDS log streams
func NewObserver() *Observer {
	r := &Observer{}

	if strings.ToLower(runtime.GOOS) == "windows" {
		r.EndOfLine = eolWindows
	} else {
		r.EndOfLine = eolUnix
	}

	return r
}

// Read a SRCDS log output stream
func (o *Observer) Read(r io.Reader) {
	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		for range o.Listen(r) {
		}
	}()
}

// Listen to a SRCDS log output stream
func (o *Observer) Listen(r io.Reader) <-chan LogEntry {
	br := bufio.NewReader(r)

	firstLine, err := br.ReadString(0x0A)
	if err == io.EOF {
		return nil
	}

	// Determine EOL delimiter as it may not match operating system's EOL
	if strings.HasSuffix(firstLine, eolWindows) {
		log.Debug().Msg("Windows EOL delimiter detected (0x0D0A).")
		o.EndOfLine = eolWindows
	} else {
		if strings.HasSuffix(firstLine, eolUnix) {
			log.Debug().Msg("Unix EOL delimiter detected (0x0A).")
		} else {
			log.Warn().Msg("Couldn't detect EOL delimiter; defaulting to Unix EOL (0x0A).")
		}

		o.EndOfLine = eolUnix
	}

	o.wg.Add(1)
	logStream := make(chan LogEntry, 6)
	runtime.GC()

	go func(firstLine string, br io.Reader, c chan<- LogEntry) {
		log.Info().Msg("Now observing the SRCDS log stream")

		defer close(c)
		defer o.wg.Done()

		o.processMessage(firstLine, c)

		scanner := bufio.NewScanner(br)
		for scanner.Scan() {
			o.processMessage(scanner.Text(), c)
		}
	}(firstLine, br, logStream)

	return logStream
}

// Wait for the SRCDS observer to exit naturally.
func (o *Observer) Wait() {
	o.wg.Wait()
}

// TryCvarAsInt attempts to return a cvar as an integer, returning a bool indicating if the provided fallback value was returned
func (o *Observer) TryCvarAsInt(name string, fallback int) (value int, nonFallback bool) {
	return o.cvars.tryInt(name, fallback)
}

// TryCvarAsString attempts to return a cvar as an integer, returning a bool indicating if the provided fallback value was returned
func (o *Observer) TryCvarAsString(name, fallback string) (value string, nonFallback bool) {
	return o.cvars.tryString(name, fallback)
}

type Observer struct {
	cvars      Cvars
	EndOfLine  string
	started    time.Time
	statistics observerStatistics
	wg         sync.WaitGroup
}

type observerStatistics struct {
	totalLines uint32
	blankLines uint32
	logLines   uint32
}

func (o *Observer) processMessage(line string, outEntries chan<- LogEntry) {
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

		if outEntries != nil {
			outEntries <- le
		}

		return
	}

	if cvarSet, ok := parseCvarResponse(line); ok {
		o.cvars.setIfWatched(cvarSet.Name, cvarSet.Value, time.Now())
		return
	}
}
