package csgo

import (
	"testing"

	"github.com/lacledeslan/sourceseer/srcds"
)

var (
	mapStateTestPlayer0 = srcds.Client{Username: "Billionairebot", SteamID: "ph1l-l4m4rr"}
	mapStateTestPlayer1 = srcds.Client{Username: "Suspendington"}
	mapStateTestPlayer2 = srcds.Client{Username: "Titanius Anglesmith", SteamID: "j0hn-d1m46610"}
)

func Test_CTPlayerJoin(testing *testing.T) {
	sut := mapState{}
	ct := sut.getCT()
	t := sut.getT()

	// Joins when empty
	sut.CTPlayerJoin(mapStateTestPlayer0)
	if ct.playerIndex(mapStateTestPlayer0) == -1 {
		testing.Errorf("CT should have player %q.", mapStateTestPlayer0)
	}

	// Joins when not empty
	sut.CTPlayerJoin(mapStateTestPlayer1)
	if ct.playerIndex(mapStateTestPlayer1) == -1 {
		testing.Errorf("CT should have player %q.", mapStateTestPlayer1)
	}

	// Doesn't re-join when already on team
	sut.CTPlayerJoin(mapStateTestPlayer1)
	if len(ct.knownPlayers) != 2 {
		testing.Errorf("Player %q should not be able to join twice.", mapStateTestPlayer1)
	}

	// Gets removed from T
	sut.TerroristPlayerJoin(mapStateTestPlayer2)
	sut.CTPlayerJoin(mapStateTestPlayer2)

	if t.playerIndex(mapStateTestPlayer2) != -1 {
		testing.Errorf("Player %q should have been removed from T when joining CT.", mapStateTestPlayer2)
	}
}

func Test_CTWinRound(testing *testing.T) {
	sut := mapState{}
	ct := sut.getCT()

	if ct.roundsWon != 0 {
		testing.Errorf("New CT Team should show 0 rounds won but have %d.", ct.roundsWon)
	}

	if sut.roundNumber != 0 {
		testing.Errorf("New mapstate should show round 0 but was %d.", sut.roundNumber)
	}

	for i := 0; i < 4; i++ {
		sut.CTWinRound()
	}

	if ct.roundsWon != 4 {
		testing.Errorf("CT Team should show 4 rounds won but have %d.", ct.roundsWon)
	}

	if sut.roundNumber != 4 {
		testing.Errorf("Mapstate should show round 4 but was %d.", sut.roundNumber)
	}
}

func Test_SwapSides(testing *testing.T) {
	sut := mapState{}
	originalCT := sut.getCT()
	originalT := sut.getT()

	sut.SwapSides()

	newCT := sut.getCT()
	newT := sut.getT()

	if &newCT == &originalCT {
		testing.Errorf("The memory address for `newCT` should not match the address for `originalCT` (%X)", &newCT)
	}

	if &newT == &originalT {
		testing.Errorf("The memory address for `newT` should not match the address for `originalT` (%X)", &newT)
	}
}

func Test_TerroristPlayerJoin(testing *testing.T) {
	sut := mapState{}
	t := sut.getT()
	ct := sut.getCT()

	// Joins when empty
	sut.TerroristPlayerJoin(mapStateTestPlayer0)
	if t.playerIndex(mapStateTestPlayer0) == -1 {
		testing.Errorf("T should have player %q.", mapStateTestPlayer0)
	}

	// Joins when not empty
	sut.TerroristPlayerJoin(mapStateTestPlayer1)
	if t.playerIndex(mapStateTestPlayer1) == -1 {
		testing.Errorf("T should have player %q.", mapStateTestPlayer1)
	}

	// Doesn't re-join when already on team
	sut.TerroristPlayerJoin(mapStateTestPlayer1)
	if len(t.knownPlayers) != 2 {
		testing.Errorf("Player %q should not be able to join twice.", mapStateTestPlayer1)
	}

	// Gets removed from CT
	sut.CTPlayerJoin(mapStateTestPlayer2)
	sut.TerroristPlayerJoin(mapStateTestPlayer2)

	if ct.playerIndex(mapStateTestPlayer2) != -1 {
		testing.Errorf("Player %q should have been removed from CT when joining T.", mapStateTestPlayer2)
	}
}

func Test_TerroristWinRound(testing *testing.T) {
	sut := mapState{}
	t := sut.getT()

	if t.roundsWon != 0 {
		testing.Errorf("New T Team should show 0 rounds won but have %d.", t.roundsWon)
	}

	if sut.roundNumber != 0 {
		testing.Errorf("New mapstate should show round 0 but was %d.", sut.roundNumber)
	}

	for i := 0; i < 4; i++ {
		sut.TerroristWinRound()
	}

	if t.roundsWon != 4 {
		testing.Errorf("T Team should show 4 rounds won but have %d.", t.roundsWon)
	}

	if sut.roundNumber != 4 {
		testing.Errorf("Mapstate should show round 4 but was %d.", sut.roundNumber)
	}
}
