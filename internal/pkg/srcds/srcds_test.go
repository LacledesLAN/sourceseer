package srcds

import (
	"testing"
	"time"
)

func ClientsAreIdentical(p1, p2 *Client) bool {
	if p1 == nil || p2 == nil {
		return p1 == nil && p2 == nil
	}

	return p1.Username == p2.Username && p1.ServerSlot == p1.ServerSlot && p1.ServerTeam == p2.ServerTeam && p1.SteamID == p2.SteamID
}

func Test_ExtractClients(t *testing.T) {
	testDatum := []struct {
		rawLogEntry string
		originator  *Client
		target      *Client
	}{
		// Garbage
		{"", nil, nil},
		{"Bender, quit destroying the universe!", nil, nil},
		// Console
		{`"Console<0><Console><Console>" say "I’m so embarrassed. I wish everybody else was dead."`, &Client{Username: "Console", ServerSlot: "0", ServerTeam: "Console", SteamID: "Console"}, nil},
		// GOTV
		{`"GOTV<2><BOT><>" connected, address ""`, &Client{Username: "GOTV", ServerSlot: "2", ServerTeam: "", SteamID: "BOT"}, nil},
		{`"GOTV<2><BOT><>" entered the game`, &Client{Username: "GOTV", ServerSlot: "2", ServerTeam: "", SteamID: "BOT"}, nil},
		{`"GOTV<2><BOT><Unassigned>" changed name to "Roberto V2.0"`, &Client{Username: "GOTV", ServerSlot: "2", ServerTeam: "Unassigned", SteamID: "BOT"}, nil},
		// Player Connect
		{`"A<7><STEAM_1:0:1234567><>" connected, address ""`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "", SteamID: "STEAM_1:0:1234567"}, nil},
		{`"A<7><STEAM_1:0:1234567><>" STEAM USERID validated`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "", SteamID: "STEAM_1:0:1234567"}, nil},
		{`"A<7><STEAM_1:0:1234567><>" entered the game`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "", SteamID: "STEAM_1:0:1234567"}, nil},
		{`"A<7><STEAM_1:0:1234567>" switched from team <Unassigned> to <CT>`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "", SteamID: "STEAM_1:0:1234567"}, nil},
		// One Player
		{`"A<7><STEAM_1:0:1234567><CT>" left buyzone with [ ]`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "CT", SteamID: "STEAM_1:0:1234567"}, nil},
		{`"A<7><STEAM_1:0:1234567><CT>" purchased "taser"`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "CT", SteamID: "STEAM_1:0:1234567"}, nil},
		{`"A<7><STEAM_1:0:1234567><CT>" left buyzone with [ weapon_knife_t weapon_glock ]`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "CT", SteamID: "STEAM_1:0:1234567"}, nil},
		{`"A<7><STEAM_1:0:1234567><CT>" threw smokegrenade [1317 601 125]`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "CT", SteamID: "STEAM_1:0:1234567"}, nil},
		{`"A<7><STEAM_1:0:1234567><CT>" [2504 -344 -289] committed suicide with "world"`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "CT", SteamID: "STEAM_1:0:1234567"}, nil},
		{`"A<7><STEAM_1:0:1234567><CT>" triggered "Got_The_Bomb"`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "CT", SteamID: "STEAM_1:0:1234567"}, nil},
		{`"A<7><STEAM_1:0:1234567><CT>" say "Blackmail is such an ugly word. I prefer extortion. The ‘x’ makes it sound cool."`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "CT", SteamID: "STEAM_1:0:1234567"}, nil},
		{`"A<7><STEAM_1:0:1234567><CT>" triggered "Dropped_The_Bomb"`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "CT", SteamID: "STEAM_1:0:1234567"}, nil},
		{`"A<7><STEAM_1:0:1234567><CT>" triggered "Planted_The_Bomb"`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "CT", SteamID: "STEAM_1:0:1234567"}, nil},
		// Two players
		//{`"A<7><STEAM_1:0:1234567><CT>" [756 -1951 -416] attacked "B<4><STEAM_1:1:9876543><TERRORIST>" [824 -1933 -416] with "glock" (damage "117") (damage_armor "0") (health "0") (armor "0") (hitgroup "head")`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "CT", SteamID: "STEAM_1:0:1234567"}, &Client{Username: "B", ServerSlot: "4", ServerTeam: "TERRORIST", SteamID: "STEAM_1:1:9876543"}},
		//{`"A<7><STEAM_1:0:1234567><CT>" [756 -1951 -416] killed "B<4><STEAM_1:1:9876543><TERRORIST>" [824 -1933 -352] with "glock" (headshot)`, &Client{Username: "A", ServerSlot: "7", ServerTeam: "CT", SteamID: "STEAM_1:0:1234567"}, &Client{Username: "B", ServerSlot: "4", ServerTeam: "TERRORIST", SteamID: "STEAM_1:1:9876543"}},
	}

	for _, testData := range testDatum {
		logEntry := LogEntry{Message: testData.rawLogEntry, Raw: testData.rawLogEntry, Timestamp: time.Now()}

		t.Run(testData.rawLogEntry, func(t *testing.T) {
			originator, target := ExtractClients(logEntry)

			if ClientsAreIdentical(originator, testData.originator) != true {
				t.Errorf("Expected originator client %#v but got %#v", testData.originator, originator)
			}

			if ClientsAreIdentical(target, testData.target) != true {
				t.Errorf("Expected target client %#v but got %#v", testData.target, target)
			}
		})
	}
}
