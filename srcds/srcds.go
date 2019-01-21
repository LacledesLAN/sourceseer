package srcds

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	// MaxHostnameLength is the maximum length allowed for srcds's hostname
	MaxHostnameLength int = 28
)

var (
	extractPlayerInfo = regexp.MustCompile(`"(.{1,32})<(\d{0,1})><([a-zA-Z0-9:_]*)><{0,1}([a-zA-Z0-9]*?)>{0,1}" ([^[\-?\d+ -?\d+ -?\d+\]]?)`) //group 1: name; group 2: slot; group 3: uid; group 4: team (if any)
	mapChange         = regexp.MustCompile(`^(?:[Ll]oading|[Ss]tarted)? map "([A-Za-z0-9_]+)".*`)
	srcdsTimeLayout   = "1/2/2006 - 15:04:05"
)

// LogEntry represents a log message from srcds
type LogEntry struct {
	Message   string
	Raw       string
	Timestamp time.Time
}

// ExtractLogEntry extracts the log entry from a srcds log message
func ExtractLogEntry(rawLogEntry string) LogEntry {
	rawLogEntry = strings.TrimSpace(rawLogEntry)

	if len(rawLogEntry) > 0 && strings.HasPrefix(rawLogEntry, "L ") && strings.Contains(rawLogEntry, ": ") {
		i := strings.Index(rawLogEntry, ": ")

		msgTime, err := time.ParseInLocation(srcdsTimeLayout, rawLogEntry[2:i], time.Local)

		if err == nil {
			return LogEntry{
				Raw:       rawLogEntry,
				Timestamp: msgTime,
				Message:   rawLogEntry[i+2:],
			}
		}
	}

	return LogEntry{Raw: rawLogEntry}
}

// ExtractClients extracts the players from a srcds log message
func ExtractClients(logEntry LogEntry) (originator, target *Client) {
	originator = nil
	target = nil

	players := extractPlayerInfo.FindAllStringSubmatch(logEntry.Message, -1)

	if len(players) >= 1 {
		originatorRaw := players[0]
		originator = &Client{Username: originatorRaw[1], ServerSlot: originatorRaw[2], ServerTeam: originatorRaw[4], SteamID: originatorRaw[3]}
	}

	if len(players) >= 2 {
		targetRaw := players[1]
		target = &Client{Username: targetRaw[1], ServerSlot: targetRaw[2], ServerTeam: targetRaw[4], SteamID: targetRaw[3]}
	}

	return
}

// WrapProc starts the SRCDS instance and connects the needed pipes
func WrapProc(srcdsArgs []string, stdIn <-chan string, stdOut, stdError chan<- string) {
	var args []string
	var srcdsProcess *exec.Cmd

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
	if stdError != nil {
		cmdStderr, err := srcdsProcess.StderrPipe()
		if err != nil {
			panic("unable to link standard error")
		}
		defer cmdStderr.Close()

		go func(reader *bufio.Reader) {
			for {
				errLine, _ := reader.ReadString('\n')
				errLine = strings.Trim(strings.TrimSuffix(errLine, "\n"), "")

				if len(errLine) > 0 {
					stdError <- errLine
				}
			}
		}(bufio.NewReader(cmdStderr))
	}

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
					stdOut <- outLine
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
	fmt.Println("calling start.....")
	fmt.Println(args)

	err := srcdsProcess.Start()
	defer srcdsProcess.Process.Kill()

	if err != nil {
		panic(err)
	}

	srcdsProcess.Wait()
}
