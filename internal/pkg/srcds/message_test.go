package srcds

import (
	"testing"
	"time"
)

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

func Test_ParseCvarValueSet(t *testing.T) {
	goodDatum := []struct {
		srcdsMessage     string
		expectedVarName  string
		expectedVarValue string
	}{
		{`server_cvar: "mp_whatever2" "-300"`, "mp_whatever2", "-300"},
		{`server_cvar: "cash_player_interact_with_hostage" "150"`, "cash_player_interact_with_hostage", "150"},
		{`server_cvar: "cash_team_rescued_hostage" "0"`, "cash_team_rescued_hostage", "0"},
		{`server_cvar: "mp_roundtime_hostage" "1.92"`, "mp_roundtime_hostage", "1.92"},
		{`server_cvar: "mp_whatever" ""`, "mp_whatever", ""},
		{`server_cvar: "cash_team_elimination_hostage_map_ct" "3000"`, "cash_team_elimination_hostage_map_ct", "3000"},
		{`server_cvar: "sm_nextmap" "de_dust2"`, "sm_nextmap", "de_dust2"},
		{`"cash_player_killed_teammate" = "-300"`, "cash_player_killed_teammate", "-300"},
		{`"cash_player_respawn_amount" = "0"`, "cash_player_respawn_amount", "0"},
		{`"sv_maxspeed" = "320"`, "sv_maxspeed", "320"},
		{`"mp_teamlist" = "hgrunt;scientist"`, "mp_teamlist", "hgrunt;scientist"},
		{`"sourcemod_version" = "1.9.0.6148"`, "sourcemod_version", "1.9.0.6148"},
		{`"sv_tags" = ""`, "sv_tags", ""},
		{`"mp_respawnwavetime" = "10.0"`, "mp_respawnwavetime", "10.0"},
		{`"metamod_version" = "1.11.0-dev+1097V"`, "metamod_version", "1.11.0-dev+1097V"},
		{`"mp_do_warmup_period" = "1" min. 0.000000 max. 1.000000 game replicated          - Whether or not to do a warmup period at the start of a match.`, "mp_do_warmup_period", "1"},
		{`"mp_maxrounds" = "7" ( def. "0" ) min. 0.000000 game notify replicated           - max number of rounds to play before server changes maps`, "mp_maxrounds", "7"},
		{`"mp_maxrounds" = "30" ( def. "0" ) min. 0.000000 game notify replicated          - max number of rounds to play before server changes maps`, "mp_maxrounds", "30"},
		{`"mp_overtime_maxrounds" = "7" ( def. "6" ) client replicated                     - When overtime is enabled play additional rounds to determine winner`, "mp_overtime_maxrounds", "7"},
	}

	for _, testData := range goodDatum {
		t.Run(testData.srcdsMessage, func(t *testing.T) {
			result, err := ParseCvarValueSet(testData.srcdsMessage)

			if err != nil {
				t.Errorf("Error result should be nil but got %q", err)
			}

			if result.Name != testData.expectedVarName {
				t.Errorf("Expected var name %q but got %q", testData.expectedVarName, result.Name)
			}

			if result.Value != testData.expectedVarValue {
				t.Errorf("Expected var name %q but got %q", testData.expectedVarValue, result.Value)
			}
		})
	}

	badDatum := []string{
		`My first clue came at 4:15, when the clock stopped.`,
		`Console<0><Console><Console>" say ""Running gamemode_competitive_server.cfg"`,
	}

	for _, testData := range badDatum {
		t.Run(testData, func(t *testing.T) {
			result, err := ParseCvarValueSet(testData)

			if err == nil {
				t.Errorf("Error result should not be nill; but go a name of '%q' with a value of '%q'.", result.Name, result.Value)
			}
		})
	}
}

func Test_parseLogEntry(t *testing.T) {
	goodData := []struct {
		actualRaw    string
		expectedMsg  string
		expectedTime time.Time
	}{
		{"L 1/2/2000 - 03:04:00: Sweet llamas of the Bahamas!", "Sweet llamas of the Bahamas!", time.Unix(946803840, 0)},
		{"L 01/2/2000 - 03:04:00: Excuse my language but I have had it with you ruffling my petticoats!", "Excuse my language but I have had it with you ruffling my petticoats!", time.Unix(946803840, 0)},
		{"L 1/02/2000 - 03:04:00: Your music is bad & you should feel bad!", "Your music is bad & you should feel bad!", time.Unix(946803840, 0)},
		{"L 01/02/2000 - 03:04:00: Did everything just taste purple for a second?", "Did everything just taste purple for a second?", time.Unix(946803840, 0)},
		{"L 01/02/2000 - 3:04:00: When you look this good, you don’t have to know anything!", "When you look this good, you don’t have to know anything!", time.Unix(946803840, 0)},
	}

	for _, testData := range goodData {
		result, err := parseLogEntry(testData.actualRaw)

		if err != nil {
			t.Errorf("Return err should be nil not '%q'", err)
		}

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

	//{"Server will auto-restart if there is a crash.", "", time.Time{}},
}
