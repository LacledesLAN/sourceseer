package csgo

import (
	"testing"

	"github.com/lacledeslan/sourceseer/srcds"
)

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

			if result.IsReady != false {
				t.Error("Default value for 'IsReady' was not properly set.")
			}
		})
	}
}

func Test_ToSrcdsClient(t *testing.T) {
	// TODO - add some tests
}
