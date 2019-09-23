package csgo

import (
	"time"

	"github.com/rs/zerolog/log"
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

type matchPhase uint16

const (
	unknown matchPhase = 1 << iota
	freezePeriod
)

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
func (m *matchInfo) reset(start time.Time) {
	m.ended = time.Time{}
	m.started = start
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

func (g *gameInfo) teamName(t team) string {
	switch t {
	case mpTeam1:
		return g.mpTeamname1
	case mpTeam2:
		return g.mpTeamname2
	default:
		return ""
	}
}

// roundsWonCurrentMatch returns the number of rounds won for the current match by the specified team
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

func (g *gameInfo) scoresCurrentMatch() (mpTeam1Wins, mpTeam2Wins lastInt) {
	if len(g.matches) == 0 {
		return 0, 0
	}

	matchIndex := len(g.matches) - 1
	if len(g.matches[matchIndex].rounds) == 0 {
		return 0, 0
	}

	mpTeam1Wins, mpTeam2Wins = lastInt(0), lastInt(0)

	for _, round := range g.matches[matchIndex].rounds {
		if round.winningTeam == mpTeam1 {
			mpTeam1Wins++
		} else if round.winningTeam == mpTeam2 {
			mpTeam2Wins++
		}
	}

	return mpTeam1Wins, mpTeam2Wins
}

func (g *gameInfo) setRoundWinner(aff affiliation, t team, trigger string) {
	if len(g.matches) == 0 {
		g.matches = []matchInfo{matchInfo{}}
	}

	matchIndex := len(g.matches) - 1
	if len(g.matches[matchIndex].rounds) == 0 {
		g.matches[matchIndex].rounds = []roundInfo{}
	}

	g.matches[matchIndex].rounds = append(g.matches[matchIndex].rounds, roundInfo{
		winningAffiliation: aff,
		winningTeam:        t,
		winningTrigger:     trigger,
	})

	lastRound := int(g.currentMatchLastCompletedRound())
	mpTeam1Wins, mpTeam2Wins := g.scoresCurrentMatch()

	log.Info().Int("match", matchIndex+1).Int("round", lastRound).Int("team1_score", int(mpTeam1Wins)).Int("team2_score", int(mpTeam2Wins)).Msgf("Round %02d won by %v (%v as %v)", lastRound, t, g.teamName(t), aff)
}

// nextMatch will end the current match and start the next; if the current match has one or fewer completed round it will be reset and reused
// TODO - better  documentation!
func (g *gameInfo) nextMatch(mapName string, start time.Time) {
	if len(g.matches) == 0 {
		g.matches = append(g.matches, matchInfo{
			mapName: mapName,
			started: start,
		})
		log.Info().Msgf("Match 01 starting on map %q", mapName)
	}
	i := len(g.matches) - 1

	if g.currentMatchLastCompletedRound() >= 1 {
		// 1+ rounds have been completed; assume we completed the last match and are advancing ot the next
		g.matches[i].ended = start
		g.matches = append(g.matches, matchInfo{
			mapName: mapName,
			started: start,
		})
		i++
		log.Info().Msgf("Match %02d starting on map %q", i+1, mapName)
		return
	}

	if len(g.matches[i].rounds) > 0 && len(g.matches[i].rounds[0].winningTrigger) > 0 {
		log.Info().Msgf("Match %02d is restarting on map %q", i+1, mapName)
	}

	// Use to reset stats; even if no round history details are being reset
	g.matches[i].reset(start)
}
