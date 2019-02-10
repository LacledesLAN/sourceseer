package csgo

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
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

func Test_CTSetScore(testing *testing.T) {
	sut := mapState{}

	if sut.ct().roundsLost != 0 {
		testing.Error("CT should have 0 rounds lost.")
	}

	if sut.ct().roundsWon != 0 {
		testing.Error("CT should have 0 rounds won.")
	}

	if sut.terrorist().roundsLost != 0 {
		testing.Error("T should have 0 rounds lost.")
	}

	if sut.terrorist().roundsWon != 0 {
		testing.Error("T should have 0 rounds won.")
	}

	sut.CTSetScore("18")

	if sut.ct().roundsWon != 18 {
		testing.Error("CT should have 18 rounds won.")
	}

	if sut.terrorist().roundsLost != 18 {
		testing.Error("T should have 18 rounds lost.")
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

func Test_RoundsCompleted(testing *testing.T) {
	sut := mapState{}

	sut.CTSetScore("1")
	sut.TerroristSetScore("0")
	if sut.RoundsCompleted() != 1 {
		testing.Errorf("Rounds completed should be equal to 1 but was %d", sut.RoundsCompleted())
	}

	sut.CTSetScore("2")
	sut.TerroristSetScore("0")
	if sut.RoundsCompleted() != 2 {
		testing.Errorf("Rounds completed should be equal to 2 but was %d", sut.RoundsCompleted())
	}

	sut.CTSetScore("2")
	sut.TerroristSetScore("1")
	if sut.RoundsCompleted() != 3 {
		testing.Errorf("Rounds completed should be equal to 3 but was %d", sut.RoundsCompleted())
	}

	sut.TerroristSetScore("7")
	if sut.RoundsCompleted() != 9 {
		testing.Errorf("Rounds completed should be equal to 9 but was %d", sut.RoundsCompleted())
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

func Test_TerroristSetScore(testing *testing.T) {
	sut := mapState{}

	if sut.ct().roundsLost != 0 {
		testing.Error("CT should have 0 rounds lost.")
	}

	if sut.ct().roundsWon != 0 {
		testing.Error("CT should have 0 rounds won.")
	}

	if sut.terrorist().roundsLost != 0 {
		testing.Error("T should have 0 rounds lost.")
	}

	if sut.terrorist().roundsWon != 0 {
		testing.Error("T should have 0 rounds won.")
	}

	sut.TerroristSetScore("65")

	if sut.ct().roundsLost != 65 {
		testing.Error("CT should have 65 rounds lost.")
	}

	if sut.terrorist().roundsWon != 65 {
		testing.Error("T should have 65 rounds won.")
	}
}
