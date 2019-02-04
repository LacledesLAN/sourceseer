package srcds

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	// MaxHostnameLength is the maximum length allowed for srcds's hostname
	MaxHostnameLength int    = 28
	srcdsTimeLayout   string = "1/2/2006 - 15:04:05"
)

type SrcdsWrapper struct {
	cvars   map[string]string
	started time.Time
}

func newWrapper() SrcdsWrapper {
	return SrcdsWrapper{}
}

func (m *SrcdsWrapper) GetCvar(name string) (value string, found bool) {
	value, found = m.cvars[name]
	return
}

func (m *SrcdsWrapper) UpdatedCvar(name, value string) {
	if _, found := m.cvars[name]; found {
		m.cvars[name] = value
	}
}

func (m *SrcdsWrapper) WatchCvar(names ...string) {
	for _, name := range names {
		name = strings.Trim(name, "")

		if len(name) == 0 {
			continue
		}

		if _, found := m.cvars[name]; !found {
			m.cvars[name] = ""
		}
	}
}

func (m *SrcdsWrapper) TryUpdateFromStdOut(logEntry LogEntry) (consumed bool) {
	if strings.HasPrefix(logEntry.Message, `server_cvar: "`) {
		result := serverCvarSetRegex.FindStringSubmatch(logEntry.Message)

		if result != nil && len(result) >= 3 {
			m.UpdatedCvar(result[1], result[2])
			return true
		}
	} else if result := serverCvarEchoRegex.FindStringSubmatch(logEntry.Message); len(result) > 1 {
		m.UpdatedCvar(result[1], result[2])
		return true
	}

	return false
}

// WrapProc starts the SRCDS instance and connects the needed pipes
func WrapProc(srcdsArgs []string, stdIn <-chan string, stdOut chan<- LogEntry) *SrcdsWrapper {
	var args []string
	var srcdsProcess *exec.Cmd

	m := newWrapper()

	if _, err := os.Stat("/app/srcds_run"); err == nil {
		// we're inside docker
	} else {
		switch os := runtime.GOOS; os {
		case "windows":
			args = []string{"powershell.exe", "-NonInteractive", "-Command", "docker", "run", "-i", "--rm", "-p 27015:27015", "-p 27015:27015/udp", "lltest/gamesvr-csgo-tourney", "./srcds_run"}
		default:
			args = []string{"docker", "run", "-i", "--rm", "-p 27015:27015", "-p 27015:27015/udp", "lltest/gamesvr-csgo-tourney", "./srcds_run"}
		}
	}

	args = append(args, srcdsArgs...)

	srcdsProcess = exec.Command(args[0], args[1:len(args)-1]...)

	// Link standard error
	cmdStderr, err := srcdsProcess.StderrPipe()
	defer cmdStderr.Close()
	go func(reader *bufio.Reader) {
		for {
			errLine, _ := reader.ReadString('\n')
			errLine = strings.Trim(strings.TrimSuffix(errLine, "\n"), "")

			if len(errLine) > 0 {
				log.Println("Received output from standard error:", errLine)
			}
		}
	}(bufio.NewReader(cmdStderr))

	// Link standard out
	if stdOut != nil {
		cmdStdout, err := srcdsProcess.StdoutPipe()
		if err != nil {
			panic("unable to link standard out")
		}
		defer cmdStdout.Close()

		go func(reader *bufio.Reader) {
			for {
				outLine, _ := reader.ReadString('\n')
				outLine = strings.Trim(strings.TrimSuffix(outLine, "\n"), "")

				if len(outLine) > 0 {
					logEntry := ExtractLogEntry(outLine)

					if !m.TryUpdateFromStdOut(logEntry) {
						stdOut <- logEntry
					}
				}
			}
		}(bufio.NewReader(cmdStdout))
	}

	// Link standard in
	if stdIn != nil {
		cmdStdin, err := srcdsProcess.StdinPipe()
		if err != nil {
			panic("unable to link standard in")
		}
		defer cmdStdin.Close()

		go func() {
			timer := time.NewTimer(time.Millisecond * 500).C
			var lastSent time.Time

			for {
				select {
				case s := <-stdIn:
					cmdStdin.Write([]byte(s + "\n"))
					lastSent = time.Now()
				case <-timer:
					// ensure srcds's buffer is flushed on a regular basis
					if time.Now().Sub(lastSent) > (time.Millisecond * 750) {
						cmdStdin.Write([]byte("\n"))
						lastSent = time.Now()
					}

					timer = time.NewTimer(time.Millisecond * 500).C
				}
			}
		}()
	}

	// Start process
	fmt.Println("Starting server using", args)

	err = srcdsProcess.Start()
	defer srcdsProcess.Process.Kill()

	if err != nil {
		panic(err)
	}

	m.started = time.Now()

	srcdsProcess.Wait()

	return &m
}
