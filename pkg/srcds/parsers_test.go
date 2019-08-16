package srcds

import (
	"testing"
	"time"
)

func Test_ParseClient(t *testing.T) {
	t.Run("Valid Cases", func(t *testing.T) {
		validCases := map[string][]struct {
			msg                 string
			expectedUsername    string
			expectedServerSlot  int16
			expectedSteamID     string
			expectedAffiliation string
		}{
			"CSGO": {
				{`"Console<0><Console><Console>"`, "Console", 0, "Console", "Console"},
				{`"&#0000106&#0000097<31><STEAM_1:0:59942879>"`, "&#0000106&#0000097", 31, "STEAM_1:0:59942879", ""},
				{`"onfocus=JaVaSCript:alert(123)<4><STEAM_1:1:8643911><CT>"`, "onfocus=JaVaSCript:alert(123)", 4, "STEAM_1:1:8643911", "CT"},
				{`"💘 | LacledesLAN.com<8><STEAM_1:0:99144862><CT>"`, "💘 | LacledesLAN.com", 8, "STEAM_1:0:99144862", "CT"},
			},
			"TF2": {
				{`"Ω≈ç√∫˜µ≤≥÷(｡◕‿◕｡)<2><[U:1:6107481]><>"`, "Ω≈ç√∫˜µ≤≥÷(｡◕‿◕｡)", 2, "U:1:6107481", ""},
				{`"panzershrek<11><[U:1:122465451]><>"`, "panzershrek", 11, "U:1:122465451", ""},
				{`"♥♥©♥♥௸atswood<42><13><[U:1:28500804]><>"`, "♥♥©♥♥௸atswood<42>", 13, "U:1:28500804", ""},
				{`"xXx360noscopesxXx<7><[U:1:122465451]><Unassigned>"`, "xXx360noscopesxXx", 7, "U:1:122465451", "Unassigned"},
				{`"r***yDestroyer9<14><[U:1:122465451]><Blue>"`, "r***yDestroyer9", 14, "U:1:122465451", "Blue"},
				{`"<LL>arcticfox012<5><[U:1:5015550]><Red>"`, "<LL>arcticfox012", 5, "U:1:5015550", "Red"},
				{`"GutsAndGlory!<16><BOT><Red>"`, "GutsAndGlory!", 16, "BOT", "Red"},
			},
		}

		for name, tests := range validCases {
			t.Run(name, func(t *testing.T) {
				for _, test := range tests {
					if actual, ok := ParseClient(test.msg); !ok {
						t.Errorf("Message %q should have successfully parsed.", test.msg)
					} else {
						if actual.Affiliation != test.expectedAffiliation {
							t.Errorf("Expected Affiliation %q but got %q.", test.expectedAffiliation, actual.Affiliation)
						}

						if actual.SteamID != test.expectedSteamID {
							t.Errorf("Expected SteamID %q but got %q.", test.expectedSteamID, actual.SteamID)
						}

						if actual.ServerSlot != test.expectedServerSlot {
							t.Errorf("Expected ServerSlot %q but got %q.", test.expectedServerSlot, actual.ServerSlot)
						}

						if actual.Username != test.expectedUsername {
							t.Errorf("Expected Username %q but got %q.", test.expectedUsername, actual.Username)
						}
					}
				}
			})
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := map[string][]string{
			"CSGO": {
				`Loading map "de_cache"`,
				`server cvars start`,
				`World triggered "Round_Start"`,
			},
			"TF2": {
				`server_cvar: "sv_alltalk" "1"`,
				`Connection to Steam servers successful.`,
				`Assigned anonymous gameserver Steam ID [A:1:2271993858:11261].`,
				`Team "Red" triggered "pointcaptured" (cp "0") (cpname "#koth_viaduct_cap") (numcappers "1") (player1 "BEan | LacledesLAN.com<2><[U:1:7609438]><Red>") (position1 "440 319 -172") `,
				`Team "Red" current score "0" with "3" players`,
				`World triggered "Round_Overtime"`,
				`World triggered "Round_Win" (winner "Blue")`,
				`[META] Loaded 0 plugins (1 already loaded)`,
			},
			"OTHER": {
				``,
				`       `,
				`With a warning label this big, you know they gotta be fun!`,
			},
		}

		for name, tests := range invalidCases {
			t.Run(name, func(t *testing.T) {
				for _, msg := range tests {
					if _, ok := ParseClient(msg); ok {
						t.Errorf("Message %q should NOT have successfully parsed.", msg)
					}
				}
			})
		}
	})
}

func Test_ParseClientLogEntry(t *testing.T) {
	t.Run("Valid Cases", func(t *testing.T) {
		validCases := map[string][]struct {
			rawMsg              string
			expectedUsername    string
			expectedServerSlot  int16
			expectedSteamID     string
			expectedAffiliation string
			expectedMsg         string
		}{
			"CSGO": {
				{`"Charles Nicole<1><STEAM_1:0:13377331><>" connected, address ""`, `Charles Nicole`, 1, "STEAM_1:0:13377331", "", `connected, address ""`},
				{`"Bill Nye<2><STEAM_1:0:13377331><>" STEAM USERID validated`, `Bill Nye`, 2, "STEAM_1:0:13377331", "", `STEAM USERID validated`},
				{`"Thomas Gold<3><STEAM_1:0:13377331><>" entered the game`, "Thomas Gold", 3, "STEAM_1:0:13377331", ``, `entered the game`},
				{`"Jane Goodall<4><STEAM_1:0:13377331><CT>" purchased "m4a1"`, "Jane Goodall", 4, "STEAM_1:0:13377331", `CT`, `purchased "m4a1"`},
				{`"Charles Babbage<5><STEAM_1:0:13377331><TERRORIST>" say "Where's sunny"`, `Charles Babbage`, 5, "STEAM_1:0:13377331", "TERRORIST", `say "Where's sunny"`},
				{`"B. Pascal<6><STEAM_1:0:13377331><CT>" left buyzone with [ ]`, `B. Pascal`, 6, "STEAM_1:0:13377331", "CT", `left buyzone with [ ]`},
				{`"Joseph Henry<7><STEAM_1:0:13377331><CT>" left buyzone with [ weapon_knife_m9_bayonet weapon_usp_silencer weapon_m4a1 kevlar(100) helmet ]`, `Joseph Henry`, 7, "STEAM_1:0:13377331", "CT", `left buyzone with [ weapon_knife_m9_bayonet weapon_usp_silencer weapon_m4a1 kevlar(100) helmet ]`},
				{`"David B.<8><STEAM_1:0:13377331><Unassigned>" disconnected (reason "David B. timed out")`, `David B.`, 8, "STEAM_1:0:13377331", "Unassigned", `disconnected (reason "David B. timed out")`},
				{`"Sally Ride<9><BOT><CT>" disconnected (reason "Kicked by Console")`, `Sally Ride`, 9, "BOT", "CT", `disconnected (reason "Kicked by Console")`},
				{`"Ronald Ross<10><STEAM_1:0:13377331>" switched from team <Unassigned> to <TERRORIST>`, `Ronald Ross`, 10, "STEAM_1:0:13377331", "", `switched from team <Unassigned> to <TERRORIST>`},
				{`"B. F. Skinner<11><BOT><TERRORIST>" say_team ".r"`, `B. F. Skinner`, 11, "BOT", "TERRORIST", `say_team ".r"`},
			},
			"TF2": {
				{`"Ω≈ç√∫˜µ≤≥÷(｡◕‿◕｡)<2><[U:1:6107481]><>" connected, address "192.168.1.210:27005"`, "Ω≈ç√∫˜µ≤≥÷(｡◕‿◕｡)", 2, "U:1:6107481", "", `connected, address "192.168.1.210:27005"`},
				{`"panzershrek<11><[U:1:122465451]><>" connected, address "192.168.1.37:27005"`, "panzershrek", 11, "U:1:122465451", "", `connected, address "192.168.1.37:27005"`},
				{`"♥♥©♥♥௸atswood<42><13><[U:1:28500804]><>" entered the game`, "♥♥©♥♥௸atswood<42>", 13, "U:1:28500804", "", `entered the game`},
				{`"xXx360noscopesxXx<7><[U:1:122465451]><Unassigned>" joined team "Blue"`, "xXx360noscopesxXx", 7, "U:1:122465451", "Unassigned", `joined team "Blue"`},
				{`"r***yDestroyer9<14><[U:1:122465451]><Blue>" changed role to "sniper"`, "r***yDestroyer9", 14, "U:1:122465451", "Blue", `changed role to "sniper"`},
				{`"<LL>arcticfox012<5><[U:1:5015550]><Red>" killed "[LL]Buddha<6><[U:1:13251124]><Blue>" with "minigun" (attacker_position "608 -871 -234") (victim_position "596 -532 -261")`, "<LL>arcticfox012", 5, "U:1:5015550", "Red", `killed "[LL]Buddha<6><[U:1:13251124]><Blue>" with "minigun" (attacker_position "608 -871 -234") (victim_position "596 -532 -261")`},
				{`"GutsAndGlory!<16><BOT><Red>" changed role to "demoman"`, "GutsAndGlory!", 16, "BOT", "Red", `changed role to "demoman"`},
			},
			"OTHER": {
				{`"Console<0><Console><Console>" say "WarMod [BFG] WarmUp Config Loaded"`, "Console", 0, "Console", "Console", `say "WarMod [BFG] WarmUp Config Loaded"`},
				{`"&#0000106&#0000097<31><STEAM_1:0:59942879>" switched from team <Unassigned> to <Hello There>`, "&#0000106&#0000097", 31, "STEAM_1:0:59942879", "", `switched from team <Unassigned> to <Hello There>`},
				{`"onfocus=JaVaSCript:alert(123)<4><STEAM_1:1:8643911><CT>" purchased "p90"`, "onfocus=JaVaSCript:alert(123)", 4, "STEAM_1:1:8643911", "CT", `purchased "p90"`},
				{`"💘 | LacledesLAN.com<8><STEAM_1:0:99144862><CT>" money change 16000-1000 = $15000 (tracked) (purchase: item_assaultsuit)`, "💘 | LacledesLAN.com", 8, "STEAM_1:0:99144862", "CT", `money change 16000-1000 = $15000 (tracked) (purchase: item_assaultsuit)`},
			},
		}

		for name, tests := range validCases {
			t.Run(name, func(t *testing.T) {
				for _, test := range tests {
					if actual, ok := ParseClientLogEntry(LogEntry{Message: test.rawMsg}); !ok {
						t.Errorf("Raw message %q should have parsed successfully but didn't", test.rawMsg)
					} else {
						if actual.Client.Username != test.expectedUsername {
							t.Errorf("Should have received Username '%q' not '%q'.", test.expectedUsername, actual.Client.Username)
						}

						if actual.Client.ServerSlot != test.expectedServerSlot {
							t.Errorf("Should have received ServerSlot '%q' not '%q'.", test.expectedServerSlot, actual.Client.ServerSlot)
						}

						if actual.Client.SteamID != test.expectedSteamID {
							t.Errorf("Should have received SteamID '%q' not '%q'.", test.expectedSteamID, actual.Client.SteamID)
						}

						if actual.Client.Affiliation != test.expectedAffiliation {
							t.Errorf("Should have received Affiliation '%q' not '%q'.", test.expectedAffiliation, actual.Client.Affiliation)
						}

						if actual.Message != test.expectedMsg {
							t.Errorf("Should have received message '%q' not '%q'.", test.expectedMsg, actual.Message)
						}
					}
				}
			})
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := map[string][]string{
			"CSGO": {
				`Loading map "de_cache"`,
				`server cvars start`,
				`World triggered "Round_Start"`,
			},
			"OTHER": {
				``,
				`      `,
				`I don't want to live on this planet anymore.`,
			},
		}

		for name, tests := range invalidCases {
			t.Run(name, func(t *testing.T) {
				for _, test := range tests {
					if _, ok := ParseClientLogEntry(LogEntry{Message: test}); ok {
						t.Errorf("Log message %q should NOT have successfully parsed.", test)
					}
				}
			})
		}
	})
}

func Test_ParseClientConnected(t *testing.T) {
	mockClient := Client{Username: "Nannybot 1.0", ServerSlot: 1, SteamID: "[U:1:28500804", Affiliation: "Red"}

	t.Run("Valid Cases", func(t *testing.T) {
		validCases := map[string][]string{
			"CSGO": {
				`connected, address ""`,
				`connected, address "192.168.1.52"`,
				`connected, address "2001:0DB8:AC10:FE01:0000:EE91:0000:0000"`,
			},
			"TF2": {
				`connected, address "none"`,
				`connected, address "192.168.1.106:27005"`,
				`connected, address "2001:0DB8:AC10:FE01:0000:EE91:0000:0000"`,
			},
		}

		for name, tests := range validCases {
			t.Run(name, func(t *testing.T) {
				for _, msg := range tests {
					if ok := ParseClientConnected(ClientLogEntry{Client: mockClient, Message: msg}); !ok {
						t.Errorf("Failed to parse valid message %q", msg)
					}
				}
			})
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := map[string][]string{
			"CSGO": []string{
				`say ""running server.cfg""`,
				`purchased "ak47"`,
				`left buyzone with [ weapon_knife_t weapon_glock weapon_ak47 weapon_molotov weapon_smokegrenade weapon_flashbang kevlar(100) helmet ]`,
				`switched from team <Unassigned> to <CT>`,
				`blinded for 0.28 by "Russell<2><STEAM_1:0:165450181><TERRORIST>" from flashbang entindex 293`,
				`threw flashbang [526 712 116] flashbang entindex 293)`,
				`say "might want to tell your friend to swap teams"`,
				`disconnected (reason "haw-ha")`,
				`disconnected`,
				`STEAM USERID validated`,
				`[-429 -1552 -40] attacked "{Apollo}Prime<8><STEAM_1:1:44223295><CT>" [229 -1544 -176] with "negev" (damage "8") (damage_armor "1") (health "92") (armor "98") (hitgroup "chest")`,
				`[503 -1627 -263] killed "{Apollo}Prime<8><STEAM_1:1:44223295><CT>" [350 -1599 -123] with "ak47"`,
			},
			"OTHER": []string{
				``,
				`You don't get to laugh.`,
			},
		}

		for name, tests := range invalidCases {
			t.Run(name, func(t *testing.T) {
				for _, test := range tests {
					if ok := ParseClientConnected(ClientLogEntry{Client: mockClient, Message: test}); ok {
						t.Errorf("Message %q should NOT have successfully parsed.", test)
					}
				}
			})
		}
	})
}

func Test_ParseClientDisconnected(t *testing.T) {
	mockClient := Client{Username: "Sinclair 2K", ServerSlot: 1, SteamID: "[U:1:28500804", Affiliation: ""}

	t.Run("Valid Cases", func(t *testing.T) {
		validCases := map[string][]struct {
			msg            string
			expectedReason string
		}{
			"CSGO": {
				{msg: `disconnected`, expectedReason: ""},
				{msg: `disconnected (reason "David B. Robertson timed out")`, expectedReason: "timed out"},
				{msg: `disconnected (reason "Kicked by Console")`, expectedReason: "Kicked by Console"},
			},
			"TF2": {
				{msg: `disconnected (reason "Kicked from server")`, expectedReason: "Kicked from server"},
				{msg: `disconnected (reason "Server shutting down")`, expectedReason: "Server shutting down"},
			},
		}

		for name, tests := range validCases {
			t.Run(name, func(t *testing.T) {
				for _, test := range tests {
					if actual, ok := ParseClientDisconnected(ClientLogEntry{Client: mockClient, Message: test.msg}); !ok {
						t.Errorf("Failed to parse valid message %q", test.msg)
					} else if string(actual) != test.expectedReason {
						t.Errorf("Expected reason %q but got %q", test.expectedReason, actual)
					}
				}
			})
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := map[string][]string{
			"CSGO": []string{
				`say ""running server.cfg""`,
				`purchased "ak47"`,
				`left buyzone with [ weapon_knife_t weapon_glock weapon_ak47 weapon_molotov weapon_smokegrenade weapon_flashbang kevlar(100) helmet ]`,
				`switched from team <Unassigned> to <CT>`,
				`blinded for 0.28 by "Russell<2><STEAM_1:0:165450181><TERRORIST>" from flashbang entindex 293`,
				`threw flashbang [526 712 116] flashbang entindex 293)`,
				`say "might want to tell your friend to swap teams"`,
				`connected, address ""`,
				`STEAM USERID validated`,
				`[-429 -1552 -40] attacked "{Apollo}Prime<8><STEAM_1:1:44223295><CT>" [229 -1544 -176] with "negev" (damage "8") (damage_armor "1") (health "92") (armor "98") (hitgroup "chest")`,
				`[503 -1627 -263] killed "{Apollo}Prime<8><STEAM_1:1:44223295><CT>" [350 -1599 -123] with "ak47"`,
			},
			"OTHER": []string{
				``,
				`You don't get to laugh.`,
			},
		}

		for name, tests := range invalidCases {
			t.Run(name, func(t *testing.T) {
				for _, test := range tests {
					if _, ok := ParseClientDisconnected(ClientLogEntry{Client: mockClient, Message: test}); ok {
						t.Errorf("Message %q should NOT have successfully parsed.", test)
					}
				}
			})
		}
	})
}

func Test_parseLogEntry(t *testing.T) {
	t.Run("Valid Cases", func(t *testing.T) {
		validCases := map[string][]struct {
			rawLog       string
			expectedTime string
			expectedMsg  string
		}{
			"OTHER": {
				{"L 1/2/2000 - 03:04:00: Sweet llamas of the Bahamas!", "01/02/2000 - 03:04:00", "Sweet llamas of the Bahamas!"},
				{"L 01/2/2000 - 03:04:00: Excuse my language but I have had it with you ruffling my petticoats!", "1/02/2000 - 03:04:00", "Excuse my language but I have had it with you ruffling my petticoats!"},
				{"L 1/02/2000 - 03:04:00: Your music is bad & you should feel bad!", "01/2/2000 - 03:04:00", "Your music is bad & you should feel bad!"},
				{"L 01/02/2000 - 03:04:00: Did everything just taste purple for a second?", "1/2/2000 - 3:04:00", "Did everything just taste purple for a second?"},
				{"L 01/02/2000 - 3:04:00: When you look this good, you don’t have to know anything!", "1/2/2000 - 03:04:00", "When you look this good, you don’t have to know anything!"},
			},
			"CSGO": {
				{`L 01/10/2015 - 10:58:00: Loading map "workshop/163589843/de_cache"`, "01/10/2015 - 10:58:00", `Loading map "workshop/163589843/de_cache"`},
				{`L 01/10/2015 - 10:58:00: server cvars start`, "01/10/2015 - 10:58:00", `server cvars start`},
				{`L 01/10/2015 - 10:58:00: "cash_player_killed_teammate" = "-3300"`, "01/10/2015 - 10:58:00", `"cash_player_killed_teammate" = "-3300"`},
				{`L 10/29/2016 - 19:11:30: Started map "de_cbble" (CRC "787242208")`, "10/29/2016 - 19:11:30", `Started map "de_cbble" (CRC "787242208")`},
				{`L 10/29/2016 - 19:11:30: server_cvar: "cash_team_win_by_defusing_bomb" "3500"`, "10/29/2016 - 19:11:30", `server_cvar: "cash_team_win_by_defusing_bomb" "3500"`},
				{`L 10/29/2016 - 19:11:33: "GOTV<2><BOT><>" connected, address ""`, "10/29/2016 - 19:11:33", `"GOTV<2><BOT><>" connected, address ""`},
				{`L 04/21/2018 - 16:59:16: " - Xerdy  O-O <7><STEAM_1:1:49462758><>" STEAM USERID validated`, "04/21/2018 - 16:59:16", `" - Xerdy  O-O <7><STEAM_1:1:49462758><>" STEAM USERID validated`},
				{`L 10/29/2016 - 18:05:19: "♥♥©♥♥௸atswood<42><13><[U:1:28500804]><>" entered the game`, "10/29/2016 - 18:05:19", `"♥♥©♥♥௸atswood<42><13><[U:1:28500804]><>" entered the game`},
				{`L 10/29/2016 - 18:05:19: "💘 | LacledesLAN.com<8><STEAM_1:0:99144862><CT>" money change 16000-1000 = $15000 (tracked) (purchase: item_assaultsuit)`, "10/29/2016 - 18:05:19", `"💘 | LacledesLAN.com<8><STEAM_1:0:99144862><CT>" money change 16000-1000 = $15000 (tracked) (purchase: item_assaultsuit)`},
				{`L 10/29/2016 - 18:05:31: "Ω≈ç√∫˜µ≤≥÷(｡◕‿◕｡)<2><[U:1:6107481]><>" switched from team <Unassigned> to <TERRORIST>`, "10/29/2016 - 18:05:31", `"Ω≈ç√∫˜µ≤≥÷(｡◕‿◕｡)<2><[U:1:6107481]><>" switched from team <Unassigned> to <TERRORIST>`},
				{`L 10/29/2016 - 18:05:31: World triggered "Game_Commencing"`, "10/29/2016 - 18:05:31", `World triggered "Game_Commencing"`},
				{`L 03/18/2017 - 23:28:08: Starting Freeze period`, "03/18/2017 - 23:28:08", `Starting Freeze period`},
				{`L 03/18/2017 - 23:28:08: World triggered "Match_Start" on "de_inferno"`, "03/18/2017 - 23:28:08", `World triggered "Match_Start" on "de_inferno"`},
				{`L 03/18/2017 - 23:28:08: Team playing "TERRORIST": SOFA_KING`, "03/18/2017 - 23:28:08", `Team playing "TERRORIST": SOFA_KING`},
				{`L 03/18/2017 - 23:28:12: "*<)JuKe <3 BRaTZ<3><STEAM_1:1:3891511><CT>" purchased "vesthelm"`, "03/18/2017 - 23:28:12", `"*<)JuKe <3 BRaTZ<3><STEAM_1:1:3891511><CT>" purchased "vesthelm"`},
				{`L 10/22/2017 - 16:12:41: "Steve<6><STEAM_1:0:22510661><CT>" [-1158 519 -55] attacked "Dan<8><BOT><TERRORIST>" [-1296 401 -56] with "hkp2000" (damage "17") (damage_armor "8") (health "83") (armor "91") (hitgroup "chest")`, "10/22/2017 - 16:12:41", `"Steve<6><STEAM_1:0:22510661><CT>" [-1158 519 -55] attacked "Dan<8><BOT><TERRORIST>" [-1296 401 -56] with "hkp2000" (damage "17") (damage_armor "8") (health "83") (armor "91") (hitgroup "chest")`},
				{`L 10/22/2017 - 16:13:29: "MikeJay @ PARTY<3><STEAM_1:1:72356891><CT>" left buyzone with [ weapon_knife_m9_bayonet weapon_usp_silencer ]`, "10/22/2017 - 16:13:29", `"MikeJay @ PARTY<3><STEAM_1:1:72356891><CT>" left buyzone with [ weapon_knife_m9_bayonet weapon_usp_silencer ]`},
				{`L 10/22/2017 - 16:22:42: World triggered "Round_Start"`, "10/22/2017 - 16:22:42", `World triggered "Round_Start"`},
				{`L 04/21/2018 - 16:58:58: World triggered "Game_Commencing"`, "04/21/2018 - 16:58:58", `World triggered "Game_Commencing"`},
				{`L 04/21/2018 - 16:58:59: World triggered "Match_Start" on "de_mirage"`, "04/21/2018 - 16:58:59", `World triggered "Match_Start" on "de_mirage"`},
				{`L 04/22/2018 - 15:29:21: "BEan [LacledesLAN.com]<3><STEAM_1:1:62160657><TERRORIST>" say "gg"`, "04/22/2018 - 15:29:21", `"BEan [LacledesLAN.com]<3><STEAM_1:1:62160657><TERRORIST>" say "gg"`},
				{`L 04/22/2018 - 15:29:21: "Tonster<4><STEAM_1:1:13023104><CT>" disconnected (reason "Disconnect")`, "04/22/2018 - 15:29:21", `"Tonster<4><STEAM_1:1:13023104><CT>" disconnected (reason "Disconnect")`},
				{`L 04/22/2018 - 15:29:18: Game Over: competitive  de_nuke score 16:9 after 39 min`, "04/22/2018 - 15:29:18", `Game Over: competitive  de_nuke score 16:9 after 39 min`},
				{`L 04/22/2018 - 15:29:18: World triggered "Round_End"`, "04/22/2018 - 15:29:18", `World triggered "Round_End"`},
				{`L 04/22/2018 - 15:29:18: Team "CT" triggered "SFUI_Notice_CTs_Win" (CT "16") (T "9")`, "04/22/2018 - 15:29:18", `Team "CT" triggered "SFUI_Notice_CTs_Win" (CT "16") (T "9")`},
				{`L 04/22/2018 - 15:29:18: Team "CT" scored "16" with "5" players`, "04/22/2018 - 15:29:18", `Team "CT" scored "16" with "5" players`},
				{`L 04/22/2018 - 15:29:18: Team "TERRORIST" scored "9" with "5" players`, "04/22/2018 - 15:29:18", `Team "TERRORIST" scored "9" with "5" players`},
			},
			"TF2": {
				{`L 08/11/2019 - 19:52:33: Log file started (file "logs/L0811004.log") (game "/app/TF2/tf") (version "5257084")`, "08/11/2019 - 19:52:33", `Log file started (file "logs/L0811004.log") (game "/app/TF2/tf") (version "5257084")`},
				{`L 08/11/2019 - 19:52:33: server_cvar: "sv_alltalk" "1"`, "08/11/2019 - 19:52:33", `server_cvar: "sv_alltalk" "1"`},
				{`L 08/11/2019 - 19:52:33: "NotMe<2><BOT><>" connected, address "none"`, "08/11/2019 - 19:52:33", `"NotMe<2><BOT><>" connected, address "none"`},
				{`L 10/21/2018 - 14:08:11: rcon from "172.30.40.7:53493": command "echo HLSW: Test"`, "10/21/2018 - 14:08:11", `rcon from "172.30.40.7:53493": command "echo HLSW: Test"`},
				{`L 10/21/2018 - 14:09:54: "BEan | LacledesLAN.com<2><[U:1:7609438]><>" connected, address "172.30.40.210:27005"`, "10/21/2018 - 14:09:54", `"BEan | LacledesLAN.com<2><[U:1:7609438]><>" connected, address "172.30.40.210:27005"`},
				{`L 08/11/2019 - 19:52:34: "CreditToTeam<3><BOT><Red>" disconnected (reason "Kicked from server")`, "08/11/2019 - 19:52:34", `"CreditToTeam<3><BOT><Red>" disconnected (reason "Kicked from server")`},
				{`L 10/21/2018 - 14:28:54: Team "Red" triggered "pointcaptured" (cp "0") (cpname "#koth_viaduct_cap") (numcappers "2") (player1 "BEan | LacledesLAN.com<2><[U:1:7609438]><Red>") (position1 "661 3 -172") (player2 "[LL]arcticfox012<5><[U:1:5015550]><Red>") (position2 "670 92 -165")`, "10/21/2018 - 14:28:54", `Team "Red" triggered "pointcaptured" (cp "0") (cpname "#koth_viaduct_cap") (numcappers "2") (player1 "BEan | LacledesLAN.com<2><[U:1:7609438]><Red>") (position1 "661 3 -172") (player2 "[LL]arcticfox012<5><[U:1:5015550]><Red>") (position2 "670 92 -165")`},
				{`L 10/21/2018 - 14:30:13: Team "Blue" triggered "pointcaptured" (cp "0") (cpname "#koth_viaduct_cap") (numcappers "2") (player1 "Snek<4><[U:1:122465451]><Blue>") (position1 "718 81 -172") (player2 "[LL]rnjmur<6><[U:1:13251124]><Blue>") (position2 "634 153 -165")`, "10/21/2018 - 14:30:13", `Team "Blue" triggered "pointcaptured" (cp "0") (cpname "#koth_viaduct_cap") (numcappers "2") (player1 "Snek<4><[U:1:122465451]><Blue>") (position1 "718 81 -172") (player2 "[LL]rnjmur<6><[U:1:13251124]><Blue>") (position2 "634 153 -165")`},
				{`L 10/21/2018 - 14:35:17: World triggered "Round_Win" (winner "Red")`, "10/21/2018 - 14:35:17", `World triggered "Round_Win" (winner "Red")`},
				{`L 10/21/2018 - 14:35:17: World triggered "Round_Length" (seconds "592.12")`, "10/21/2018 - 14:35:17", `World triggered "Round_Length" (seconds "592.12")`},
				{`L 10/21/2018 - 14:35:17: Team "Red" current score "1" with "3" players`, "10/21/2018 - 14:35:17", `Team "Red" current score "1" with "3" players`},
				{`L 10/21/2018 - 14:35:17: Team "Blue" current score "1" with "3" players`, "10/21/2018 - 14:35:17", `Team "Blue" current score "1" with "3" players`},
				{`L 10/21/2018 - 14:35:24: "[LL]red<6><[U:1:13251124]><Blue>" changed role to "scout"`, "10/21/2018 - 14:35:24", `"[LL]red<6><[U:1:13251124]><Blue>" changed role to "scout"`},
				{`L 10/21/2018 - 14:39:56: Vote succeeded "NextLevel pl_millstone_event"`, "10/21/2018 - 14:39:56", `Vote succeeded "NextLevel pl_millstone_event"`},
				{`L 10/21/2018 - 14:47:03: World triggered "Game_Over" reason "Reached Time Limit"`, "10/21/2018 - 14:47:03", `World triggered "Game_Over" reason "Reached Time Limit"`},
				{`L 10/21/2018 - 14:47:03: Team "Red" final score "2" with "3" players`, "10/21/2018 - 14:47:03", `Team "Red" final score "2" with "3" players`},
				{`L 10/21/2018 - 14:47:03: Team "Blue" final score "1" with "4" players`, "10/21/2018 - 14:47:03", `Team "Blue" final score "1" with "4" players`},
				{`L 10/21/2018 - 14:47:14: Log file closed.`, "10/21/2018 - 14:47:14", `Log file closed.`},
			},
		}

		for name, tests := range validCases {
			t.Run(name, func(t *testing.T) {
				for _, test := range tests {
					if actual, ok := parseLogEntry(test.rawLog); !ok {
						t.Errorf("Log %q should have successfully parsed.", test.rawLog)
					} else {
						// did the time parse?
						if expectedTimestamp, err := time.ParseInLocation(srcdsTimeLayout, test.expectedTime, time.Local); err != nil {
							t.Fatalf("Test data is bad - couldn't parse timestamp %q", test.expectedTime)
						} else if expectedTimestamp != actual.Timestamp {
							logMessage := test.rawLog
							if len(logMessage) > 56 {
								logMessage = logMessage[:56] + "..."
							}

							t.Errorf("Expected timestamp %q but got %q from %q", expectedTimestamp, actual.Timestamp, logMessage)
						}

						// did the message parse?
						if actual.Message != test.expectedMsg {
							t.Errorf("Expected message of %q but got %q", test.expectedMsg, actual.Message)
						}
					}
				}
			})
		}
	})

	invalidCases := []string{
		``,
		`My first clue came at 4:15, when the clock stopped.`,
	}

	t.Run("Invalid Cases", func(t *testing.T) {
		for _, test := range invalidCases {
			if _, ok := parseLogEntry(test); ok {
				t.Errorf("Raw string %q should NOT have successfully parsed.", test)
			}
		}
	})
}

func Test_parsEchoCvar(t *testing.T) {
	t.Run("Valid Cases", func(t *testing.T) {
		validCases := []struct {
			msg           string
			expectedName  string
			expectedValue string
		}{
			{`sv_stopspeed - 80`, "sv_stopspeed", "80"},
			{`weapon_sound_falloff_multiplier - 1.0`, "weapon_sound_falloff_multiplier", "1.0"},
			{`sv_weapon_encumbrance_per_item - 0.85`, "sv_weapon_encumbrance_per_item", "0.85"},
			{`spec_replay_leadup_time - 5.3438`, "spec_replay_leadup_time", "5.3438"},
			{`mp_t_default_melee - weapon_knife`, "mp_t_default_melee", "weapon_knife"},
			{`mp_global_damage_per_second - 0.0`, "mp_global_damage_per_second", "0.0"},
			{`mp_ct_default_secondary - weapon_hkp2000`, "mp_ct_default_secondary", "weapon_hkp2000"},
			{`sv_buy_status_override - -1`, "sv_buy_status_override", "-1"},
			{`sv_i_can_has_negative_float - -1.55`, "sv_i_can_has_negative_float", "-1.55"},
			{`sv_i_can_has_colon_separated_values - red;blue;green`, "sv_i_can_has_colon_separated_values", "red;blue;green"},
			{`mp_i_can_has_ip_address - 192.168.1.1`, "mp_i_can_has_ip_address", "192.168.1.1"},
			{`mp_can_has_empty_string -`, "mp_can_has_empty_string", ""},
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

		for _, test := range validCases {
			t.Run(test.msg, func(t *testing.T) {
				if actual, ok := parsEchoCvar(test.msg); !ok {
					t.Errorf("Message %q should have parsed successfully.", test.msg)
				} else {
					if actual.Name != test.expectedName {
						t.Errorf("Expected var name %q but got %q.", test.expectedName, actual.Name)
					}

					if actual.Value != test.expectedValue {
						t.Errorf("Expected var value %q but got %q.", test.expectedValue, actual.Value)
					}
				}
			})
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := map[string][]string{
			"CSGO": []string{
				`Loading map "de_nuke"`,
				`server cvars start`,
				`"Console<0><Console><Console>" say ""running gamemode_competitive_server.cfg""`,
				`World triggered "Round_Start"`,
				`"BlondeQuack<5><STEAM_1:1:07523420><CT>" money change 10900-200 = $10700 (tracked) (purchase: weapon_flashbang)`,
				`Molotov projectile spawned at -111.981644 -1925.359863 -347.735901, velocity 668.248840 -685.101196 179.108994`,
				`World triggered "Match_Start" on "de_nuke"`,
			},
			"OTHER": []string{
				``,
				`My first clue came at 4:15, when the clock stopped.`,
			},
		}

		for name, tests := range invalidCases {
			t.Run(name, func(t *testing.T) {
				for _, test := range tests {
					if _, ok := parsEchoCvar(test); ok {
						t.Errorf("Message %q should NOT have successfully parsed.", test)
					}
				}
			})
		}
	})
}

func Test_paresEchoServerCvar(t *testing.T) {
	t.Run("Valid Cases", func(t *testing.T) {
		validCases := []struct {
			msg              string
			expectedVarName  string
			expectedVarValue string
		}{
			{`server_cvar: "mp_whatever2" "-300"`, "mp_whatever2", "-300"},
			{`server_cvar: "cash_player_interact_with_hostage" "150"`, "cash_player_interact_with_hostage", "150"},
			{`server_cvar: "cash_team_rescued_hostage" "0"`, "cash_team_rescued_hostage", "0"},
			{`server_cvar: "mp_roundtime_hostage" "1.92"`, "mp_roundtime_hostage", "1.92"},
			{`server_cvar: "sv_negative_float" "-3.14"`, "sv_negative_float", "-3.14"},
			{`server_cvar: "sv_negative_integer" "-7"`, "sv_negative_integer", "-7"},
			{`server_cvar: "mp_empty_string" ""`, "mp_empty_string", ""},
			{`server_cvar: "cash_team_elimination_hostage_map_ct" "3000"`, "cash_team_elimination_hostage_map_ct", "3000"},
			{`server_cvar: "sm_nextmap" "de_dust2"`, "sm_nextmap", "de_dust2"},
			{`server_cvar: "tv_colon_separated_values" "red;green;blue"`, "tv_colon_separated_values", "red;green;blue"},
			{`server_cvar: "tf_server_identity_disable_quickplay" "1"`, "tf_server_identity_disable_quickplay", "1"},
			{`server_cvar: "sv_tags" "alltalk,cp,noquickplay"`, "sv_tags", "alltalk,cp,noquickplay"},
		}

		for _, test := range validCases {
			t.Run(test.msg, func(t *testing.T) {
				if actual, ok := paresEchoServerCvar(test.msg); !ok {
					t.Errorf("Message %q should have been parsed as a echo of a server cvar", test.msg)
				} else {
					if actual.Name != test.expectedVarName {
						t.Errorf("Expected var name %q but got %q", test.expectedVarName, actual.Name)
					}

					if actual.Value != test.expectedVarValue {
						t.Errorf("Expected var name %q but got %q", test.expectedVarValue, actual.Value)
					}
				}
			})
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := map[string][]string{
			"CSGO": []string{
				`Loading map "de_nuke"`,
				`server cvars start`,
				`"Console<0><Console><Console>" say ""running gamemode_competitive_server.cfg""`,
				`World triggered "Round_Start"`,
				`"BlondeQuack<5><STEAM_1:1:07523420><CT>" money change 10900-200 = $10700 (tracked) (purchase: weapon_flashbang)`,
				`Molotov projectile spawned at -111.981644 -1925.359863 -347.735901, velocity 668.248840 -685.101196 179.108994`,
				`World triggered "Match_Start" on "de_nuke"`,
			},
			"OTHER": []string{
				``,
				`My first clue came at 4:15, when the clock stopped.`,
			},
		}

		for name, tests := range invalidCases {
			t.Run(name, func(t *testing.T) {
				for _, test := range tests {
					if _, ok := paresEchoServerCvar(test); ok {
						t.Errorf("Message %q should NOT have successfully parsed.", test)
					}
				}
			})
		}
	})
}
