package csgo

import (
	"io"
	"strings"

	"github.com/lacledeslan/sourceseer/pkg/srcds"
	"github.com/rs/zerolog/log"
)

//Observer for watching CSGO log streams
type Observer interface {
	Start()
}

// NewReader for observing streaming CSGO data
func NewReader(r io.Reader, mpHalftime, mpMaxRounds, mpMaxOvertimeRounds int) Observer {
	o := &observer{
		srcdsObserver: srcds.NewReader(r),
	}

	o.srcdsObserver.AddCvarWatcherDefault("mp_halftime", string(mpHalftime))
	o.srcdsObserver.AddCvarWatcherDefault("mp_maxrounds", string(mpMaxRounds))
	o.srcdsObserver.AddCvarWatcherDefault("mp_overtime_maxrounds", string(mpMaxOvertimeRounds))

	return o
}

// Starts the observer
func (o *observer) Start() {
	for le := range o.srcdsObserver.Start() {
		o.processLogEntry(le)
	}
}

type observer struct {
	players struct {
		mpTeam1    srcds.Clients
		mpTeam2    srcds.Clients
		unassigned srcds.Clients
	}
	game          gameInfo
	srcdsObserver srcds.Observer
}

func (o *observer) processLogEntry(le srcds.LogEntry) {
	if clientLog, ok := srcds.ParseClientLogEntry(le); ok {
		if _, ok := parseClientSay(clientLog); ok {
			// TODO: process the client saying something
			return
		}

		if ok := srcds.ParseClientConnected(clientLog); ok {
			o.playerJoined(unassigned, clientLog.Client)
			return
		}

		if m, ok := parseClientSetAffiliation(clientLog); ok {
			o.playerJoined(m.to, clientLog.Client)
			return
		}

		if _, ok := srcds.ParseClientDisconnected(clientLog); ok {
			o.playerDropped(clientLog.Client)
		}

		return
	}

	if parseStartingFreezePeriod(le) {
		log.Info().Msg("Starting Freeze Period")
		return
	}

	if worldLog, ok := parseWorldTrigger(le); ok {
		if mapName, ok := parseWorldTriggerMatchStart(worldLog); ok {
			o.game.nextMatch(mapName, le.Timestamp)
		}

		if parseWorldTriggerRoundStart(worldLog) {
			log.Info().Msg("Round Start")
		}

		if parseWorldTriggerGameCommencing(worldLog) {
			log.Info().Msg("Game Commencing")
		}

		return
	}

	if strings.HasPrefix(le.Message, "Team") {
		if msg, ok := parseTeamTriggered(le); ok {
			team := o.getTeam(msg.affiliation)
			o.game.setRoundWinner(msg.affiliation, team, msg.trigger)

			// Let's see if a team won
			maxrounds, _ := o.srcdsObserver.TryCvarAsInt("mp_maxrounds", defaultMpMaxrounds)
			otMaxrounds, _ := o.srcdsObserver.TryCvarAsInt("mp_overtime_maxrounds", defaultMpOvertimeMaxrounds)

			if winThreshold := calculateLastRoundWinThreshold(maxrounds, otMaxrounds, o.game.currentMatchLastCompletedRound()); o.game.currentMatchLastCompletedRound() >= winThreshold {
				matchNum := len(o.game.matches)
				roundNum := int(o.game.currentMatchLastCompletedRound())
				mpTeam1Wins, mpTeam2Wins := o.game.scoresCurrentMatch()
				winningTeam := spectator

				if mpTeam1Wins >= winThreshold {
					winningTeam = mpTeam1
				} else if mpTeam2Wins >= winThreshold {
					winningTeam = mpTeam2
				} else {
					return
				}

				log.Info().Int("match", matchNum).Int("round", roundNum).Int("team1_score", int(mpTeam1Wins)).Int("team2_score", int(mpTeam2Wins)).Msgf("Match %02d clinched by %v (%v)", matchNum, winningTeam, o.game.teamName(winningTeam))
			}

			return
		}

		if msg, ok := parseTeamSetName(le); ok {
			o.setTeamname(msg.affiliation, msg.teamName)
		}

		return
	}

	// WarMod Warning
	if strings.HasPrefix(le.Message, "[WarMod_BFG]") {
		if strings.Contains(le.Message, `", "event": "log_start", `) {
			log.Warn().Msg("WarMod BFG detected; there are multiple bugs with running WarMod across multiple matches.")
		}
	}
}

// getTeam returns the team (mp_team1 / mp_team2 / unassigned)
// TODO: -- needs unit tests
func (o *observer) getTeam(aff affiliation) team {
	if aff != counterterrorist && aff != terrorist {
		return ""
	}

	mpHalftime, _ := o.srcdsObserver.TryCvarAsInt("mp_halftime", defaultMpHalftime)
	mpMaxrounds, _ := o.srcdsObserver.TryCvarAsInt("mp_maxrounds", defaultMpMaxrounds)
	mpOvertimeMaxrounds, _ := o.srcdsObserver.TryCvarAsInt("mp_overtime_maxrounds", defaultMpOvertimeMaxrounds)
	completedRounds := o.game.currentMatchLastCompletedRound()

	if calculateSidesAreCurrentlySwitched(mpHalftime, mpMaxrounds, mpOvertimeMaxrounds, completedRounds) {
		if aff == counterterrorist {
			return mpTeam2
		}

		return mpTeam1
	}

	if aff == counterterrorist {
		return mpTeam1
	}

	return mpTeam2
}

// TODO: -- needs unit tests
func (o *observer) playerDropped(c srcds.Client) {
	o.players.mpTeam1.ClientDropped(c)
	o.players.mpTeam2.ClientDropped(c)
	o.players.unassigned.ClientDropped(c)
	log.Info().Str("SteamID", c.SteamID).Msgf("Client %q disconnected.", c.Username)
}

// TODO: -- needs unit tests
func (o *observer) playerJoined(aff affiliation, c srcds.Client) {
	team := o.getTeam(aff)

	switch team {
	case mpTeam1:
		if o.players.mpTeam1.HasClient(c) {
			return
		}

		o.players.mpTeam1.ClientJoined(c)
		o.players.mpTeam2.ClientDropped(c)
		o.players.unassigned.ClientDropped(c)
	case mpTeam2:
		if o.players.mpTeam2.HasClient(c) {
			return
		}

		o.players.mpTeam1.ClientDropped(c)
		o.players.mpTeam2.ClientJoined(c)
		o.players.unassigned.ClientDropped(c)
	default:
		if o.players.unassigned.HasClient(c) {
			return
		}

		o.players.unassigned.ClientJoined(c)

		if o.players.mpTeam1.HasClient(c) {
			o.players.mpTeam1.ClientDropped(c)
			log.Info().Str("SteamID", c.SteamID).Msgf("Client %q dropped from mpTeam1 and joined unassigned.", c.Username)
			return
		}

		if o.players.mpTeam2.HasClient(c) {
			o.players.mpTeam2.ClientDropped(c)
			log.Info().Str("SteamID", c.SteamID).Msgf("Client %q dropped from mpTeam2 and joined unassigned.", c.Username)
			return
		}

		log.Debug().Str("SteamID", c.SteamID).Msgf("Client %q connected and joined unassigned.", c.Username)
		return
	}

	log.Info().Str("SteamID", c.SteamID).Msgf("Client %q joined %v.", c.Username, team)
}

func (o *observer) setTeamname(aff affiliation, name string) {
	team := o.getTeam(aff)
	name = strings.TrimSpace(name)

	switch team {
	case mpTeam1:
		if o.game.mpTeamname1 == name {
			return
		}

		o.game.mpTeamname1 = name
		log.Info().Msgf("Team %q is playing as %v (currently %v)", name, mpTeam1, aff)
	case mpTeam2:
		if o.game.mpTeamname2 == name {
			return
		}

		o.game.mpTeamname2 = name
		log.Info().Msgf("Team %q is playing as %v (currently %v)", name, mpTeam2, aff)
	default:
		log.Warn().Msgf("Cannot set a team name of %q for affiliation %q.", name, aff)
	}
}
