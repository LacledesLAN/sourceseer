package csgo

import (
	"github.com/lacledeslan/sourceseer/pkg/srcds"
)

type Server struct {
	srcds *srcds.Server
	Observer
}

// Start the CSGO server
func (s *Server) Start() {
	for le := range s.srcds.Start() {
		s.processLogEntry(le)
	}
}

func (s *Server) Wait() {
	s.srcds.Wait()
}

// NewServer for CSGO
func NewServer(arg string, args ...string) (*Server, error) {
	srcds, err := srcds.NewServer(arg, args...)
	if err != nil {
		return nil, err
	}

	srcds.AddCvarWatcher("mp_halftime", "mp_maxrounds", "mp_overtime_maxrounds")

	s := &Server{
		srcds: srcds,
	}

	return s, nil
}
