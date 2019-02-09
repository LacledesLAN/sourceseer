package srcds

import (
	"testing"
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
