package csgo

import (
	"testing"
	"time"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
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

func Test_parseLoadingMap(t *testing.T) {
	datum := []struct {
		srcdsMessage string
		expectedMap  string
	}{
		{`Loading map "de_nuke"`, "de_nuke"},
		{`Loading map "de_dust2"`, "de_dust2"},
	}

	for _, testData := range datum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result, err := parseLoadingMap(srcds.LogEntry{Message: testData.srcdsMessage})

			if err != nil {
				t.Errorf("Parsing `%q` should not have resulted in an error of `%q`.", testData.srcdsMessage, err)
			}

			if result != testData.expectedMap {
				t.Errorf("Expected '%q but got '%q'.", testData.expectedMap, result)
			}
		})
	}
}

func Test_parsePlayerSay(t *testing.T) {
	datum := []struct {
		srcdsMessage    string
		expectedChannel PlayerSaidChannel
		expectedMessage string
	}{
		{`"Malfunctioning Eddie<2><STEAM_1:1:86753090><CT>" say_team "I'm malfunctioning so badly, I'm practically giving these cars away!"`, ChannelAffiliation, `I'm malfunctioning so badly, I'm practically giving these cars away!`},
		//{``, "", ``},
		//{``, "", ``},
	}
	for _, testData := range datum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result, err := parsePlayerSay(srcds.LogEntry{Message: testData.srcdsMessage})

			if err != nil {
				t.Errorf("Unable to parse player say; received error '%q'.", err)
			}

			//
			// check player info
			//

			if result.channel != testData.expectedChannel {
				t.Errorf("Expected channel %q but got %q", testData.expectedChannel, result.channel)
			}

			if result.msg != testData.expectedMessage {
				t.Errorf("Expected msg `%q` but got `%q`.", testData.expectedMessage, result.msg)
			}
		})
	}
}

func Test_parseTeamScored(t *testing.T) {
	datum := []struct {
		srcdsMessage        string
		expectedAffiliation string
		expectedScore       int
		expectedPlayerCount int
	}{
		{`Team "CT" scored "0" with "72" players`, "CT", 0, 72},
		{`Team "TERRORIST" scored "196" with "512" players`, "TERRORIST", 196, 512},
		{`Team "CT" scored "86" with "0" players`, "CT", 86, 0},
	}

	for _, testData := range datum {
		result, err := parseTeamScored(srcds.LogEntry{Message: testData.srcdsMessage})

		if err != nil {
			t.Errorf("Parsing `%q` should not have resulted in an error of `%q`.", testData.srcdsMessage, err)
		}

		if result.teamAffiliation != testData.expectedAffiliation {
			t.Errorf("Expected affiliation %q but got %q", testData.expectedAffiliation, result.teamAffiliation)
		}

		if result.teamScore != testData.expectedScore {
			t.Errorf("Expected team score of %q but got %q", testData.expectedScore, result.teamScore)
		}

		if result.teamPlayerCount != testData.expectedPlayerCount {
			t.Errorf("Expected team player count of %q but got %q", testData.expectedPlayerCount, result.teamPlayerCount)
		}
	}
}

func Test_parseTeamSetSide(t *testing.T) {
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
			result, err := parseTeamSetSide(srcds.LogEntry{Message: testData.srcdsMessage})

			if err != nil {
				t.Errorf("Parsing `%q` should not have resulted in an error of `%q`.", testData.srcdsMessage, err)
			}

			if result.teamAffiliation != testData.expectedSide {
				t.Errorf("Expected side %q but got %q", testData.expectedSide, result.teamAffiliation)
			}

			if result.teamName != testData.expectedTeam {
				t.Errorf("Expected team name of %q but got %q", testData.expectedTeam, result.teamName)
			}
		})
	}
}

func Test_parseTeamTriggered(t *testing.T) {
	datum := []struct {
		srcdsMessage        string
		expectedAffiliation string
		expectedTrigger     string
		expectedCTScore     int
		expectedTScore      int
	}{
		{`Team "CT" triggered "SFUI_Notice_Bomb_Defused" (CT "21") (T "7")`, "CT", "SFUI_Notice_Bomb_Defused", 21, 7},
		{`Team "CT" triggered "SFUI_Notice_CTs_Win" (CT "124") (T "0")`, "CT", "SFUI_Notice_CTs_Win", 124, 0},
		{`Team "CT" triggered "SFUI_Notice_Target_Saved" (CT "12") (T "3")`, "CT", "SFUI_Notice_Target_Saved", 12, 3},
		{`Team "TERRORIST" triggered "SFUI_Notice_Target_Bombed" (CT "0") (T "5")`, "TERRORIST", "SFUI_Notice_Target_Bombed", 0, 5},
		{`Team "TERRORIST" triggered "SFUI_Notice_Terrorists_Win" (CT "6") (T "23")`, "TERRORIST", "SFUI_Notice_Terrorists_Win", 6, 23},
	}

	for _, testData := range datum {
		result, err := parseTeamTriggered(srcds.LogEntry{Message: testData.srcdsMessage})

		if err != nil {
			t.Errorf("Parsing `%q` should not have resulted in an error of `%q`.", testData.srcdsMessage, err)
		}

		if result.teamAffiliation != testData.expectedAffiliation {
			t.Errorf("Parsing `%q` should resulted in an affiliation of `%q` not `%q`.", testData.srcdsMessage, testData.expectedAffiliation, result.teamAffiliation)
		}

		if result.triggered != testData.expectedTrigger {
			t.Errorf("Parsing `%q` should resulted in an triggered of `%q` not `%q`.", testData.srcdsMessage, testData.expectedTrigger, result.triggered)
		}

		if result.ctScore != testData.expectedCTScore {
			t.Errorf("Parsing `%q` should resulted in a ct score of `%d` not `%d`.", testData.srcdsMessage, testData.expectedCTScore, result.ctScore)
		}

		if result.terroristScore != testData.expectedTScore {
			t.Errorf("Parsing `%q` should resulted in a terrorist score of `%d` not `%d`.", testData.srcdsMessage, testData.expectedTScore, result.terroristScore)
		}
	}
}

func Test_parseWorldTriggered(t *testing.T) {
	inTimeSpan := func(start, end, check time.Time) bool {
		return check.After(start) && check.Before(end)
	}

	// Round Start
	r, _ := parseWorldTriggered(srcds.LogEntry{Message: `World triggered "Round_Start"`})
	if r.trigger != RoundStart {
		t.Error("Should have RoundStart trigger")
	}

	// Round End
	r, _ = parseWorldTriggered(srcds.LogEntry{Message: `World triggered "Round_End"`})
	if r.trigger != RoundEnd {
		t.Error("Should have RoundEnd trigger")
	}

	// Restart Round
	r, _ = parseWorldTriggered(srcds.LogEntry{Message: `World triggered "Restart_Round_(25_seconds)"`})
	if r.trigger != RoundRestarting {
		t.Error("Should have RoundRestarting trigger")
	}
	if !inTimeSpan(time.Now().Add(23*time.Second), time.Now().Add(26*time.Second), r.eta) {
		t.Error("ETA was out of range")
	}

	// Match Start
	r, _ = parseWorldTriggered(srcds.LogEntry{Message: `World triggered "Match_Start" on "de_nuke"`})
	if r.trigger != MatchStart {
		t.Error("Should have MatchStart trigger")
	}

	// Game Commencing
	r, _ = parseWorldTriggered(srcds.LogEntry{Message: `World triggered "Game_Commencing"`})
	if r.trigger != GameCommencing {
		t.Error("Should have GameCommencing trigger")
	}
}
