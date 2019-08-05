package csgo

import (
	"time"
)

// affiliation represents a player's assigned side (CT/T)
type affiliation string

const (
	unassigned       affiliation = "UNASSIGNED"
	counterterrorist affiliation = "CT"
	terrorist        affiliation = "TERRORIST"
)

// lastInt is an integer where the value is scoped to the previous round
type lastInt int

// matchInfo contains statistics about a match
type matchInfo struct {
	ended   time.Time
	mapName string
	rounds  []roundInfo
	started time.Time
}

// roundInfo contains statistics about a round
type roundInfo struct {
	winningAffiliation affiliation
	winningTeam        team
	winningTrigger     string
}

// team represents a players' team (mp_team1/mp_team2)
type team string

const (
	mpTeam1   team = "mp_team1"
	mpTeam2   team = "mp_team2"
	spectator team = "spectator"
)

// reset the match information; preserving the map name
func (m *matchInfo) reset() {
	m.ended = time.Time{}
	m.started = time.Now()
	m.rounds = []roundInfo{}
}

type gameInfo struct {
	matches     []matchInfo
	mpTeamname1 string
	mpTeamname2 string
}

func (g *gameInfo) currentMatchLastCompletedRound() lastInt {
	if len(g.matches) == 0 {
		return 0
	}
	matchIndex := len(g.matches) - 1

	return lastInt(len(g.matches[matchIndex].rounds))
}

func (g *gameInfo) roundsWonCurrentMatch(t team) int {
	if len(g.matches) == 0 {
		return 0
	}

	matchIndex := len(g.matches) - 1
	if len(g.matches[matchIndex].rounds) == 0 {
		return 0
	}

	r := 0
	for _, round := range g.matches[matchIndex].rounds {
		if t == round.winningTeam {
			r++
		}
	}

	return r
}

func (g *gameInfo) setRoundWinner(a affiliation, t team, trigger string) {
	if len(g.matches) == 0 {
		g.matches = []matchInfo{matchInfo{}}
	}

	matchIndex := len(g.matches) - 1
	if len(g.matches[matchIndex].rounds) == 0 {
		g.matches[matchIndex].rounds = []roundInfo{}
	}

	g.matches[matchIndex].rounds = append(g.matches[matchIndex].rounds, roundInfo{
		winningAffiliation: a,
		winningTeam:        t,
		winningTrigger:     trigger,
	})
}

// restart the current match
func (g *gameInfo) restartMatch() {
	if len(g.matches) == 0 {
		return
	}

	g.matches[len(g.matches)-1].reset()
}

// nextMatch will end the current match and start the next; if the current match has one or fewer completed round it will be reset and reused
func (g *gameInfo) nextMatch(mapName string) {
	if len(g.matches) == 0 {
		g.matches = append(g.matches, matchInfo{
			mapName: mapName,
		})
	}
	i := len(g.matches) - 1

	// 2+ rounds have been completed; assume we can advance to next match
	if g.currentMatchLastCompletedRound() > 1 {
		g.matches[i].ended = time.Now()
		g.matches = append(g.matches, matchInfo{
			mapName: mapName,
		})
		i++
	}

	g.matches[i].reset()
}
