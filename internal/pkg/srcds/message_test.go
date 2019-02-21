package srcds

import (
	"testing"
	"time"
)

func Test_serverCvarEchoRegex(t *testing.T) {
	//TODO - what happens when variable value has either whitespace or a double quote in it?

	datum := []struct {
		srcdsMessage     string
		expectedVarName  string
		expectedVarValue string
	}{
		{`"cash_player_killed_teammate" = "-300"`, "cash_player_killed_teammate", "-300"},
		{`"cash_player_respawn_amount" = "0"`, "cash_player_respawn_amount", "0"},
		{`"sv_maxspeed" = "320"`, "sv_maxspeed", "320"},
		{`"mp_teamlist" = "hgrunt;scientist"`, "mp_teamlist", "hgrunt;scientist"},
		{`"sourcemod_version" = "1.9.0.6148"`, "sourcemod_version", "1.9.0.6148"},
		{`"sv_tags" = ""`, "sv_tags", ""},
		{`"mp_respawnwavetime" = "10.0"`, "mp_respawnwavetime", "10.0"},
		{`"metamod_version" = "1.11.0-dev+1097V"`, "metamod_version", "1.11.0-dev+1097V"},
		{`"mp_do_warmup_period" = "1" min. 0.000000 max. 1.000000 game replicated          - Whether or not to do a warmup period at the start of a match.`, "mp_do_warmup_period", "1"},
		{`"mp_maxrounds" = "30" ( def. "0" ) min. 0.000000 game notify replicated          - max number of rounds to play before server changes maps`, "mp_maxrounds", "30"},
	}

	for _, testData := range datum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result := serverCvarEchoRegex.FindStringSubmatch(testData.srcdsMessage)

			if len(result) == 0 {
				t.Errorf("Was unable to parse %q", testData.srcdsMessage)
			} else {
				if result[1] != testData.expectedVarName {
					t.Errorf("Expected var name %q but got %q", testData.expectedVarName, result[1])
				}

				if result[2] != testData.expectedVarValue {
					t.Errorf("Expected var name %q but got %q", testData.expectedVarValue, result[2])
				}
			}
		})
	}
}

func Test_parseClientConnected(t *testing.T) {
	goodDatum := []string{
		`"Beelzebot | theinfosphere.org<2><STEAM_1:0:8675309><>" connected, address ""`,
		`"[FUTURAMA] Hedonismbot<3><STEAM_1:0:8675309><>" entered the game`,
	}

	for _, testString := range goodDatum {
		t.Run(testString, func(t *testing.T) {
			_, err := parseClientConnected(LogEntry{Message: testString})

			if err != nil {
				t.Errorf("Should have parse successfully but got error: %q", err)
			}
		})
	}

}

func Test_parseClientDisconnected(t *testing.T) {
	goodDatum := []string{
		`"Don<17><BOT><CT>" disconnected (reason "Kicked by Console")`,
		`"Don<17><BOT><CT>" disconnected`,
	}

	for _, testString := range goodDatum {
		t.Run(testString, func(t *testing.T) {
			_, err := parseClientDisconnected(LogEntry{Message: testString})

			if err != nil {
				t.Errorf("Should have parse successfully but got error: %q", err)
			}
		})
	}
}

func Test_serverCvarSetPattern(t *testing.T) {
	datum := []struct {
		srcdsMessage     string
		expectedVarName  string
		expectedVarValue string
	}{
		{`server_cvar: "mp_whatever2" "-300"`, "mp_whatever2", "-300"},
		{`server_cvar: "cash_player_interact_with_hostage" "150"`, "cash_player_interact_with_hostage", "150"},
		{`server_cvar: "cash_team_rescued_hostage" "0"`, "cash_team_rescued_hostage", "0"},
		{`server_cvar: "mp_roundtime_hostage" "1.92"`, "mp_roundtime_hostage", "1.92"},
		{`server_cvar: "mp_whatever" ""`, "mp_whatever", ""},
	}

	for _, testData := range datum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result := serverCvarSetRegex.FindStringSubmatch(testData.srcdsMessage)

			if result[1] != testData.expectedVarName {
				t.Errorf("Expected var name %q but got %q", testData.expectedVarName, result[1])
			}

			if result[2] != testData.expectedVarValue {
				t.Errorf("Expected var name %q but got %q", testData.expectedVarValue, result[2])
			}
		})
	}
}

func Test_parseLogEntry(t *testing.T) {
	testDatum := []struct {
		actualRaw    string
		expectedMsg  string
		expectedTime time.Time
	}{
		{"Server will auto-restart if there is a crash.", "", time.Time{}},
		{"L 1/2/2000 - 03:04:00: Sweet llamas of the Bahamas!", "Sweet llamas of the Bahamas!", time.Unix(946803840, 0)},
		{"L 01/2/2000 - 03:04:00: Excuse my language but I have had it with you ruffling my petticoats!", "Excuse my language but I have had it with you ruffling my petticoats!", time.Unix(946803840, 0)},
		{"L 1/02/2000 - 03:04:00: Your music is bad & you should feel bad!", "Your music is bad & you should feel bad!", time.Unix(946803840, 0)},
		{"L 01/02/2000 - 03:04:00: Did everything just taste purple for a second?", "Did everything just taste purple for a second?", time.Unix(946803840, 0)},
		{"L 01/02/2000 - 3:04:00: When you look this good, you don’t have to know anything!", "When you look this good, you don’t have to know anything!", time.Unix(946803840, 0)},
	}

	for _, testData := range testDatum {
		result := parseLogEntry(testData.actualRaw)

		if result.Message != testData.expectedMsg {
			t.Errorf("Expected message of %q but got %q", testData.expectedMsg, result.Message)
		}

		if result.Raw != testData.actualRaw {
			t.Errorf("Expected raw message of %q but got %q", testData.actualRaw, result.Raw)
		}

		if result.Timestamp != testData.expectedTime {
			t.Errorf("Expected timestamp of %q but got %q", testData.expectedTime, result.Timestamp)
		}
	}
}
