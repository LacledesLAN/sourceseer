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

// LogEntryProcessor represents a function that can parse log entires; returning false when the log entry has been consumed or its effects undone.
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

// AddLogProcessor to top of the log processor stack.
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

// GetCvar value and a boolean as to if the value was found or not.
func (s *SRCDS) GetCvar(name string) (value string, found bool) {
	value, found = s.cvars[name]
	return
}

func (s *SRCDS) GetCvarAsInt(name string) (value int, err error) {
	v, found := s.cvars[name]

	if !found {
		return 0, errors.New("cvar '" + name + "' was not found.")
	}

	return strconv.Atoi(v)
}

// New creates a new CSGO server instance.
func New(osArgs []string) (SRCDS, error) {
	s := SRCDS{
		cvars:      make(map[string]string),
		launchArgs: osArgs,
	}
	s.logProcessorStack = s.processLogEntry

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
				// ensure srcds buffer get flushed at regular interval
				if time.Now().Sub(lastSent) > (time.Millisecond * 750) {
					stdIn.Write([]byte("\n"))
					lastSent = time.Now()
				}

				timer = time.NewTimer(time.Millisecond * 500).C
			}
		}
	}()
}

func (s *SRCDS) linkStdOut(i io.ReadCloser) {
	go func(reader *bufio.Reader) {
		for {
			outLine, _ := reader.ReadString('\n')
			outLine = strings.Trim(strings.TrimSuffix(outLine, "\n"), "")

			if len(outLine) > 0 {
				le := parseLogEntry(outLine)

				if len(le.Message) > 0 {
					s.logProcessorStack(le)
				} else {
					//fmt.Println(outLine)
				}
			}
		}
	}(bufio.NewReader(i))
}

func (s *SRCDS) processLogEntry(le LogEntry) (keepProcessing bool) {

	cvarSet, err := parseCvarValueSet(le)
	if err != nil {
		s.updatedCvar(cvarSet.name, cvarSet.value)
		return false
	}

	return true
}

func (s *SRCDS) updatedCvar(name, value string) {
	if _, found := s.cvars[name]; found {
		s.cvars[name] = value
	}
}
