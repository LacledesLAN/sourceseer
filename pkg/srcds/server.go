package srcds

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

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

// Command returns the Server struct for launching and wrapping a SRCDS instance.
func Command(name string, arg ...string) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		s := make(chan os.Signal, 2)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
		lastSigInt := time.Time{}

		for sig := range s {
			if sig == syscall.SIGINT && time.Since(lastSigInt) > time.Second*8 {
				fmt.Fprintf(os.Stderr, "Press ctrl^c again in the next 8 seconds to terminate.")
				lastSigInt = time.Now()
				continue
			}

			fmt.Fprintf(os.Stderr, "Terminating process.")
			cancel()
			time.Sleep(2 * time.Second)
			os.Exit(0)
		}
	}(cancel)

	return CommandContext(ctx, name, arg...)
}

// CommandContext is like command but includes a context.
func CommandContext(ctx context.Context, name string, arg ...string) (*Server, error) {
	cmd := exec.CommandContext(ctx, name, arg...)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmdStdIn, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	server := &Server{
		cmdIn:    make(chan string, 6),
		observer: *newReader(stdOut),
	}

	go func(wc io.WriteCloser, s *Server) {
		prev := time.Time{}
		ticker := time.NewTicker(90 * time.Millisecond)
		defer wc.Close()
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Send EOL to flush the process's standard out buffer
				if time.Since(prev) >= 100*time.Millisecond {
					prev = time.Now()
					wc.Write([]byte(s.endOfLine))
				}
			case cmd := <-server.cmdIn:
				// Send the command to process's standard in
				prev = time.Now()
				wc.Write([]byte(cmd + s.endOfLine))
			case <-ctx.Done():
				// Politely ask SRCDS to shutdown
				wc.Write([]byte("say server shutting down" + s.endOfLine))
				wc.Write([]byte("quit" + s.endOfLine))
				time.Sleep(500 * time.Millisecond)
				return
			}
		}
	}(cmdStdIn, server)

	go func(r io.ReadCloser) {
		defer r.Close()
		s := bufio.NewScanner(r)

		for s.Scan() {
			server.cmdIn <- s.Text()
		}
	}(os.Stdin)

	return server, nil
}

// Server represents an interactive SRCDS instance
type Server struct {
	cmdIn chan string
	observer
}
