package csgo

import "testing"

func Test_calculateLastRoundWinThreshold(t *testing.T) {
	testCases := map[string]struct {
		mpMaxRounds         int
		mpOvertimeMaxRounds int
		scenarios           map[lastInt][]lastInt // map[winTreshold]lastCompletedRounds
	}{
		"Normal Rounds Clinchable": {
			mpMaxRounds:         1,
			mpOvertimeMaxRounds: 2,
			scenarios: map[lastInt][]lastInt{
				1: []lastInt{-9999, -1, 0, 1, 2, 9999},
			},
		},
		"Normal Rounds Not Clinchable - Overtime Clinchable": {
			mpMaxRounds:         2,
			mpOvertimeMaxRounds: 3,
			scenarios: map[lastInt][]lastInt{
				2: []lastInt{-9999, -1, 0, 1, 2},
				3: []lastInt{3, 4, 9999},
			},
		},
		"Normal Rounds Not Clinchable - Overtime Not Clinchable": {
			mpMaxRounds:         2,
			mpOvertimeMaxRounds: 2,
			scenarios: map[lastInt][]lastInt{
				2:  []lastInt{-9999, -1, 0, 1, 2},
				3:  []lastInt{3, 4},
				4:  []lastInt{5, 6},
				5:  []lastInt{7, 8},
				59: []lastInt{115, 116},
			},
		},
		"Hasty Settings": {
			mpMaxRounds:         4,
			mpOvertimeMaxRounds: 3,
			scenarios: map[lastInt][]lastInt{
				3: []lastInt{-9999, 0, 1, 2, 3, 4},
				4: []lastInt{5, 6, 7, 9999},
			},
		},
		"Default Settings": {
			mpMaxRounds:         30,
			mpOvertimeMaxRounds: 6,
			scenarios: map[lastInt][]lastInt{
				16: []lastInt{-9999, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
				19: []lastInt{31, 32, 33, 34, 35, 36},
				22: []lastInt{37, 38, 39, 40, 41, 42},
				28: []lastInt{49},
				31: []lastInt{55},
				34: []lastInt{61},
				37: []lastInt{67},
				40: []lastInt{73},
				43: []lastInt{79},
				49: []lastInt{91, 92, 93, 94, 95, 96},
			},
		},
		"Prevailing Tourney Settings": {
			mpMaxRounds:         30,
			mpOvertimeMaxRounds: 7,
			scenarios: map[lastInt][]lastInt{
				16: []lastInt{-9999, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
				19: []lastInt{31, 32, 33, 34, 35, 36, 37, 9999},
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			for expected, lastCompletedRounds := range test.scenarios {
				for _, lastCompletedRound := range lastCompletedRounds {
					actual := calculateLastRoundWinThreshold(test.mpMaxRounds, test.mpOvertimeMaxRounds, lastCompletedRound)

					if actual != expected {
						t.Errorf("With `mp_maxrounds` = `%d` and `mp_overtime_maxrounds` = `%d` and the last completed round being `%d` a calculated win threshold should be %d not %d!",
							test.mpMaxRounds, test.mpOvertimeMaxRounds, lastCompletedRound, expected, actual)
					}
				}
			}
		})
	}
}

func Test_calcOvertimePeriodNumber(t *testing.T) {
	testCases := map[string]struct {
		mpMaxRounds          int
		mpOvertimeMaxRounds  int
		overtimePeriodRounds map[int][]lastInt // map[overtimeNumber]lastCompletedRounds
	}{
		"Hasty Settings": {
			mpMaxRounds:         4,
			mpOvertimeMaxRounds: 3,
			overtimePeriodRounds: map[int][]lastInt{
				0: {0, 1, 2, 3},
				1: {4, 5, 6, 7, 8, 9, 10, 9999}, // OT rounds 7+ will never happen; should still report OT period 1
			},
		},
		"Default Settings": {
			mpMaxRounds:         30,
			mpOvertimeMaxRounds: 6,
			overtimePeriodRounds: map[int][]lastInt{
				0:  {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
				1:  {30, 31, 32, 33, 34, 35},
				2:  {36, 37, 38, 39, 40, 41},
				3:  {42, 43, 44, 45, 46, 47},
				9:  {78, 79, 80, 81, 82, 83},
				12: {96, 97, 98, 99, 100, 101},
			},
		},
		"Prevailing Tourney Settings": {
			mpMaxRounds:         30,
			mpOvertimeMaxRounds: 7,
			overtimePeriodRounds: map[int][]lastInt{
				0: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
				1: {30, 31, 32, 33, 34, 35, 36, 37, 39, 39, 40, 9999}, // OT rounds 37+ will never happen; should still report OT period 1
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			for expectedOTNumber, lastCompletedRounds := range test.overtimePeriodRounds {
				for _, lastCompletedRound := range lastCompletedRounds {
					actual := calcOvertimePeriodNumber(test.mpMaxRounds, test.mpOvertimeMaxRounds, lastCompletedRound)
					if actual != expectedOTNumber {
						t.Errorf("With `mp_maxrounds` = `%d`, `mp_overtime_maxrounds` = `%d`, and `last completed round = %d` the overtime period number should be %v not %v",
							test.mpMaxRounds, test.mpOvertimeMaxRounds, lastCompletedRound, expectedOTNumber, actual)
					}
				}
			}
		})
	}
}

func Test_calculateSidesAreCurrentlySwitched(t *testing.T) {
	type scenario struct {
		expectTeamSwitch    bool
		lastCompletedRounds []lastInt
	}

	testCases := map[string]struct {
		mpHalftime          int
		mpMaxRounds         int
		mpOvertimeMaxRounds int
		scenarios           map[string]scenario
	}{
		"Hasty Settings": {
			mpHalftime:          1,
			mpMaxRounds:         4,
			mpOvertimeMaxRounds: 3,
			scenarios: map[string]scenario{
				"First-Half":            {false, []lastInt{0, 1}},
				"Second-Half":           {true, []lastInt{2, 3, 4}},
				"Overtime - First-Half": {false, []lastInt{5, 6}},
			},
		},
		"Default Settings": {
			mpHalftime:          1,
			mpMaxRounds:         30,
			mpOvertimeMaxRounds: 6,
			scenarios: map[string]scenario{
				"First-Half":       {false, []lastInt{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
				"Second-Half":      {true, []lastInt{15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}},
				"OT 1 First-Half":  {true, []lastInt{30, 31, 32}},
				"OT 1 Second-Half": {false, []lastInt{33, 34, 35}},
				"OT 2 First-Half":  {false, []lastInt{36, 37, 38}},
				"OT 2 Second-Half": {true, []lastInt{39, 40, 41}},
				"OT 3 First-Half":  {true, []lastInt{42, 43, 44}},
				"OT 3 Second-Half": {false, []lastInt{45, 46, 47}},
			},
		},
		"Prevailing Tourney Settings": {
			mpHalftime:          1,
			mpMaxRounds:         30,
			mpOvertimeMaxRounds: 7,
			scenarios: map[string]scenario{
				"First-Half":             {false, []lastInt{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
				"Second-Half":            {true, []lastInt{15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}},
				"Overtime - First-Half":  {true, []lastInt{30, 31, 32}},
				"Overtime - Second-Half": {false, []lastInt{33, 34, 35, 36}},
			},
		},
	}

	for testName, test := range testCases {
		for scenarioName, scenario := range test.scenarios {
			t.Run(testName+"-"+scenarioName, func(t *testing.T) {
				for _, lastCompletedRound := range scenario.lastCompletedRounds {
					actualTeamSwitche := calculateSidesAreCurrentlySwitched(test.mpHalftime, test.mpMaxRounds, test.mpOvertimeMaxRounds, lastCompletedRound)

					if actualTeamSwitche != scenario.expectTeamSwitch {
						t.Errorf("With `mp_halftime = %d`, `mp_maxrounds` = `%d`, `mp_overtime_maxrounds` = `%d`, and `last completed round = %d` the team sides swapped should be: %v",
							test.mpHalftime, test.mpMaxRounds, test.mpOvertimeMaxRounds, lastCompletedRound, scenario.expectTeamSwitch)
					}
				}
			})
		}
	}
}
