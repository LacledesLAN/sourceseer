package srcds

import "testing"

func Test_Clients(t *testing.T) {
	c0 := Client{Username: "Lulubelle 7", SteamID: "7r355 m4cn31ll3", ServerSlot: "", ServerTeam: ""}
	c1 := Client{Username: "Animatronio", SteamID: "d4v1d h3rm4n", ServerSlot: "3", ServerTeam: ""}
	c2 := Client{Username: "Parts Hilton", SteamID: "7h3 7h13f 0f b46h34d", ServerSlot: "", ServerTeam: "b46h34d"}

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
