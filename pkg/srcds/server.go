package srcds

import (
	"bufio"
	"context"
	"fmt"
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

// Server represents an interactive SRCDS instance
type Server struct {
	*Observer
	process *exec.Cmd
	cmdIn   chan string
	wg      sync.WaitGroup
}

// NewServer for interacting with a SRCDS instance
func NewServer() *Server {
	s := &Server{
		Observer: NewObserver(),
		cmdIn:    make(chan string, 4),
	}

	return s
}

// SetExec prepares the SRCDS instance for execution using the given arguments
func (s *Server) SetExec(path string, args ...string) error {
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		/// TODO: add back in safety requiring signal to be sent twice in x seconds?
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, syscall.SIGINT)
		defer cancel()
		defer close(sig)

		<-sig
	}(cancel)

	return s.SetExecContext(ctx, path, args...)
}

// SetExecContext is like SetExec but includes a context
func (s *Server) SetExecContext(ctx context.Context, arg string, args ...string) error {
	var osArgs []string

	if runtime.GOOS == "windows" {
		osArgs = append(osArgs, "powershell.exe", "-NonInteractive", "-Command")
	}

	osArgs = append(osArgs, arg)
	if len(args) > 0 {
		osArgs = append(osArgs, args...)
	}

	s.process = exec.Command(osArgs[0], osArgs[1:len(osArgs)]...)

	s.linkStdIn(ctx)

	return nil
}

func (s *Server) linkStdIn(ctx context.Context) error {
	cmdStdIn, err := s.process.StdinPipe()
	if err != nil {
		return fmt.Errorf("Couldn't connect to process's standard in pipe: %w", err)
	}

	// connect to the process's standard in
	go func(wc io.WriteCloser, cmdIn <-chan string) {
		defer wc.Close()
		s.wg.Add(1)
		defer s.wg.Done()
		prev := time.Time{}
		ticker := time.NewTicker(175 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				// Politely ask SRCDS to shutdown
				log.Info().Msg("Attempting to gracefully shut down the SRCDS server")
				io.WriteString(wc, "say server shutting down"+s.EndOfLine)
				time.Sleep(250 * time.Millisecond)
				io.WriteString(wc, "quit"+s.EndOfLine)
				time.Sleep(750 * time.Millisecond)
				return
			case cmd := <-cmdIn:
				// Send the command to process's standard in
				prev = time.Now()
				log.Info().Msgf("Sending command %q to SRCDS", cmd)
				io.WriteString(wc, cmd+s.EndOfLine)
			case <-ticker.C:
				// Send EOL to flush the process's standard out buffer
				if time.Since(prev) >= 100*time.Millisecond {
					prev = time.Now()
					io.WriteString(wc, s.EndOfLine)
				}
			}
		}
	}(cmdStdIn, s.cmdIn)

	// grab the terminal's standard in
	go func(r io.ReadCloser) {
		defer r.Close()
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			s.SendCommand(scanner.Text())
		}
	}(os.Stdin)

	return nil
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
func (s *Server) SendCommand(l string) {
	l = strings.TrimSpace(l)
	if len(l) > 0 {
		s.cmdIn <- l
	}
}

// Listen starts the SRCDS server, processes its output, and returns its log stream
func (s *Server) Listen() (<-chan LogEntry, error) {
	if s.process == nil {
		return nil, errors.New("Exec was never set")
	}

	cmdStdOut, err := s.process.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect to process's standard out pipe: %w", err)
	}

	if err := s.process.Start(); err != nil {
		return nil, fmt.Errorf("Problem executing process: %w", err)
	}

	log.Debug().Msg("Server execution started")

	return s.Observer.Listen(cmdStdOut), nil
}

// Read starts the SRCDS server and processes its output
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

// Wait blocks until the SRCDS server stops executing
func (s *Server) Wait() {
	s.wg.Wait()
}
