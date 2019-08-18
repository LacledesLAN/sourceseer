package srcds

import (
	"bufio"
	"context"
	"os/exec"
	"strings"
	"time"
)

// RefreshCvars attempts to trigger SRCDS into echoing all watched cvars to the log stream.
func (s *Server) RefreshCvars() {
	go func(s *Server) {
		for _, name := range s.cvars.getNames() {
			s.SendCommand(name)
			time.Sleep(10 * time.Millisecond)
		}
	}(s)
}

// SendCommand to the interactive SRCDS instance
func (s *Server) SendCommand(cmd string) {

}

// Command for observing and interacting with a SRCDS instance
func Command(name string, osArgs ...string) (*Server, error) {
	ctx, _ := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, name, osArgs...)

	return wrapProcess(cmd)
}

// Server represents an interactive SRCDS instance
type Server struct {
	cmdIn chan string
	observer
}

func wrapProcess(cmd *exec.Cmd) (*Server, error) {
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	//if stdErr, err := cmd.StderrPipe(); err == nil {
	//// TODO read and process standard error
	//}

	return &Server{
		cmdIn:    make(chan string, 9),
		observer: *newReader(stdOut),
	}, nil
}

func (s *Server) linkStdOut(r *bufio.Reader) {
	go func(r *bufio.Reader) {
		for {
			outLine, _ := r.ReadString('\n')
			outLine = strings.Trim(strings.TrimSuffix(outLine, "\n"), "")

			if len(outLine) > 0 {

				Log(SRCDSOther, outLine)
			}
		}
	}(r)
}
