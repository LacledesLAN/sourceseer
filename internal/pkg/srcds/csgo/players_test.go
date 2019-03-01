package csgo

import (
	"testing"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

func Test_Clients(t *testing.T) {
	p0 := Player{Client: srcds.Client{Username: "Lulubelle 7", SteamID: "7r355 m4cn31ll3", ServerSlot: "", ServerTeam: ""}}
	p1 := Player{Client: srcds.Client{Username: "Animatronio", SteamID: "d4v1d h3rm4n", ServerSlot: "3", ServerTeam: ""}}
	p2 := Player{Client: srcds.Client{Username: "Parts Hilton", SteamID: "7h3 7h13f 0f b46h34d", ServerSlot: "", ServerTeam: "b46h34d"}}

	var sut Players

	sut.PlayerJoined(p0)
	sut.PlayerJoined(p1)
	sut.PlayerJoined(p1) // Make sure client can't join twice
	sut.PlayerJoined(p2)

	if len(sut) != 3 {
		t.Error("Should have 3 players.")
	}

	if !sut.HasPlayer(p1) {
		t.Errorf("Player %q should have been found.", p1.Username)
	}

	sut.PlayerDropped(p1)
	sut.PlayerDropped(p1) // Make sure doesn't panic

	if len(sut) != 2 {
		t.Errorf("Should have 2 players not %d.", len(sut))
	}

	if sut.HasPlayer(p1) {
		t.Errorf("PLayer %q should not have been found.", p1.Username)
	}

}

func Test_playerFromSrcdsClient(t *testing.T) {
	testDatum := []struct {
		input    srcds.Client
		expected Player
	}{
		{input: srcds.Client{Username: "Robo-Puppy", SteamID: "b1lly w357", ServerSlot: "12", ServerTeam: "blu"},
			expected: Player{Client: srcds.Client{Username: "Robo-Puppy", SteamID: "b1lly w357", ServerSlot: "12", ServerTeam: "blu"}}},
	}

	for _, testData := range testDatum {
		t.Run(testData.input.Username, func(t *testing.T) {
			result := playerFromSrcdsClient(testData.input)

			if result.Username != testData.input.Username {
				t.Errorf("Username %q did not carry over.", testData.input.Username)
			}

			if result.SteamID != testData.input.SteamID {
				t.Errorf("SteamID %q did not carry over.", testData.input.SteamID)
			}

			if result.ServerSlot != testData.input.ServerSlot {
				t.Errorf("ServerSlot %q did not carry over.", testData.input.ServerSlot)
			}

			if result.ServerTeam != testData.input.ServerTeam {
				t.Errorf("ServerTeam %q did not carry over.", testData.input.ServerTeam)
			}

			if len(result.flags) > 0 {
				t.Error("Flags slice should be empty.")
			}
		})
	}
}
