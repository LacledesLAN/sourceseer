package srcds

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

const (
	// MaxHostnameLength for all SRCDS servers
	MaxHostnameLength int = 28

	optionEchoSRCDSCvar  bool   = false
	optionEchoSRCDSError bool   = true
	optionEchoSRCDSLog   bool   = true
	optionEchoSRCDSOther bool   = true
	srcdsTimeLayout      string = "1/2/2006 - 15:04:05"
)

// SRCDS represents a source dedicated server
type SRCDS struct {
	CmdIn    chan string
	game     Game
	started  time.Time
	finished time.Time
}

// Game represents a SRCDS game
type Game interface {
	ClientConnected(Client)
	ClientDisconnected(ClientDisconnected)
	ClientMessage(ClientMessage)
	CmdSender() chan string
	CvarSet(name, value string)
	LaunchArgs() []string
	LogReceiver(LogEntry)
}

// New creates and wraps around srcds instance.
func New(srcdsGame Game) (SRCDS, error) {
	s := SRCDS{
		CmdIn: make(chan string, 12),
		game:  srcdsGame,
	}

	return s, nil
}

// Start the instance of the SRCDS; connecting a channel to its standard input stream.
func (s *SRCDS) Start(osArgs []string) error {

	srcdsProcess := exec.Command(osArgs[0], append(osArgs[1:len(osArgs)], s.game.LaunchArgs()...)...)

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

			if len(errLine) > 0 && optionEchoSRCDSError {
				fmt.Println("[SRCDS ERR ]", errLine)
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
				le, err := parseLogEntry(outLine)

				if err == nil {
					s.processLogEntry(le)
					if optionEchoSRCDSLog {
						fmt.Println("[SRCDS LOG ]", le.Message)
					}
				} else if cvarSet, err := ParseCvarValueSet(outLine); err == nil {
					s.game.CvarSet(cvarSet.Name, cvarSet.Value)

					if optionEchoSRCDSCvar {
						fmt.Println("[SRCDS CVAR]", outLine)
					}
				} else {
					if optionEchoSRCDSOther {
						fmt.Println("[SRCDS     ]", outLine)
					}
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
	}(stdIn, s.game.CmdSender())

	// Start SRCDS
	fmt.Print("\n\n/======================================================================================\\\n")
	fmt.Print("[SOURCESEER] Starting using", srcdsProcess.Args, "\n\n")

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

func (s *SRCDS) processLogEntry(le LogEntry) {
	if strings.HasPrefix(le.Message, `"`) {
		if clientMsg, err := parseClientMessage(le); err == nil {
			if client, err := parseClientConnected(clientMsg); err == nil {
				s.game.ClientConnected(client)
				return
			}

			if clientDisconnected, err := parseClientDisconnected(clientMsg); err == nil {
				s.game.ClientDisconnected(clientDisconnected)
				return
			}

			s.game.ClientMessage(clientMsg)
			return
		}
	}

	s.game.LogReceiver(le)
}
