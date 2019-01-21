package csgo

import (
	"testing"

	"github.com/lacledeslan/sourceseer/srcds"
)

var (
	teamStateTestPlayer0 = srcds.Client{Username: "Lulubelle 7", SteamID: "7r355-m4cn31ll3"}
	teamStateTestPlayer1 = srcds.Client{Username: "Daisy-Mae 128K"}
	teamStateTestPlayer2 = srcds.Client{Username: "The Crushinator", SteamID: "m4ur1c3-l4m4rch3"}
)

func Test_PlayerJoin(t *testing.T) {
	sut := teamState{}

	sut.PlayerJoin(teamStateTestPlayer0)
	if len(sut.knownPlayers) != 1 {
		t.Error("Should have 1 player")
	}

	sut.PlayerJoin(teamStateTestPlayer1)
	if len(sut.knownPlayers) != 2 {
		t.Error("Should have 2 players")
	}

	sut.PlayerJoin(teamStateTestPlayer0)
	if len(sut.knownPlayers) != 2 {
		t.Error("Should still have 2 players")
	}

	sut.PlayerJoin(teamStateTestPlayer2)
	if len(sut.knownPlayers) != 3 {
		t.Error("Should still have 3 players")
	}

	sut.PlayerJoin(teamStateTestPlayer0)
	if len(sut.knownPlayers) != 3 {
		t.Error("Should still have 3 player")
	}
}

func Test_PlayerRemove(t *testing.T) {
	sut := teamState{}

	sut.PlayerRemove(teamStateTestPlayer2) // Make sure there's no panic

	sut.PlayerJoin(teamStateTestPlayer1)
	sut.PlayerRemove(teamStateTestPlayer1)

	if len(sut.knownPlayers) != 0 {
		t.Errorf("Should have 0 players not %d", len(sut.knownPlayers))
	}

	sut.PlayerJoin(teamStateTestPlayer0)
	sut.PlayerJoin(teamStateTestPlayer1)
	sut.PlayerJoin(teamStateTestPlayer2)
	sut.PlayerRemove(teamStateTestPlayer1)
	if len(sut.knownPlayers) != 2 {
		t.Errorf("Should have 2 players not %d", len(sut.knownPlayers))
	}

	sut.PlayerRemove(teamStateTestPlayer0)
	if sut.playerIndex(teamStateTestPlayer2) < 0 {
		t.Errorf("Player %q should be found in %q", teamStateTestPlayer2, sut.knownPlayers)
	}

	sut.PlayerRemove(teamStateTestPlayer1) // Make sure there's no panic
}

func Test_SetName(t *testing.T) {
	testDatum := []struct {
		actualName   string
		expectedName string
	}{
		{"", "Unspecified"},
		{"  ", "Unspecified"},
		{"\t\r\n", "Unspecified"},
		{"", "Unspecified"},
		{"   aaa\t\n\r", "aaa"},
		{"   a A  a\t\n\r", "a_A_a"},
	}

	sut := teamState{}

	for _, testData := range testDatum {
		t.Run(testData.actualName, func(t *testing.T) {
			sut.SetName(testData.actualName)
			result := sut.name

			if result != testData.expectedName {
				t.Errorf("Expected name %q but got %q", testData.expectedName, result)
			}
		})
	}
}
