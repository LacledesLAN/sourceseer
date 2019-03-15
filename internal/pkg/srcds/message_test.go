package srcds

import (
	"testing"
	"time"
)

func Test_parseClientConnected(t *testing.T) {
	client := Client{Username: "Nannybot 1.0", ServerSlot: "1", SteamID: "[U:1:28500804", Affiliation: ""}
	goodMessages := []string{
		`connected, address ""`,
		`connected, address "192.168.1.52"`,
		`connected, address "2001:0DB8:AC10:FE01:0000:EE91:0000:0000"`,
		`entered the game`,
	}

	for _, testMsg := range goodMessages {
		t.Run(testMsg, func(t *testing.T) {
			_, err := parseClientConnected(ClientMessage{Client: client, Message: testMsg})

			if err != nil {
				t.Errorf("Should have parse successfully but got error: %q", err)
			}
		})
	}

}

func Test_parseClientDisconnected(t *testing.T) {
	client := Client{Username: "Sinclair 2K", ServerSlot: "1", SteamID: "[U:1:28500804", Affiliation: ""}
	goodDatum := []struct {
		msg            string
		expectedReason string
	}{
		{`disconnected (reason "Kicked by Console")`, "Kicked by Console"},
		{`disconnected`, ""},
	}

	for _, testData := range goodDatum {
		t.Run(testData.msg, func(t *testing.T) {
			cltDisconnect, err := parseClientDisconnected(ClientMessage{Client: client, Message: testData.msg})

			if err != nil {
				t.Errorf("Should have parse successfully but got error: %q", err)
			}

			if cltDisconnect.Reason != testData.expectedReason {
				t.Errorf("Should have received reason '%q' but got reason '%q'.", testData.expectedReason, cltDisconnect.Reason)
			}
		})
	}
}

func Test_parseClientMessage(t *testing.T) {
	goodDatum := []struct {
		rawMsg              string
		expectedUsername    string
		expectedServerSlot  string
		expectedSteamID     string
		expectedAffiliation string
		expectedMsg         string
	}{
		// CSGO
		{`"Console<0><Console><Console>" say "WarMod [BFG] WarmUp Config Loaded"`, "Console", "0", "Console", "Console", `say "WarMod [BFG] WarmUp Config Loaded"`},
		{`"&#0000106&#0000097<31><STEAM_1:0:59942879>" switched from team <Unassigned> to <TERRORIST>`, "&#0000106&#0000097", "31", "STEAM_1:0:59942879", "", `switched from team <Unassigned> to <TERRORIST>`},
		{`"onfocus=JaVaSCript:alert(123)<4><STEAM_1:1:8643911><CT>" purchased "p90"`, "onfocus=JaVaSCript:alert(123)", "4", "STEAM_1:1:8643911", "CT", `purchased "p90"`},
		{`"💘 | LacledesLAN.com<8><STEAM_1:0:99144862><CT>" money change 16000-1000 = $15000 (tracked) (purchase: item_assaultsuit)`, "💘 | LacledesLAN.com", "8", "STEAM_1:0:99144862", "CT", `money change 16000-1000 = $15000 (tracked) (purchase: item_assaultsuit)`},
		// TF2
		{`"Ω≈ç√∫˜µ≤≥÷(｡◕‿◕｡)<2><[U:1:6107481]><>" connected, address "192.168.1.210:27005"`, "Ω≈ç√∫˜µ≤≥÷(｡◕‿◕｡)", "2", "U:1:6107481", "", `connected, address "192.168.1.210:27005"`},
		{`"panzershrek<11><[U:1:122465451]><>" connected, address "192.168.1.37:27005"`, "panzershrek", "11", "U:1:122465451", "", `connected, address "192.168.1.37:27005"`},
		{`"♥♥©♥♥௸atswood<42><13><[U:1:28500804]><>" entered the game`, "♥♥©♥♥௸atswood<42>", "13", "U:1:28500804", "", `entered the game`},
		{`"xXx360noscopesxXx<7><[U:1:122465451]><Unassigned>" joined team "Blue"`, "xXx360noscopesxXx", "7", "U:1:122465451", "Unassigned", `joined team "Blue"`},
		{`"r***yDestroyer9<14><[U:1:122465451]><Blue>" changed role to "sniper"`, "r***yDestroyer9", "14", "U:1:122465451", "Blue", `changed role to "sniper"`},
		{`"<LL>arcticfox012<5><[U:1:5015550]><Red>" killed "[LL]rnjmur<6><[U:1:13251124]><Blue>" with "minigun" (attacker_position "608 -871 -234") (victim_position "596 -532 -261")`, "<LL>arcticfox012", "5", "U:1:5015550", "Red", `killed "[LL]rnjmur<6><[U:1:13251124]><Blue>" with "minigun" (attacker_position "608 -871 -234") (victim_position "596 -532 -261")`},
	}

	for _, testData := range goodDatum {
		t.Run(testData.rawMsg, func(t *testing.T) {
			clientMsg, err := parseClientMessage(LogEntry{Message: testData.rawMsg})

			if err != nil {
				t.Errorf("Should have parsed successfully but got error: %q", err)
			}

			if clientMsg.Client.Username != testData.expectedUsername {
				t.Errorf("Should have received Username '%q' not '%q'.", testData.expectedUsername, clientMsg.Client.Username)
			}

			if clientMsg.Client.ServerSlot != testData.expectedServerSlot {
				t.Errorf("Should have received ServerSlot '%q' not '%q'.", testData.expectedServerSlot, clientMsg.Client.ServerSlot)
			}

			if clientMsg.Client.SteamID != testData.expectedSteamID {
				t.Errorf("Should have received SteamID '%q' not '%q'.", testData.expectedSteamID, clientMsg.Client.SteamID)
			}

			if clientMsg.Client.Affiliation != testData.expectedAffiliation {
				t.Errorf("Should have received Affiliation '%q' not '%q'.", testData.expectedAffiliation, clientMsg.Client.Affiliation)
			}

			if clientMsg.Message != testData.expectedMsg {
				t.Errorf("Should have received message '%q' not '%q'.", testData.expectedMsg, clientMsg.Message)
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
