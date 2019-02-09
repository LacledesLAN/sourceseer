package srcds

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

const (
	MaxHostnameLength int    = 28
	srcdsTimeLayout   string = "1/2/2006 - 15:04:05"
)

type LogEntryProcessor func(LogEntry) (keepProcessing bool)

// SRCDS represents a source dedicated server
type SRCDS struct {
	cmdIn             chan string
	cvars             map[string]string
	launchArgs        []string
	logProcessorStack LogEntryProcessor
	started           time.Time
	finished          time.Time
}

func New(osArgs []string) (SRCDS, error) {
	s := SRCDS{
		cvars:      make(map[string]string),
		launchArgs: osArgs,
	}
	s.logProcessorStack = s.processLogEntry

	return s, nil
}

func (s *SRCDS) AddLogProcessor(p LogEntryProcessor) {
	if p != nil {
		prev := s.logProcessorStack
		s.logProcessorStack = func(le LogEntry) (keepProcessing bool) {
			if !prev(le) {
				return false
			}

			return p(le)
		}
	}
}

func (s *SRCDS) Start(cmdInd chan string) {
	srcdsProcess := exec.Command(s.launchArgs[0], s.launchArgs[1:len(s.launchArgs)-1]...)

	stdErr, err := srcdsProcess.StderrPipe()
	if err != nil {
		defer stdErr.Close()
		s.linkStdErr(stdErr)
	}

	stdOut, err := srcdsProcess.StdoutPipe()
	if err != nil {
		fmt.Println("Unable to link std out!")
		panic(err)
	}
	defer stdOut.Close()
	s.linkStdOut(stdOut)

	// go routine to grab stdin, combine with timer, send to process
	stdIn, err := srcdsProcess.StdinPipe()
	if err != nil {
		fmt.Println("Unable to link standard in!")
		panic(err)
	}
	defer stdIn.Close()
	s.linkStdIn(stdIn)

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
}

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

func (s *SRCDS) AddLaunchArg(args ...string) {
	for _, arg := range args {
		arg = strings.Trim(arg, "")
		if len(arg) > 0 {
			s.launchArgs = append(s.launchArgs, arg)
		}
	}
}

func (s *SRCDS) GetCvar(name string) (value string, found bool) {
	value, found = s.cvars[name]
	return
}

func (s *SRCDS) processLogEntry(le LogEntry) (keepProcessing bool) {
	fmt.Println("RRRRR ----- ", le.Message)

	if strings.HasPrefix(le.Message, `server_cvar: "`) {
		result := serverCvarSetRegex.FindStringSubmatch(le.Message)

		if result != nil && len(result) >= 3 {
			s.updatedCvar(result[1], result[2])
			return false
		}
	} else if result := serverCvarEchoRegex.FindStringSubmatch(le.Message); len(result) > 1 {
		s.updatedCvar(result[1], result[2])
		return false
	}

	return true
}

func (s *SRCDS) ReconcileCvars() {
	go func() {
		for name := range s.cvars {
			s.cmdIn <- name
		}
	}()
}

func (s *SRCDS) updatedCvar(name, value string) {
	if _, found := s.cvars[name]; found {
		s.cvars[name] = value
	}
}

func (s *SRCDS) linkStdErr(e io.ReadCloser) {
	go func(reader *bufio.Reader) {
		for {
			errLine, _ := reader.ReadString('\n')
			errLine = strings.Trim(strings.TrimSuffix(errLine, "\n"), "")

			if len(errLine) > 0 {
				log.Println("Standard Error:>", errLine)
			}
		}
	}(bufio.NewReader(e))
}

func (s *SRCDS) linkStdOut(i io.ReadCloser) {
	go func(reader *bufio.Reader) {
		for {
			outLine, _ := reader.ReadString('\n')
			outLine = strings.Trim(strings.TrimSuffix(outLine, "\n"), "")

			if len(outLine) > 0 {
				le := ExtractLogEntry(outLine)

				if len(le.Message) > 0 {
					s.logProcessorStack(le)
				} else {
					fmt.Println(outLine)
				}
			}
		}
	}(bufio.NewReader(i))
}

func (s *SRCDS) linkStdIn(stdIn io.WriteCloser) {
	timer := time.NewTimer(time.Millisecond * 500).C
	var lastSent time.Time

	if s.cmdIn == nil {
		s.cmdIn = make(chan string, 12)
	}

	go func() {
		for {
			select {
			case s := <-s.cmdIn:
				stdIn.Write([]byte(s))
				stdIn.Write([]byte("\n"))
				lastSent = time.Now()
			case <-timer:
				// ensure srcds's buffer is flushed on a regular basis
				if time.Now().Sub(lastSent) > (time.Millisecond * 750) {
					stdIn.Write([]byte("\n"))
					lastSent = time.Now()
				}

				timer = time.NewTimer(time.Millisecond * 500).C
			}
		}
	}()
}
