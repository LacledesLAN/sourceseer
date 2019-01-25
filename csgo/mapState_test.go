package csgo

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/lacledeslan/sourceseer/srcds"
)

var (
	mapStateTestPlayer0 = Player{Client: srcds.Client{Username: "Billionairebot", SteamID: "ph1l-l4m4rr"}}
	mapStateTestPlayer1 = Player{Client: srcds.Client{Username: "Suspendington"}}
	mapStateTestPlayer2 = Player{Client: srcds.Client{Username: "Titanius Anglesmith", SteamID: "j0hn-d1m46610"}}
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	returnCode := m.Run()
	os.Exit(returnCode)
}

func Test_CTWonRound(testing *testing.T) {
	sut := mapState{}
	ct := sut.ct()
	t := sut.terrorist()

	if ct.roundsWon != 0 {
		testing.Errorf("New CT Team should show 0 rounds won but have %d.", ct.roundsWon)
	}

	if sut.roundNumber != 0 {
		testing.Errorf("New mapstate should show round 0 but was %d.", sut.roundNumber)
	}

	roundsToTest := uint8(rand.Intn(14) + 2)

	for i := 0; i < int(roundsToTest); i++ {
		sut.CTWonRound()
	}

	if ct.roundsWon != roundsToTest {
		testing.Errorf("CT Team should show %d rounds won but have %d.", roundsToTest, ct.roundsWon)
	}

	if sut.roundNumber != roundsToTest {
		testing.Errorf("Mapstate should show round %d but was %d.", roundsToTest, sut.roundNumber)
	}

	if t.roundsLost != roundsToTest {
		testing.Errorf("Terrorist Team should show %d rounds lost but have %d.", roundsToTest, ct.roundsLost)
	}
}

func Test_PlayerDropped(testing *testing.T) {
	sut := mapState{}

	sut.PlayerJoinedCT(mapStateTestPlayer0)
	sut.PlayerDropped(mapStateTestPlayer0)

	ct := sut.ct()
	if ct.HasPlayer(mapStateTestPlayer0) {
		testing.Errorf("CT should be empty but was %v.", ct)
	}

	sut.PlayerJoinedTerrorist(mapStateTestPlayer0)
	sut.PlayerDropped(mapStateTestPlayer0)

	t := sut.terrorist()
	if t.HasPlayer(mapStateTestPlayer0) {
		testing.Errorf("T should be empty but was %v.", t)
	}
}

func Test_PlayerJoinedCT(testing *testing.T) {
	sut := mapState{}
	ct := sut.ct()

	// Joins when empty
	sut.PlayerJoinedCT(mapStateTestPlayer0)
	if !ct.HasPlayer(mapStateTestPlayer0) {
		testing.Errorf("CT should have player %q.", mapStateTestPlayer0.Username)
	}

	// Joins when not empty
	sut.PlayerJoinedCT(mapStateTestPlayer1)
	if !ct.HasPlayer(mapStateTestPlayer1) {
		testing.Errorf("CT should have player %q.", mapStateTestPlayer0.Username)
	}

	// Doesn't re-join when already on team
	sut.PlayerJoinedCT(mapStateTestPlayer1)
	if ct.PlayerCount() != 2 {
		testing.Errorf("CT should have 2 players but had %d.", ct.PlayerCount())
	}

	// Gets removed from T
	sut.PlayerJoinedTerrorist(mapStateTestPlayer2)
	sut.PlayerJoinedCT(mapStateTestPlayer2)
	if !ct.HasPlayer(mapStateTestPlayer2) {
		testing.Errorf("CT should have player %q.", mapStateTestPlayer2.Username)
	}

	t := sut.terrorist()
	if t.HasPlayer(mapStateTestPlayer2) {
		testing.Errorf("T should not have player %q.", mapStateTestPlayer2.Username)
	}
}

func Test_PlayerJoinedTerrorist(testing *testing.T) {
	sut := mapState{}
	t := sut.terrorist()

	// Joins when empty
	sut.PlayerJoinedTerrorist(mapStateTestPlayer0)
	if !t.HasPlayer(mapStateTestPlayer0) {
		testing.Errorf("T should have player %q.", mapStateTestPlayer0.Username)
	}

	// Joins when not empty
	sut.PlayerJoinedTerrorist(mapStateTestPlayer1)
	if !t.HasPlayer(mapStateTestPlayer1) {
		testing.Errorf("T should have player %q.", mapStateTestPlayer0.Username)
	}

	// Doesn't re-join when already on team
	sut.PlayerJoinedTerrorist(mapStateTestPlayer1)
	if t.PlayerCount() != 2 {
		testing.Errorf("T should have 2 players but had %d.", t.PlayerCount())
	}

	// Gets removed from CT
	sut.PlayerJoinedCT(mapStateTestPlayer2)
	sut.PlayerJoinedTerrorist(mapStateTestPlayer2)
	if !t.HasPlayer(mapStateTestPlayer2) {
		testing.Errorf("T should have player %q.", mapStateTestPlayer2.Username)
	}

	ct := sut.ct()
	if ct.HasPlayer(mapStateTestPlayer2) {
		testing.Errorf("CT should not have player %q.", mapStateTestPlayer2.Username)
	}
}

func Test_TeamsSwappedSides(testing *testing.T) {
	sut := mapState{}
	originalCT := sut.ct()
	originalT := sut.terrorist()

	sut.TeamsSwappedSides()

	newCT := sut.ct()
	newT := sut.terrorist()

	if &newCT == &originalCT {
		testing.Errorf("The memory address for `newCT` should not match the address for `originalCT` (%X)", &newCT)
	}

	if &newT == &originalT {
		testing.Errorf("The memory address for `newT` should not match the address for `originalT` (%X)", &newT)
	}
}

func Test_TerroristWonRound(testing *testing.T) {
	sut := mapState{}
	ct := sut.ct()
	t := sut.terrorist()

	if t.roundsWon != 0 {
		testing.Errorf("New T Team should show 0 rounds won but have %d.", t.roundsWon)
	}

	if sut.roundNumber != 0 {
		testing.Errorf("New mapstate should show round 0 but was %d.", sut.roundNumber)
	}

	roundsToTest := uint8(rand.Intn(30) + 3)

	for i := 0; i < int(roundsToTest); i++ {
		sut.TerroristWonRound()
	}

	if t.roundsWon != roundsToTest {
		testing.Errorf("T Team should show %d rounds won but have %d.", roundsToTest, t.roundsWon)
	}

	if sut.roundNumber != roundsToTest {
		testing.Errorf("Mapstate should show round %d but was %d.", roundsToTest, sut.roundNumber)
	}

	if ct.roundsLost != roundsToTest {
		testing.Errorf("CT Team should show %d rounds lost but have %d.", roundsToTest, sut.roundNumber)
	}
}
