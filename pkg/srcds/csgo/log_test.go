package csgo

import (
	"bufio"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func Test_ParseLogs(t *testing.T) {
	logDir := "." + string(os.PathSeparator) + "testdata" + string(os.PathSeparator)

	tests := map[string]struct {
		mpHalftime          int
		mpMaxRounds         int
		mpMaxOvertimeRounds int
		expected            observerStatistics
	}{
		//all_team_triggers.log
		"mpteam1_domination.log": {mpHalftime: 1, mpMaxRounds: 30, mpMaxOvertimeRounds: 7,
			expected: observerStatistics{roundsCompleted: 16}},

		"mpteam1_overtimewin.log": {mpHalftime: 1, mpMaxRounds: 30, mpMaxOvertimeRounds: 7,
			expected: observerStatistics{roundsCompleted: 35}},

		"mpteam2_clinch.log": {mpHalftime: 1, mpMaxRounds: 30, mpMaxOvertimeRounds: 7,
			expected: observerStatistics{roundsCompleted: 26}},

		"tourney_1x_match.log": {mpHalftime: 1, mpMaxRounds: 30, mpMaxOvertimeRounds: 7,
			expected: observerStatistics{roundsCompleted: 19}},

		"tourney_3map_clinch.log": {mpHalftime: 1, mpMaxRounds: 30, mpMaxOvertimeRounds: 7,
			expected: observerStatistics{roundsCompleted: 51}},

		"tourney_5x_match.log": {mpHalftime: 1, mpMaxRounds: 30, mpMaxOvertimeRounds: 7,
			expected: observerStatistics{roundsCompleted: 88}},

		"warmup_120seconds.log": {mpHalftime: 1, mpMaxRounds: 4, mpMaxOvertimeRounds: 3,
			expected: observerStatistics{roundsCompleted: 4}},

		"warmup_disabled.log": {mpHalftime: 1, mpMaxRounds: 4, mpMaxOvertimeRounds: 3,
			expected: observerStatistics{roundsCompleted: 7}},

		"warmup_manual.log": {mpHalftime: 1, mpMaxRounds: 4, mpMaxOvertimeRounds: 3,
			expected: observerStatistics{roundsCompleted: 3}},
	}

	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	for logFile, test := range tests {
		t.Run(logFile, func(t *testing.T) {
			logFile = logDir + logFile

			file, err := os.Open(logFile)
			if err != nil {
				t.Fatalf("Could not open log %q file for parsing: %e", logFile, err)
			}
			defer file.Close()

			reader := bufio.NewReader(file)

			csgo := NewReader(reader, test.mpHalftime, test.mpMaxRounds, test.mpMaxOvertimeRounds)
			csgo.Start()

			//if csgo.statistics.roundsStarted != test.expected.roundsStarted {
			//	t.Errorf("Expected %02d rounds to have been started but saw %02d", test.expected.roundsStarted, csgo.statistics.roundsStarted)
			//}

			if csgo.statistics.roundsCompleted != test.expected.roundsCompleted {
				t.Errorf("Expected %02d rounds to have been completed, not %02d.", test.expected.roundsCompleted, csgo.statistics.roundsCompleted)
			}

			//if csgo.statistics.matchesStarted != test.expected.matchesStarted {
			//	t.Errorf("Expected %02d matches to have been started but saw %02d", test.expected.matchesStarted, csgo.statistics.matchesStarted)
			//}
		})
	}
}
