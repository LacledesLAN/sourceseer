package srcds

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	// MaxHostnameLength for all SRCDS servers
	MaxHostnameLength int    = 28
	srcdsTimeLayout   string = "1/2/2006 - 15:04:05"
)

// SRCDS represents a source dedicated server
type SRCDS struct {
	cmdIn      chan string
	cvars      map[string]string
	game       GameServer
	launchArgs []string
	started    time.Time
	finished   time.Time
}

type GameServer interface {
	ClientConnected(Client)
	ClientDisconnected(ClientDisconnected)
	CmdSender() chan string
	LogReceiver(LogEntry)
}

// AddCvarWatch to watch for and update from the log stream.
func (s *SRCDS) AddCvarWatch(names ...string) {
	for _, name := range names {
		name = strings.Trim(name, "")

		if len(name) > 0 {
			if _, found := s.cvars[name]; !found {
				s.cvars[name] = ""
			}
		}
	}
}

// AddLaunchArg to be used when initializing the SRCDS instance.
func (s *SRCDS) AddLaunchArg(args ...string) {
	for _, arg := range args {
		arg = strings.Trim(arg, "")
		if len(arg) > 0 {
			s.launchArgs = append(s.launchArgs, arg)
		}
	}
}

// GetCvar value and a boolean as to if the value was found or not.
func (s *SRCDS) GetCvar(name string) (value string, found bool) {
	value, found = s.cvars[name]
	return
}

// GetCvarAsInt attempts to return a cvar as an integer
func (s *SRCDS) GetCvarAsInt(name string) (value int, err error) {
	v, found := s.cvars[name]

	if !found {
		return 0, errors.New("cvar '" + name + "' was not found.")
	}

	return strconv.Atoi(v)
}

// New creates and wraps around srcds instance.
func New(gameRunner GameServer, osArgs []string) (SRCDS, error) {
	s := SRCDS{
		cmdIn:      make(chan string, 12),
		cvars:      make(map[string]string),
		launchArgs: osArgs,
	}

	return s, nil
}

// RefreshCvars triggers SRCDS to echo all watched cvars to the log stream.
func (s *SRCDS) RefreshCvars() {
	go func() {
		for name := range s.cvars {
			s.cmdIn <- name
		}
	}()
}

// Start the instance of the SRCDS; connecting a channel to its standard input stream.
func (s *SRCDS) Start() error {
	srcdsProcess := exec.Command(s.launchArgs[0], s.launchArgs[1:len(s.launchArgs)-1]...)

	// link standard error
	stdErr, err := srcdsProcess.StderrPipe()
	if err != nil {
		return errors.New("Unable to link standard error")
	}
	defer stdErr.Close()
	go func(reader *bufio.Reader) {
		for {
			errLine, _ := reader.ReadString('\n')
			errLine = strings.Trim(strings.TrimSuffix(errLine, "\n"), "")

			if len(errLine) > 0 {
				log.Println("Standard Error:>", errLine)
			}
		}
	}(bufio.NewReader(stdErr))

	// link standard out
	stdOut, err := srcdsProcess.StdoutPipe()
	if err != nil {
		return errors.New("Unable to link standard out")
	}
	defer stdOut.Close()
	go func(reader *bufio.Reader) {
		for {
			outLine, _ := reader.ReadString('\n')
			outLine = strings.Trim(strings.TrimSuffix(outLine, "\n"), "")

			if len(outLine) > 0 {
				le := parseLogEntry(outLine)

				if len(le.Message) > 0 {
					s.processLogEntry(le)
				}
			}
		}
	}(bufio.NewReader(stdOut))

	// go routine to grab stdin, combine with timer, send to process
	stdIn, err := srcdsProcess.StdinPipe()
	if err != nil {
		return errors.New("Unable to link standard in")
	}
	defer stdIn.Close()

	go func(downStream io.WriteCloser, upStream chan string) {
		timer := time.NewTimer(time.Millisecond * 500).C
		var lastSent time.Time

		for {
			select {
			case s := <-upStream:
				downStream.Write([]byte(s))
				downStream.Write([]byte("\n"))
				lastSent = time.Now()
			case <-timer:
				// ensure srcds buffer get flushed at regular interval
				if time.Now().Sub(lastSent) > (time.Millisecond * 750) {
					downStream.Write([]byte("\n"))
					lastSent = time.Now()
				}

				timer = time.NewTimer(time.Millisecond * 500).C
			}
		}
	}(stdIn, s.cmdIn)

	// Start SRCDS
	fmt.Println("Starting srcds using", s.launchArgs)

	s.started = time.Now()
	err = srcdsProcess.Start()
	defer srcdsProcess.Process.Kill()

	if err != nil {
		panic(err)
	}

	err = srcdsProcess.Wait()
	if err != nil {
		panic(err)
	}

	s.finished = time.Now()

	return nil
}

func (s *SRCDS) processLogEntry(le LogEntry) (keepProcessing bool) {

	cvarSet, err := parseCvarValueSet(le)
	if err != nil {
		if _, found := s.cvars[cvarSet.name]; found {
			s.cvars[cvarSet.name] = cvarSet.value
		}

		return false
	}

	if strings.HasPrefix(le.Message, `"`) {
		client, err := parseClientConnected(le)
		if err != nil {
			s.game.ClientConnected(client)
			return false
		}

		clientDisconnected, err := parseClientDisconnected(le)
		if err != nil {
			s.game.ClientDisconnected(clientDisconnected)
			return false
		}
	}

	s.game.LogReceiver(le)

	return true
}
