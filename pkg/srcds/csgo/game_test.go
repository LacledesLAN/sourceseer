package csgo

import (
	"testing"
	"time"
)

func Test_lastCompletedRound(t *testing.T) {

}

func Test_setRoundWinner(t *testing.T) {
	m := &gameInfo{}

	m.setRoundWinner(counterterrorist, mpTeam1, "SFUI_Notice_CTs_Win")
	m.setRoundWinner(terrorist, mpTeam2, "SFUI_Notice_Terrorists_Win")
	m.setRoundWinner(counterterrorist, mpTeam1, "SFUI_Notice_CTs_Win")
	m.setRoundWinner(counterterrorist, mpTeam1, "SFUI_Notice_CTs_Win")

	if m.currentMatchLastCompletedRound() != 4 {
		t.Errorf("Expected last completed round to be %d not %d.", 4, m.currentMatchLastCompletedRound())
	}

	if m.roundsWonCurrentMatch(mpTeam1) != 3 {
		t.Errorf("Expected mpTeam1 wins to be %d not %d.", 3, m.roundsWonCurrentMatch(mpTeam1))
		t.Errorf("%v", m)
	}

	if m.roundsWonCurrentMatch(mpTeam2) != 1 {
		t.Errorf("Expected mpTeam2 wins to be %d not %d.", 1, m.roundsWonCurrentMatch(mpTeam2))
	}

	m.setRoundWinner(terrorist, mpTeam2, "SFUI_Notice_Terrorists_Win")
	m.setRoundWinner(terrorist, mpTeam2, "SFUI_Notice_Terrorists_Win")
	m.setRoundWinner(terrorist, mpTeam2, "SFUI_Notice_Terrorists_Win")

	if m.currentMatchLastCompletedRound() != 7 {
		t.Errorf("Expected last completed round to be %d not %d.", 7, m.currentMatchLastCompletedRound())
	}

	if m.roundsWonCurrentMatch(mpTeam2) != 4 {
		t.Errorf("Expected mpTeam2 wins to be %d not %d.", 4, m.roundsWonCurrentMatch(mpTeam2))
	}
}

func Test_roundsWonCurrentMatch(t *testing.T) {

}

func Test_matchInfo_resetMatch(t *testing.T) {
	expectedMapName := "...ship's ready except for this cup holder, and I should have that done in 12 hours."

	mock := &matchInfo{
		ended:   time.Now(),
		mapName: expectedMapName,
		started: time.Now(),
		rounds: []roundInfo{
			{winningAffiliation: counterterrorist, winningTeam: mpTeam1, winningTrigger: "Comets, the icebergs of the sky"},
			{winningAffiliation: terrorist, winningTeam: mpTeam2, winningTrigger: "Problem solved. You two fight to the death and I’ll cook the loser."},
		},
	}

	mock.reset()

	if !mock.ended.IsZero() {
		t.Errorf("Ended time did not get reset; was %q.", mock.ended)
	}

	if mock.mapName != expectedMapName {
		t.Errorf("Map name wasn't preserved; was %q", mock.mapName)
	}

	if len(mock.rounds) != 0 {
		t.Errorf("rounds[] should have been reset buy was %v", mock.rounds)
	}

	if mock.started.IsZero() {
		t.Error("Started time got reset.")
	}
}

func Test_gameInfo_restartMatch(t *testing.T) {
	//TODO!
}

func Test_gameInfo_nextMatch(t *testing.T) {
	//TODO!
	//create new entry doesn't crash
	//doesn't advance when zero rounds
}
