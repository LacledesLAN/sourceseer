package csgo

import (
	"testing"

	"github.com/lacledeslan/sourceseer/srcds"
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
