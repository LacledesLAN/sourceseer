package srcds

import (
	"testing"
)

func Test_Clients(t *testing.T) {
	c0 := Client{Username: "Lulubelle 7", SteamID: "7r355 m4cn31ll3", ServerSlot: "", Affiliation: ""}
	c1 := Client{Username: "Animatronio", SteamID: "d4v1d h3rm4n", ServerSlot: "3", Affiliation: ""}
	c2 := Client{Username: "Parts Hilton", SteamID: "7h3 7h13f 0f b46h34d", ServerSlot: "", Affiliation: "b46h34d"}

	var sut Clients

	sut.ClientJoined(c0)
	sut.ClientJoined(c1)
	sut.ClientJoined(c2)
	sut.ClientJoined(c2) // Make sure client can't join twice

	if len(sut) != 3 {
		t.Error("Should have 3 clients.")
	}

	if !sut.HasClient(c1) {
		t.Errorf("Client %q should have been found.", c1.Username)
	}

	sut.ClientDropped(c1)
	sut.ClientDropped(c1) // Make sure doesn't panic

	if len(sut) != 2 {
		t.Errorf("Should have 2 clients not %d.", len(sut))
	}

	if sut.HasClient(c1) {
		t.Errorf("Client %q should not have been found.", c1.Username)
	}
}

func Test_ParseClient(t *testing.T) {
	datum := []struct {
		actual   string
		expected Client
	}{
		{`"Lulubelle 7<6><7r355:m4cn31ll3><CT>"`, Client{Username: "Lulubelle 7", SteamID: "7r355:m4cn31ll3", ServerSlot: "6", Affiliation: "CT"}},
		{`"Lulubelle 7<48><BOT><CT>"`, Client{Username: "Lulubelle 7", SteamID: "BOT", ServerSlot: "48", Affiliation: "CT"}},
	}

	for _, testData := range datum {
		t.Run(testData.actual, func(t *testing.T) {
			c, err := ParseClient(testData.actual)

			if err != nil {
				t.Error("Reason: ", err)
			}

			if c.Username != testData.expected.Username {
				t.Errorf("Expected Username '%q' but got '%q' instead.", testData.expected.Username, c.Username)
			}

			if c.SteamID != testData.expected.SteamID {
				t.Errorf("Expected SteamID '%q' but got '%q' instead.", testData.expected.SteamID, c.SteamID)
			}

			if c.ServerSlot != testData.expected.ServerSlot {
				t.Errorf("Expected ServerSlot '%q' but got '%q' instead.", testData.expected.ServerSlot, c.ServerSlot)
			}

			if c.Affiliation != testData.expected.Affiliation {
				t.Errorf("Expected ServerTeam '%q' but got '%q' instead.", testData.expected.Affiliation, c.Affiliation)
			}
		})
	}
}
