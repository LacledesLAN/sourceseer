# SourceSeer Developer Notes for Counter-Strike: Global Offensive

## Useful CSGO Commands

| Command           | What               |
| ----------------- | ------------------ |
| `mp_warmup_end`   | Ends the warmup.   |
| `mp_warmup_start` | Starts the warmup. |

## Game Loop

[]			- Loading map "de_dust2"
[server cvars]		- server cvars start
			- server cvars end
[match-start]		- World triggered "Match_Start" on "de_train"
[]			- Team playing "CT": Retirement
[]			- Team playing "TERRORIST": Not_fast
[round-start]		- World triggered "Round_Start"
[]			- Team "CT" scored "11" with "2" players
[]			- Team "TERRORIST" scored "10" with "1" players
[round-end]		- World triggered "Round_End
[match-end]		- Game Over: competitive de_train score 11:16 after 26 min

## Config File Load Order

1. `autoexec.cfg` - executed before the first map starts.
2. `server.cfg` - executed every map change
3. `gamemode_competitive_server.cfg` - executed at every map change

## Code Scraps

// Max amount of time for players to readyup for the very first map (seconds)
maxReadyupConnectTime = time.Minute * 7
// Max amount of seconds for players to readyup for knife rounds 2+
maxReadyUpKnifeTime = time.Minute * 3
// If all players not ready by this time period auto start (seconds)
maxReadyUpPlayTime = time.Minute * 3


func Wrapper(osArgs []string) Reactor {
	if len() < 2 {
		return ReadWriter{}, errors.New("Not enough arguments")
	}

	// Prepare process
	process := exec.Command(osArgs[0], osArgs[1:len(osArgs)]...)

	// Link Standard Error
	if stdErr, err := process.StderrPipe(); err != nil {
		return ReadWriter{}, errors.New("Unable to link standard error")
	} else {
		defer stdErr.Close()

		go func(r *bufio.Reader) {
			for {
				errLine, _ := r.ReadString('\n')
				errLine = strings.TrimSpace(strings.TrimSuffix(errLine, "\n"), "")

				Log(SRCDSError, errLine)
			}
		}(bufio.NewReader(stdErr))
	}

	out := make(chan LogEntry, 12)
	//Link Standard Out
	if stdOut, err := process.StdoutPipe(); err!= nil {
		return ReadWriter{}, errors.New("Unable to link standard out")
	} else {

	}

	cmd := make(chan string, 6)
	//Link Standard In

	s.started = time.Now()
	if err := process.Start(); err != nil {
		return err
	}
	defer s.process.Process.Kill()

	if err := process.Wait(); err != nil {
		return err
	}

	return nil
}
