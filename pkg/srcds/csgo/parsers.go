package csgo

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/LacledesLAN/sourceseer/pkg/srcds"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func parseAffiliation(af string) (affiliation affiliation, ok bool) {
	switch a := strings.ToUpper(strings.TrimSpace(af)); a {
	case "":
		return unassigned, true
	case "CT":
		return counterterrorist, true
	case "TERRORIST":
		return terrorist, true
	case "UNASSIGNED":
		return unassigned, true
	}

	return unassigned, false
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var playerSayRegex = regexp.MustCompile(`^(say_team|say) "(.+)"$`)

type sayChannel int

const (
	//ChannelGlobal is seen by everyone in the csgo server
	ChannelGlobal sayChannel = iota + 1
	//ChannelAffiliation is seen by anyone in the relevant team
	ChannelAffiliation
)

// clientSaid is sent whenever a player sends a text message using the `say` command
type clientSaid struct {
	channel sayChannel
	msg     string
}

func parseClientSay(clientLog srcds.ClientLogEntry) (clientSaid, bool) {
	tokens := playerSayRegex.FindStringSubmatch(clientLog.Message)

	if len(tokens) != 3 {
		return clientSaid{}, false
	}

	r := clientSaid{
		msg:     tokens[2],
		channel: ChannelGlobal,
	}

	if tokens[1] == "say_team" {
		r.channel = ChannelAffiliation
	}

	return r, true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var clientSetAffiliationRegex = regexp.MustCompile(`^switched from team <([a-zA-Z]*)> to <([a-zA-Z]*)>$`)

// clientSwitchedAffiliation is sent whenever a client/player changes their affiliated team
type clientSetAffiliation struct {
	to   affiliation
	from affiliation
}

func parseClientSetAffiliation(clientLog srcds.ClientLogEntry) (clientSetAffiliation, bool) {
	tokens := clientSetAffiliationRegex.FindStringSubmatch(clientLog.Message)

	if len(tokens) != 3 {
		return clientSetAffiliation{}, false
	}

	r := clientSetAffiliation{}
	r.to, _ = parseAffiliation(tokens[2])
	r.from, _ = parseAffiliation(tokens[1])

	return r, true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var gameOverRegex = regexp.MustCompile(`^Game Over: (\w+)[ ]+(\w+) score (\d+):(\d+) after (\d+) min$`)

// gameOver is (sometimes?) sent when the game is over
type gameOver struct {
	mode           string
	mapName        string
	score1         int
	score2         int
	minutesElapsed int
}

func parseGameOver(le srcds.LogEntry) (m gameOver, ok bool) {
	tokens := gameOverRegex.FindStringSubmatch(le.Message)

	if len(tokens) != 6 {
		return gameOver{}, false
	}

	r := gameOver{}
	r.minutesElapsed, _ = strconv.Atoi(tokens[5])
	r.score2, _ = strconv.Atoi(tokens[4])
	r.score1, _ = strconv.Atoi(tokens[3])
	r.mapName = tokens[2]
	r.mode = tokens[1]

	return r, true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var loadingMapRegex = regexp.MustCompile(`^Loading map "(\w+)"$`)

func parseLoadingMap(le srcds.LogEntry) (mapName string, ok bool) {
	tokens := loadingMapRegex.FindStringSubmatch(le.Message)

	if len(tokens) != 2 {
		return "", false
	}

	return tokens[1], true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func parseStartingFreezePeriod(le srcds.LogEntry) (ok bool) {
	return strings.HasPrefix(le.Message, `Starting Freeze period`)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var teamScoredRegex = regexp.MustCompile(`^Team "(CT|TERRORIST)" scored "(\d+)" with "(\d+)" players$`)

// teamScored is sent when either the Counter-Terrorist or the Terrorist win a round
type teamScored struct {
	affiliation affiliation
	Score       int
	PlayerCount int
}

func parseTeamScored(le srcds.LogEntry) (teamScored, bool) {
	tokens := teamScoredRegex.FindStringSubmatch(le.Message)

	if len(tokens) != 4 {
		return teamScored{}, false
	}

	r := teamScored{}
	r.PlayerCount, _ = strconv.Atoi(tokens[3])
	r.Score, _ = strconv.Atoi(tokens[2])
	r.affiliation, _ = parseAffiliation(tokens[1])

	return r, true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var teamSetNameRegex = regexp.MustCompile(`^Team playing "(.{1,})": (.{1,})$`)

// teamSetName is sent whenever team information is set
type teamSetName struct {
	affiliation affiliation
	teamName    string
}

func parseTeamSetName(le srcds.LogEntry) (teamSetName, bool) {
	tokens := teamSetNameRegex.FindStringSubmatch(le.Message)

	if len(tokens) != 3 {
		return teamSetName{}, false
	}

	r := teamSetName{}
	r.teamName = strings.TrimSpace(tokens[2])
	r.affiliation, _ = parseAffiliation(tokens[1])

	return r, true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var teamTriggeredRegex = regexp.MustCompile(`^Team "(CT|TERRORIST)" triggered \"(SFUI_Notice_[A-Za-z_]{4,34})\" \(CT \"([\d]{1,4})\"\) \(T \"([\d]{1,4})\"\)$`)

// TeamTriggered is sent when a team wins a rounds
type teamTriggered struct {
	affiliation    affiliation
	trigger        string
	ctScore        int
	terroristScore int
}

func parseTeamTriggered(le srcds.LogEntry) (teamTriggered, bool) {
	tokens := teamTriggeredRegex.FindStringSubmatch(le.Message)

	if len(tokens) != 5 {
		return teamTriggered{}, false
	}

	r := teamTriggered{}

	r.terroristScore, _ = strconv.Atoi(tokens[4])
	r.ctScore, _ = strconv.Atoi(tokens[3])
	r.trigger = tokens[2]
	r.affiliation, _ = parseAffiliation(tokens[1])

	return r, true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var worldTriggerRegex = regexp.MustCompile(`^World triggered ("[\S]+".*)`)

type logWorldTrigger string

func parseWorldTrigger(le srcds.LogEntry) (msg logWorldTrigger, ok bool) {
	tokens := worldTriggerRegex.FindStringSubmatch(le.Message)

	if len(tokens) != 2 {
		return "", false
	}

	return logWorldTrigger(tokens[1]), true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func parseWorldTriggerGameCommencing(msg logWorldTrigger) (ok bool) {
	return strings.HasPrefix(string(msg), `"Game_Commencing"`)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func parseWorldTriggerRoundEnd(msg logWorldTrigger) (ok bool) {
	return strings.HasPrefix(string(msg), `"Round_End"`)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func parseWorldTriggerRoundStart(msg logWorldTrigger) (ok bool) {
	return strings.HasPrefix(string(msg), `"Round_Start"`)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var worldTriggerMatchStartRegex = regexp.MustCompile(`^"Match_Start" on "([\w]*)"$`)

func parseWorldTriggerMatchStart(msg logWorldTrigger) (mapName string, ok bool) {
	tokens := worldTriggerMatchStartRegex.FindStringSubmatch(string(msg))

	if len(tokens) != 2 {
		return "", false
	}

	return tokens[1], true
}
