package csgo

import (
	"testing"

	"github.com/lacledeslan/sourceseer/srcds"
)

var (
	teamStateTestClient0 = Player{Client: srcds.Client{Username: "Lulubelle 7", SteamID: "7r355-m4cn31ll3"}}
	teamStateTestClient1 = Player{Client: srcds.Client{Username: "Daisy-Mae 128K"}}
	teamStateTestClient2 = Player{Client: srcds.Client{Username: "The Crushinator", SteamID: "m4ur1c3-l4m4rch3"}}
)

func Test_HasPlayer(t *testing.T) {
	sut := teamState{}

	sut.PlayerJoined(teamStateTestClient0)
	sut.PlayerJoined(teamStateTestClient2)

	if !sut.HasPlayer(teamStateTestClient2) {
		t.Errorf("Team should have player %q", teamStateTestClient2.Username)
	}

	if sut.HasPlayer(teamStateTestClient1) {
		t.Errorf("Team should not have player %q", teamStateTestClient1.Username)
	}
}

func Test_PlayerCount(t *testing.T) {
	sut := teamState{}

	if sut.PlayerCount() != uint8(0) {
		t.Error("Client count should be 0.")
	}

	sut.PlayerJoined(teamStateTestClient0)
	sut.PlayerJoined(teamStateTestClient1)
	sut.PlayerJoined(teamStateTestClient2)

	if sut.PlayerCount() != uint8(3) {
		t.Error("Client count should be 3.")
	}
}

func TestMapState_PlayerDropped(t *testing.T) {
	sut := teamState{}

	sut.PlayerJoined(teamStateTestClient0)
	sut.PlayerJoined(teamStateTestClient1)
	sut.PlayerJoined(teamStateTestClient2)

	sut.PlayerDropped(teamStateTestClient0)
	sut.PlayerDropped(teamStateTestClient2)

	if sut.HasPlayer(teamStateTestClient0) {
		t.Errorf("Team should not have player %q", teamStateTestClient0.Username)
	}

	if !sut.HasPlayer(teamStateTestClient1) {
		t.Errorf("Team should have player %q", teamStateTestClient1.Username)
	}

	if sut.HasPlayer(teamStateTestClient2) {
		t.Errorf("Team should not have player %q", teamStateTestClient2.Username)
	}

	sut.PlayerDropped(teamStateTestClient1)
	if sut.HasPlayer(teamStateTestClient1) {
		t.Errorf("Team should not have player %q", teamStateTestClient1.Username)
	}
}

func Test_PlayerJoined(t *testing.T) {
	sut := teamState{}

	sut.PlayerJoined(teamStateTestClient0)
	sut.PlayerJoined(teamStateTestClient2)

	if !sut.HasPlayer(teamStateTestClient2) {
		t.Errorf("Team should have player %q", teamStateTestClient2.Username)
	}
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
