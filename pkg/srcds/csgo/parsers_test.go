package csgo

import (
	"testing"

	"github.com/LacledesLAN/sourceseer/pkg/srcds"
)

func Test_parseAffiliation(t *testing.T) {
	validCases := []struct {
		input    string
		expected affiliation
	}{
		{``, unassigned},
		{`ct `, counterterrorist},
		{` CT`, counterterrorist},
		{` terrorist`, terrorist},
		{`Terrorist `, terrorist},
		{`TERRORIST`, terrorist},
		{`unassigned `, unassigned},
		{` Unassigned`, unassigned},
		{`UNASSIGNED`, unassigned},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, test := range validCases {
			if actual, ok := parseAffiliation(test.input); !ok {
				t.Errorf("String affiliation %q should have successfully parsed.", test.input)
			} else if actual != test.expected {
				t.Errorf("Expected affiliation %q but got %q.", test.expected, actual)
			}
		}
	})

	invalidCases := []string{
		`I’m so embarrassed. I wish everybody else was dead.`,
	}

	t.Run("Invalid Cases", func(t *testing.T) {
		for _, input := range invalidCases {
			if _, ok := parseAffiliation(input); ok {
				t.Errorf("String affiliation %q should NOT have successfully parsed.", input)
			}
		}
	})
}

func Test_parseClientSay(t *testing.T) {
	mockClient := srcds.Client{
		Username:    "AA",
		SteamID:     "BB",
		ServerSlot:  "CC",
		Affiliation: "DD",
	}

	validCases := []struct {
		rawMsg          string
		expectedMsg     string
		expectedChannel sayChannel
	}{
		{`say ""running server.cfg""`, `"running server.cfg"`, ChannelGlobal},
		{`say "1 PING GUYS"`, "1 PING GUYS", ChannelGlobal},
		{`say_team "save rush?"`, "save rush?", ChannelAffiliation},
		{`say_team "mid rush"`, "mid rush", ChannelAffiliation},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, test := range validCases {
			if actual, ok := parseClientSay(srcds.ClientLogEntry{Client: mockClient, Message: test.rawMsg}); !ok {
				t.Errorf("Message %q should have successfully parsed.", test.rawMsg)
			} else {
				if actual.channel != test.expectedChannel {
					t.Errorf("Expected client say channel %q but got %q", test.expectedChannel, actual.channel)
				}

				if actual.msg != test.expectedMsg {
					t.Errorf("Expected messsage %q but got %q.", test.expectedMsg, actual.msg)
				}
			}
		}
	})

	invalidCases := []string{
		``,
		`          `,
		`dropped "knife"`,
		`picked up "knife"`,
		`triggered "Got_The_Bomb"`,
		`left buyzone with [ weapon_knife_karambit kevlar(98) ]`,
		`attacked "BigBop Lil' Bop<14><STEAM_1:1:32971431><TERRORIST>" [-1406 221 -60] with "knife" (damage "30") (damage_armor "2") (health "57") (armor "95") (hitgroup "generic")`,
		`disconnected`,
		`switched from team <Unassigned> to <TERRORIST>`, "&#0000106&#0000097", "31", "STEAM_1:0:59942879", "", `switched from team <Unassigned> to <TERRORIST>`,
		`money change 16000-1000 = $15000 (tracked) (purchase: item_assaultsuit)`,
		`connected, address ""`,
	}

	t.Run("Invalid Cases", func(t *testing.T) {
		for _, msg := range invalidCases {
			if _, ok := parseClientSay(srcds.ClientLogEntry{Client: mockClient, Message: msg}); ok {
				t.Errorf("Message %q should NOT have successfully parsed.", msg)
			}
		}
	})
}

func Test_parseClientSetAffiliation(t *testing.T) {
	mockClient := srcds.Client{
		Username:    "AA",
		SteamID:     "BB",
		ServerSlot:  "CC",
		Affiliation: "DD",
	}

	validCases := []struct {
		msg          string
		expectedFrom affiliation
		expectedTo   affiliation
	}{
		{"switched from team <CT> to <CT>", counterterrorist, counterterrorist},
		{"switched from team <CT> to <TERRORIST>", counterterrorist, terrorist},
		{"switched from team <CT> to <Unassigned>", counterterrorist, unassigned},
		{"switched from team <TERRORIST> to <CT>", terrorist, counterterrorist},
		{"switched from team <TERRORIST> to <TERRORIST>", terrorist, terrorist},
		{"switched from team <TERRORIST> to <Unassigned>", terrorist, unassigned},
		{"switched from team <Unassigned> to <CT>", unassigned, counterterrorist},
		{"switched from team <Unassigned> to <TERRORIST>", unassigned, terrorist},
		{"switched from team <Unassigned> to <Unassigned>", unassigned, unassigned},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, test := range validCases {
			if actual, ok := parseClientSetAffiliation(srcds.ClientLogEntry{Client: mockClient, Message: test.msg}); !ok {
				t.Errorf("Message %q should have successfully parsed.", test.msg)
			} else {
				if actual.to != test.expectedTo {
					t.Errorf("Expected %q not %q for 'to' from message %q.", test.expectedTo, actual.to, test.msg)
				}

				if actual.from != test.expectedFrom {
					t.Errorf("Expected %q not %q for 'from' from message %q.", test.expectedFrom, actual.from, test.msg)
				}
			}
		}
	})

	invalidCases := []string{
		``,
		`         `,
		`dropped "knife"`,
		`picked up "knife"`,
		`triggered "Got_The_Bomb"`,
		`left buyzone with [ weapon_knife_karambit kevlar(98) ]`,
		`attacked "BigBop Lil' Bop<14><STEAM_1:1:32971431><TERRORIST>" [-1406 221 -60] with "knife" (damage "30") (damage_armor "2") (health "57") (armor "95") (hitgroup "generic")`,
		`disconnected`,
		`money change 16000-1000 = $15000 (tracked) (purchase: item_assaultsuit)`,
		`connected, address ""`,
		`say "1 PING BOYS"`,
		`say_team "save rush?"`,
	}

	t.Run("Invalid Cases", func(t *testing.T) {
		for _, test := range invalidCases {
			if _, ok := parseClientSetAffiliation(srcds.ClientLogEntry{Client: mockClient, Message: test}); ok {
				t.Errorf("Message %q should NOT have successfully parsed.", test)
			}
		}
	})
}

func Test_parseGameOver(t *testing.T) {
	validCases := []struct {
		msg            string
		expectedMode   string
		expectedMap    string
		expectedScore1 int
		expectedScore2 int
		expectedTime   int
	}{
		{"Game Over: competitive de_nuke score 18:19 after 34 min", "competitive", "de_nuke", 18, 19, 34},
		{"Game Over: competitive  de_nuke score 18:19 after 34 min", "competitive", "de_nuke", 18, 19, 34},
		{"Game Over: competitive de_inferno score 3:16 after 38 min", "competitive", "de_inferno", 3, 16, 38},
		{"Game Over: competitive  de_inferno score 3:16 after 38 min", "competitive", "de_inferno", 3, 16, 38},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, test := range validCases {
			if actual, ok := parseGameOver(srcds.LogEntry{Message: test.msg}); !ok {
				t.Errorf("Message %q should have successfully parsed.", test.msg)
			} else {
				if actual.mode != test.expectedMode {
					t.Errorf("Expected mode %q but got %q from message %q.", test.expectedMode, actual.mode, test.msg)
				}

				if actual.mapName != test.expectedMap {
					t.Errorf("Expected map %q but got %q from message %q.", test.expectedMap, actual.mapName, test.msg)
				}

				if actual.score1 != test.expectedScore1 {
					t.Errorf("Expected score 1 %q but got %q from message %q.", test.expectedScore1, actual.score1, test.msg)
				}

				if actual.score2 != test.expectedScore2 {
					t.Errorf("Expected score 2 %q but got %q from message %q.", test.expectedScore2, actual.score2, test.msg)
				}

				if actual.minutesElapsed != test.expectedTime {
					t.Errorf("Expected minutes elapsed %q but got %q from message %q.", test.expectedTime, actual.minutesElapsed, test.msg)
				}
			}
		}
	})

	invalidCases := []string{
		``,
		`     `,
		`"r_VehicleViewDampen" = "1"`,
		`server_cvar: "sv_competitive_official_5v5" "1"`,
		`"GOTV<2><BOT><Unassigned>" changed name to "zLLTV_CSGO_BRACKET_05"`,
		`rcon from "172.30.40.10:49932": command "logaddress_add "172.30.40.10:7130""`,
		`"^v^AustinPowers<41><STEAM_1:1:24945146><TERRORIST>" left buyzone with [ ]`,
		`"^v^AustinPowers<41><STEAM_1:1:24945146>" switched from team <Unassigned> to <CT>`,
		`"^v^AustinPowers<41><STEAM_1:1:24945146><Unassigned>" triggered "clantag" (value "")`,
		`Team playing "CT": HELLO`,
	}

	t.Run("Invalid Cases", func(t *testing.T) {
		for _, msg := range invalidCases {
			if _, ok := parseGameOver(srcds.LogEntry{Message: msg}); ok {
				t.Errorf("Message %q should NOT have successfully parsed.", msg)
			}
		}
	})
}

func Test_parseLoadingMap(t *testing.T) {
	validCases := []struct {
		msg         string
		expectedMap string
	}{
		{`Loading map "de_nuke"`, "de_nuke"},
		{`Loading map "de_dust2"`, "de_dust2"},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, test := range validCases {
			if actualMap, ok := parseLoadingMap(srcds.LogEntry{Message: test.msg}); !ok {
				t.Errorf("Message %q should have successfully parsed.", test.msg)
			} else if actualMap != test.expectedMap {
				t.Errorf("Expected map %q but got map %q from message %q.", test.expectedMap, actualMap, test.msg)
			}
		}
	})

	invalidCases := []string{
		``,
		`        `,
		`"mp_tournament" = "0"`,
		`server_cvar: "mp_roundtime" "1.9167"`,
		`rcon from "192.168.1.107:61968": command "echo HLSW: Test"`,
		`"Console<0><Console><Console>" say "WarMod [BFG] WarmUp Config Loaded"`,
		`"stiff 17<7><STEAM_1:2:66421616><CT>" triggered "clantag" (value "man from uncle")`,
		`Team "CT" scored "0" with "5" players`,
		`Team "TERRORIST" scored "0" with "5" players`,
		`Team "CT" triggered "SFUI_Notice_Bomb_Defused" (CT "21") (T "7")`,
		`Team "TERRORIST" triggered "SFUI_Notice_Bomb_Defused" (CT "21") (T "7")`,
		`World triggered "Round_Start"`,
		`World triggered "Round_End"`,
	}

	t.Run("Invalid Cases", func(t *testing.T) {
		for _, msg := range invalidCases {
			if _, ok := parseLoadingMap(srcds.LogEntry{Message: msg}); ok {
				t.Errorf("Message %q should NOT have successfully parsed.", msg)
			}
		}
	})
}

func Test_parseStartingFreezePeriod(t *testing.T) {
	t.Run("Valid Cases", func(t *testing.T) {
		validCases := []string{
			`Starting Freeze period`,
		}

		for _, s := range validCases {
			if !parseStartingFreezePeriod(srcds.LogEntry{Message: s}) {
				t.Errorf("Message %q should have successfully parsed.", s)
			}
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := []string{
			``,
			`      `,
			`Kittens give Morbo gas.`,
		}

		for _, s := range invalidCases {
			if parseStartingFreezePeriod(srcds.LogEntry{Message: s}) {
				t.Errorf("Message %q should NOT have successfully parsed.", s)
			}
		}
	})
}

func Test_parseTeamScored(t *testing.T) {
	validCases := []struct {
		msg                 string
		expectedAffiliation affiliation
		expectedScore       int
		expectedPlayerCount int
	}{
		{`Team "CT" scored "0" with "5" players`, counterterrorist, 0, 5},
		{`Team "CT" scored "2" with "121" players`, counterterrorist, 2, 121},
		{`Team "CT" scored "999" with "999" players`, counterterrorist, 999, 999},
		{`Team "CT" scored "9999" with "9999" players`, counterterrorist, 9999, 9999},
		{`Team "TERRORIST" scored "0" with "5" players`, terrorist, 0, 5},
		{`Team "TERRORIST" scored "2" with "121" players`, terrorist, 2, 121},
		{`Team "TERRORIST" scored "999" with "999" players`, terrorist, 999, 999},
		{`Team "TERRORIST" scored "9999" with "9999" players`, terrorist, 9999, 9999},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, test := range validCases {
			if actual, ok := parseTeamScored(srcds.LogEntry{Message: test.msg}); !ok {
				t.Errorf("Message %q should have successfully parsed.", test.msg)
			} else {
				if actual.affiliation != test.expectedAffiliation {
					t.Errorf("Expected affiliation %q but got %q from message %q.", test.expectedAffiliation, actual.affiliation, test.msg)
				}

				if actual.Score != test.expectedScore {
					t.Errorf("Expected score %d but got %d from message %q.", test.expectedScore, actual.Score, test.msg)
				}

				if actual.PlayerCount != test.expectedPlayerCount {
					t.Errorf("Expected player count %d but got %d from message %q.", test.expectedPlayerCount, actual.PlayerCount, test.msg)
				}
			}
		}
	})

	invalidCases := []string{
		``,
		`          `,
		`Team "" scored "0" with "5" players`,
		`Team " " scored "60" with "60" players`,
		`Team "CT" triggered "SFUI_Notice_Bomb_Defused" (CT "21") (T "7")`,
		`Team "TERRORIST" triggered "SFUI_Notice_Bomb_Defused" (CT "21") (T "7")`,
		`Team "CT" triggered "SFUI_Notice_CTs_Win" (CT "124") (T "0")`,
		`Team "TERRORIST" triggered "SFUI_Notice_CTs_Win" (CT "124") (T "0")`,
		`Team "CT" triggered "SFUI_Notice_Terrorists_Win" (CT "0") (T "1")`,
		`Team "TERRORIST" triggered "SFUI_Notice_Terrorists_Win" (CT "0") (T "1")`,
	}

	t.Run("Invalid Cases", func(t *testing.T) {
		for _, msg := range invalidCases {
			if _, ok := parseTeamScored(srcds.LogEntry{Message: msg}); ok {
				t.Errorf("Message %q should NOT have successfully parsed.", msg)
			}
		}
	})
}

func Test_parseTeamSetName(t *testing.T) {
	validCases := []struct {
		msg          string
		expectedSide affiliation
		expectedName string
	}{
		{`Team playing "CT": New New Yorkers`, counterterrorist, "New New Yorkers"},
		{`Team playing "TERRORIST": Thunder Cougar Falcon Birds`, terrorist, "Thunder Cougar Falcon Birds"},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, test := range validCases {
			if actual, ok := parseTeamSetName(srcds.LogEntry{Message: test.msg}); !ok {
				t.Errorf("Message %q should have successfully parsed.", test.msg)
			} else {
				if actual.affiliation != test.expectedSide {
					t.Errorf("Expected affiliation %q but got %q from message %q.", test.expectedSide, actual.affiliation, test.msg)
				}

				if actual.teamName != test.expectedName {
					t.Errorf("Expected team name %q but got %q from message %q.", test.expectedName, actual.teamName, test.msg)
				}
			}
		}
	})

	invalidCases := []string{
		``,
		`         `,
		`Team "" scored "0" with "5" players`,
		`Team " " scored "60" with "60" players`,
		`Team "CT" triggered "SFUI_Notice_Bomb_Defused" (CT "21") (T "7")`,
		`Team "TERRORIST" triggered "SFUI_Notice_Bomb_Defused" (CT "21") (T "7")`,
		`Team "CT" triggered "SFUI_Notice_CTs_Win" (CT "124") (T "0")`,
		`Team "TERRORIST" triggered "SFUI_Notice_CTs_Win" (CT "124") (T "0")`,
		`Team "CT" triggered "SFUI_Notice_Terrorists_Win" (CT "0") (T "1")`,
		`Team "TERRORIST" triggered "SFUI_Notice_Terrorists_Win" (CT "0") (T "1")`,
		`Team "CT" scored "0" with "5" players`,
		`Team "TERRORIST" scored "0" with "5" players`,
	}

	t.Run("Invalid Cases", func(t *testing.T) {
		for _, msg := range invalidCases {
			if _, ok := parseTeamSetName(srcds.LogEntry{Message: msg}); ok {
				t.Errorf("Message %q should NOT have successfully parsed.", msg)
			}
		}
	})
}

func Test_parseTeamTriggered(t *testing.T) {
	validCases := []struct {
		msg                 string
		expectedAffiliation affiliation
		expectedTrigger     string
		expectedCTScore     int
		expectedTScore      int
	}{
		{`Team "CT" triggered "SFUI_Notice_Bomb_Defused" (CT "21") (T "7")`, counterterrorist, "SFUI_Notice_Bomb_Defused", 21, 7},
		{`Team "CT" triggered "SFUI_Notice_CTs_Win" (CT "124") (T "0")`, counterterrorist, "SFUI_Notice_CTs_Win", 124, 0},
		{`Team "CT" triggered "SFUI_Notice_Target_Saved" (CT "12") (T "3")`, counterterrorist, "SFUI_Notice_Target_Saved", 12, 3},
		{`Team "TERRORIST" triggered "SFUI_Notice_Target_Bombed" (CT "0") (T "5")`, terrorist, "SFUI_Notice_Target_Bombed", 0, 5},
		{`Team "TERRORIST" triggered "SFUI_Notice_Terrorists_Win" (CT "6") (T "23")`, terrorist, "SFUI_Notice_Terrorists_Win", 6, 23},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, test := range validCases {
			if actual, ok := parseTeamTriggered(srcds.LogEntry{Message: test.msg}); !ok {
				t.Errorf("Message %q should have successfully parsed.", test.msg)
			} else {
				if actual.affiliation != test.expectedAffiliation {
					t.Errorf("Expected affiliation %q but got %q from message %q.", test.expectedAffiliation, actual.affiliation, test.msg)
				}

				if actual.trigger != test.expectedTrigger {
					t.Errorf("Expected trigger %q but got %q from message %q.", test.expectedTrigger, actual.trigger, test.msg)
				}

				if actual.ctScore != test.expectedCTScore {
					t.Errorf("Expected CT score %d but got %d from message %q.", test.expectedCTScore, actual.ctScore, test.msg)
				}

				if actual.terroristScore != test.expectedTScore {
					t.Errorf("Expected Terrorist score %d but got %d from message %q.", test.expectedTScore, actual.terroristScore, test.msg)
				}
			}
		}
	})

	invalidCases := []string{
		``,
		`           `,
		`Team "" scored "0" with "5" players`,
		`Team " " scored "60" with "60" players`,
		`Team "CT" scored "0" with "5" players`,
		`Team "TERRORIST" scored "0" with "5" players`,
		`Team playing "CT": New New Yorkers`,
		`Team playing "TERRORIST": Thunder Cougar Falcon Birds`,
	}

	t.Run("Invalid Cases", func(t *testing.T) {
		for _, msg := range invalidCases {
			if _, ok := parseTeamTriggered(srcds.LogEntry{Message: msg}); ok {
				t.Errorf("Message %q should NOT have successfully parsed.", msg)
			}
		}
	})
}

func Test_parseWorldTrigger(t *testing.T) {
	validCases := []struct {
		msg      string
		expected string
	}{
		{`World triggered "Game_Commencing"`, `"Game_Commencing"`},
		{`World triggered "Match_Start" on "de_mirage"`, `"Match_Start" on "de_mirage"`},
		{`World triggered "Round_Start"`, `"Round_Start"`},
		{`World triggered "Round_End"`, `"Round_End"`},
		{`World triggered "Restart_Round_(1_second)"`, `"Restart_Round_(1_second)"`},
		{`World triggered "Restart_Round_(3_seconds)"`, `"Restart_Round_(3_seconds)"`},
		{`World triggered "Restart_Round_(99_seconds)"`, `"Restart_Round_(99_seconds)"`},
		{`World triggered "Restart_Round_(999_seconds)"`, `"Restart_Round_(999_seconds)"`},
		{`World triggered "Restart_Round_(9999_seconds)"`, `"Restart_Round_(9999_seconds)"`},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, test := range validCases {
			if actual, ok := parseWorldTrigger(srcds.LogEntry{Message: test.msg}); !ok {
				t.Errorf("Message %q should have successfully parsed.", test.msg)
			} else if string(actual) != test.expected {
				t.Errorf("Expected %q but got %q from message %q.", test.expected, string(actual), test.msg)
			}
		}
	})

	invalidCases := []string{
		``,
		`        `,
		`"r_VehicleViewDampen" = "1"`,
		`server_cvar: "sv_competitive_official_5v5" "1"`,
		`"GOTV<2><BOT><Unassigned>" changed name to "zLLTV_CSGO_BRACKET_05"`,
		`rcon from "172.30.40.10:49932": command "logaddress_add "172.30.40.10:7130""`,
		`"^v^AustinPowers<41><STEAM_1:1:24945146><TERRORIST>" left buyzone with [ ]`,
		`"^v^AustinPowers<41><STEAM_1:1:24945146>" switched from team <Unassigned> to <CT>`,
		`"^v^AustinPowers<41><STEAM_1:1:24945146><Unassigned>" triggered "clantag" (value "")`,
		`Team playing "CT": HELLO`,
	}

	t.Run("Invalid Cases", func(t *testing.T) {
		for _, msg := range invalidCases {
			if _, ok := parseWorldTrigger(srcds.LogEntry{Message: msg}); ok {
				t.Errorf("Message %q should NOT have successfully parsed.", msg)
			}
		}
	})
}

func Test_parseWorldTriggerGameCommencing(t *testing.T) {
	t.Run("Valid Cases", func(t *testing.T) {
		validCases := []logWorldTrigger{
			`"Game_Commencing"`,
		}

		for _, s := range validCases {
			if !parseWorldTriggerGameCommencing(s) {
				t.Errorf("Message %q should have successfully parsed.", s)
			}
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := []logWorldTrigger{
			``,
			`      `,
			`Kittens give Morbo gas.`,
		}

		for _, s := range invalidCases {
			if parseWorldTriggerGameCommencing(s) {
				t.Errorf("Message %q should NOT have successfully parsed.", s)
			}
		}
	})
}

func Test_parseWorldTriggerRoundEnd(t *testing.T) {
	t.Run("Valid Cases", func(t *testing.T) {
		validCases := []logWorldTrigger{
			`"Round_End"`,
		}

		for _, s := range validCases {
			if !parseWorldTriggerRoundEnd(s) {
				t.Errorf("Message %q should have successfully parsed.", s)
			}
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := []logWorldTrigger{
			``,
			`      `,
			`Kittens give Morbo gas.`,
		}

		for _, s := range invalidCases {
			if parseWorldTriggerRoundEnd(s) {
				t.Errorf("Message %q should NOT have successfully parsed.", s)
			}
		}
	})
}

func Test_parseWorldTriggerRoundStart(t *testing.T) {
	t.Run("Valid Cases", func(t *testing.T) {
		validCases := []logWorldTrigger{
			`"Round_Start"`,
		}

		for _, s := range validCases {
			if !parseWorldTriggerRoundStart(s) {
				t.Errorf("Message %q should have successfully parsed.", s)
			}
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := []logWorldTrigger{
			``,
			`      `,
			`Kittens give Morbo gas.`,
		}

		for _, s := range invalidCases {
			if parseWorldTriggerRoundStart(s) {
				t.Errorf("Message %q should NOT have successfully parsed.", s)
			}
		}
	})
}

func Test_parseWorldTriggerMatchStart(t *testing.T) {
	validCases := []struct {
		msg      logWorldTrigger
		expected string
	}{
		{`"Match_Start" on "de_dust2"`, `de_dust2`},
		{`"Match_Start" on "de_inferno"`, `de_inferno`},
		{`"Match_Start" on "de_mirage"`, `de_mirage`},
		{`"Match_Start" on "de_nuke"`, `de_nuke`},
		{`"Match_Start" on "de_overpass"`, `de_overpass`},
		{`"Match_Start" on "de_train"`, `de_train`},
		{`"Match_Start" on "de_vertigo"`, `de_vertigo`},
		{`"Match_Start" on "dm_zz_top"`, `dm_zz_top`},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, test := range validCases {
			if actual, ok := parseWorldTriggerMatchStart(test.msg); !ok {
				t.Errorf("Message %q should have successfully parsed.", test.msg)
			} else if actual != test.expected {
				t.Errorf("Expected map %q but got %q from message %q.", test.expected, actual, test.msg)
			}
		}
	})

	invalidCases := []logWorldTrigger{
		``,
		`    `,
		`"Game_Commencing"`,
		`"Round_Start"`,
		`"Round_End"`,
		`"Restart_Round_(1_second)"`,
		`"Restart_Round_(22_second)"`,
		`"Restart_Round_(333_second)"`,
		`"Restart_Round_(4444_second)"`,
	}

	t.Run("Invalid Cases", func(t *testing.T) {
		for _, msg := range invalidCases {
			if _, ok := parseWorldTriggerMatchStart(msg); ok {
				t.Errorf("Message %q should NOT have successfully parsed.", msg)
			}
		}
	})
}
