package csgo

import (
	"testing"
)

func Test_gameOverRegex(t *testing.T) {
	datum := []struct {
		srcdsMessage   string
		expectedMode   string
		expectedMap    string
		expectedScore1 string
		expectedScore2 string
		expectedTime   string
	}{
		{"Game Over: competitive  de_nuke score 18:19 after 34 min", "competitive", "de_nuke", "18", "19", "34"},
		{"Game Over: competitive  de_inferno score 3:16 after 38 min", "competitive", "de_inferno", "3", "16", "38"},
	}

	for _, testData := range datum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result := gameOverRegex.FindStringSubmatch(testData.srcdsMessage)

			resultMode := result[1]
			if resultMode != testData.expectedMode {
				t.Errorf("Expected mode %q but got %q", testData.expectedMode, resultMode)
			}

			resultMap := result[2]
			if resultMap != testData.expectedMap {
				t.Errorf("Expected map %q but got %q", testData.expectedMap, resultMap)
			}

			resultScore1 := result[3]
			if resultScore1 != testData.expectedScore1 {
				t.Errorf("Expected score 1 of %q but got %q", testData.expectedScore1, resultScore1)
			}

			resultScore2 := result[4]
			if resultScore2 != testData.expectedScore2 {
				t.Errorf("Expected score 2 of %q but got %q", testData.expectedScore2, resultScore2)
			}

			resultTime := result[5]
			if resultTime != testData.expectedTime {
				t.Errorf("Expected team %q but got %q", testData.expectedTime, resultTime)
			}
		})
	}
}

func Test_loadingMapRegex(t *testing.T) {
	datum := []struct {
		srcdsMessage string
		expectedMap  string
	}{
		{`Loading map "de_nuke"`, "de_nuke"},
		{`Loading map "de_dust2"`, "de_dust2"},
	}

	for _, testData := range datum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result := loadingMapRegex.FindStringSubmatch(testData.srcdsMessage)
			resultMap := result[1]

			if resultMap != testData.expectedMap {
				t.Errorf("Expected map %q but got %q", testData.expectedMap, resultMap)
			}
		})
	}
}

func Test_matchStartPattern(t *testing.T) {
	datum := []struct {
		srcdsMessage string
		expectedMap  string
	}{
		{`World triggered "Match_Start" on "de_nuke"`, "de_nuke"},
	}

	for _, testData := range datum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result := matchStartRegex.FindStringSubmatch(testData.srcdsMessage)
			resultMap := result[1]

			if resultMap != testData.expectedMap {
				t.Errorf("Expected map %q but got %q", testData.expectedMap, resultMap)
			}
		})
	}
}

func Test_teamScoredPattern(t *testing.T) {
	datum := []struct {
		srcdsMessage        string
		expectedSide        string
		expectedScore       string
		expectedPlayerCount string
	}{
		{`Team "CT" scored "0" with "7" players`, "CT", "0", "7"},
		{`Team "TERRORIST" scored "16" with "5" players`, "TERRORIST", "16", "5"},
	}

	for _, testData := range datum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result := teamScoredRegex.FindStringSubmatch(testData.srcdsMessage)
			resultSide := result[1]
			resultScore := result[2]
			resultPlayerCount := result[3]

			if resultSide != testData.expectedSide {
				t.Errorf("Expected side %q but got %q", testData.expectedSide, resultSide)
			}

			if resultScore != testData.expectedScore {
				t.Errorf("Expected team score of %q but got %q", testData.expectedScore, resultScore)
			}

			if resultPlayerCount != testData.expectedPlayerCount {
				t.Errorf("Expected team player count of %q but got %q", testData.expectedPlayerCount, resultPlayerCount)
			}
		})
	}
}

func Test_teamSetSidePattern(t *testing.T) {
	datum := []struct {
		srcdsMessage string
		expectedSide string
		expectedTeam string
	}{
		{`Team playing "CT": a`, "CT", "a"},
		{`Team playing "TERRORIST": them`, "TERRORIST", "them"},
	}

	for _, testData := range datum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result := teamSetSideRegex.FindStringSubmatch(testData.srcdsMessage)
			resultAffiliation := result[1]
			resultTeamName := result[2]

			if resultAffiliation != testData.expectedSide {
				t.Errorf("Expected side %q but got %q", testData.expectedSide, resultAffiliation)
			}

			if resultTeamName != testData.expectedTeam {
				t.Errorf("Expected team %q but got %q", testData.expectedTeam, resultTeamName)
			}
		})
	}
}

func Test_worldTriggeredMatchStartRegex(t *testing.T) {
	datum := []struct {
		srcdsMessage string
		expectedMap  string
	}{
		{`World triggered "Match_Start" on "de_nuke"`, "de_nuke"},
	}

	for _, testData := range datum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result := worldTriggeredMatchStartRegex.FindStringSubmatch(testData.srcdsMessage)
			resultMap := result[1]

			if resultMap != testData.expectedMap {
				t.Errorf("Expected map %q but got %q", testData.expectedMap, resultMap)
			}
		})
	}
}

func Test_worldTriggeredRoundRestartRegex(t *testing.T) {
	datum := []struct {
		srcdsMessage    string
		expectedSeconds string
	}{
		{`World triggered "Restart_Round_(1_second)"`, "1"},
		{`World triggered "Restart_Round_(34_second)"`, "34"},
		{`World triggered "Restart_Round_(191_second)"`, "191"},
	}

	for _, testData := range datum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result := worldTriggeredRoundRestartRegex.FindStringSubmatch(testData.srcdsMessage)
			resultTime := result[1]

			if resultTime != testData.expectedSeconds {
				t.Errorf("Expected %q seconds but got %q", testData.expectedSeconds, resultTime)
			}
		})
	}
}
