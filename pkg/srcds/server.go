package srcds

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// NewServer returns the Server struct for launching and wrapping a SRCDS instance.
func NewServer(arg string, args ...string) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		/// TODO: add back in safety requiring signal to be sent twice in x seconds?
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, syscall.SIGINT)
		defer cancel()
		defer close(sig)

		<-sig
	}(cancel)

	return NewServerContext(ctx, arg, args...)
}

// NewServerContext is like command but includes a context.
func NewServerContext(ctx context.Context, arg string, args ...string) (*Server, error) {
	var osArgs []string
	switch os := runtime.GOOS; os {
	case "windows":
		osArgs = append(osArgs, "powershell.exe", "-NonInteractive", "-Command")
	}

	osArgs = append(osArgs, arg)
	if len(args) > 0 {
		osArgs = append(osArgs, args...)
	}

	cmd := exec.Command(osArgs[0], osArgs[1:len(osArgs)]...)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Errorf("Couldn't connect to process's standard out pipe: %w", err)
	}

	cmdStdIn, err := cmd.StdinPipe()
	if err != nil {
		return nil, errors.Errorf("Couldn't connect to process's standard in pipe: %w", err)
	}

	server := wrapProcess(ctx, stdOut, cmdStdIn)

	// grab the terminal's standard in; send it to srcds
	go func(r io.ReadCloser) {
		defer r.Close()
		s := bufio.NewScanner(r)

		for s.Scan() {
			server.SendCommand(s.Text())
		}
	}(os.Stdin)

	server.start = func() {
		err = cmd.Start()
		if err != nil {
			panic("Couldn't start process")
		}
	}

	return server, nil
}

// Server represents an interactive SRCDS instance
type Server struct {
	cmdIn chan string
	*observer
	waitGroup sync.WaitGroup
	start     func()
}

// RefreshWatchedCvars triggers SRCDS into echoing all watched cvars to the log stream.
func (s *Server) RefreshWatchedCvars() {
	go func(s *Server) {
		for _, name := range s.cvars.getNames() {
			s.SendCommand(name)
			time.Sleep(10 * time.Millisecond)
		}
	}(s)
}

// SendCommand to the interactive SRCDS instance.
func (s *Server) SendCommand(cmd string) {
	c := strings.TrimSpace(cmd)
	if len(c) > 0 {
		s.cmdIn <- c
	}
}

// Wait blocks until the server shuts down
func (s *Server) Wait() {
	s.waitGroup.Wait()
	s.observer.Wait()
}

// Start the SRCDS Server
func (s *Server) Start() <-chan LogEntry {
	if s.start == nil {
		panic("srcds > server > start function was not instantiated.")
	}

	c := s.observer.start()
	s.start()
	return c
}

func newServer(stdOut io.Reader) *Server {
	return &Server{
		cmdIn:    make(chan string, 6),
		observer: newReader(stdOut),
	}
}

// wrap a csgo srcds process
func wrapProcess(ctx context.Context, processStdOut io.Reader, processStdIn io.WriteCloser) *Server {
	server := newServer(processStdOut)

	// connect to the process's standard in
	go func(wc io.WriteCloser, s *Server) {
		s.waitGroup.Add(1)
		prev := time.Time{}
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		defer wc.Close()
		defer s.waitGroup.Done()

		for {
			select {
			case <-ctx.Done():
				// Politely ask SRCDS to shutdown
				log.Info().Msg("Attempting to gracefully shut down the SRCDS server")
				io.WriteString(wc, "say server shutting down"+s.endOfLine)
				io.WriteString(wc, "quit"+s.endOfLine)
				time.Sleep(750 * time.Millisecond)
				return
			case cmd := <-server.cmdIn:
				// Send the command to process's standard in
				prev = time.Now()
				log.Info().Msgf("Sending command %q to SRCDS", cmd)
				io.WriteString(wc, cmd+s.endOfLine)
			case <-ticker.C:
				// Send EOL to flush the process's standard out buffer
				if time.Since(prev) >= 100*time.Millisecond {
					prev = time.Now()
					io.WriteString(wc, s.endOfLine)
				}
			}
		}
	}(processStdIn, server)

	return server
}
