package csgo

import (
	"fmt"
	"sync"

	"github.com/lacledeslan/sourceseer/pkg/srcds"
)

type Server struct {
	srcds *srcds.Server
	Observer
	wg sync.WaitGroup
}

func NewServer() *Server {
	s := &Server{
		srcds: srcds.NewServer(),
	}

	s.Observer.srcdsObserver = s.srcds.Observer

	s.srcds.AddCvarWatcher("mp_halftime", "mp_maxrounds", "mp_overtime_maxrounds")

	return s
}

func (s *Server) SetExec(arg string, args ...string) error {
	err := s.srcds.SetExec(arg, args...)
	if err != nil {
		return fmt.Errorf("Unable to SetExec for CSGO Server: %w", err)
	}

	return nil
}

// Read starts the CSGO server and processes its output
func (s *Server) Read() error {
	c, err := s.Listen()
	if err != nil {
		return err
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for range c {
		}
	}()

	return nil
}

// Listen starts the CSGO server, processes its output, and returns its log stream
func (s *Server) Listen() (<-chan srcds.LogEntry, error) {
	c, err := s.srcds.Listen()
	if err != nil {
		return nil, fmt.Errorf("Couldn't listen to SRCDS server: %w", err)
	}

	s.wg.Add(1)
	logStream := make(chan srcds.LogEntry, 6)
	go func(in <-chan srcds.LogEntry, out chan<- srcds.LogEntry) {
		defer s.wg.Done()
		defer close(out)
		for le := range in {
			s.processLogEntry(le)
			out <- le
		}
	}(c, logStream)

	return logStream, nil
}

// Wait blocks until the CSGO server stops executing
func (s *Server) Wait() {
	s.wg.Wait()
}

func (s *Server) serverProcessLogEntry(le srcds.LogEntry) {

}
