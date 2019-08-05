package srcds

import "time"

// Server represents an interactive SRCDS instance
type Server interface {
	Observer
	RefreshCvars()
	SendCmd(cmd string)
}

// RefreshCvars attempts to trigger SRCDS into echoing all watched cvars to the log stream.
func (s *server) RefreshCvars() {
	go func(s *server) {
		for name := range s.cvars {
			s.SendCommand(name)
			time.Sleep(10 * time.Millisecond)
		}
	}(s)
}

// SendCommand to the interactive SRCDS instance
func (s *server) SendCommand(cmd string) {
	return
}

// Wrapper for observing and interacting with a SRCDS instance
func Wrapper(osArgs ...string) Server {
	return nil
}

type server struct {
	bufCmdIn chan string
	observer
}

func newServer() *server {
	return &server{
		bufCmdIn: make(chan string, 9),
		observer: *newObserver(),
	}
}
