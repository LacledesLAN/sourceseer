package csgo

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

const (
	clientConnectedPattern           = `^(".*") (?:entered the game|connected, address "")$`
	clientDisconnectedPattern        = `^(".*") disconnected(?: \(reason \"([\w ]{1,})\"\))?$`
	clientSwitchedAffiliationPattern = `^(".*") switched from team <([a-zA-Z]*)> to <([a-zA-Z]*)>$`
	gameOverPattern                  = `^Game Over: (\w+)[ ]+(\w+) score (\d+):(\d+) after (\d+) min$`
	loadingMapPattern                = `^Loading map "(\w+)"$`
	playerSayPattern                 = `^(".*") (say_team|say) "(.+)"$`
	teamScoredPattern                = `^Team "(CT|TERRORIST)" scored "(\d+)" with "(\d+)" players$`
	teamSetSidePattern               = `^Team playing "(.{1,})": (.{1,})$`
	teamTriggeredPattern             = `^Team "(CT|TERRORIST)" triggered \"(SFUI_Notice_[A-Za-z_]{4,34})\" \(CT \"([\d]{1,3})\"\) \(T \"([\d]{1,3})\"\)$`
	worldTriggeredPattern            = `^World triggered \"(\w+)(?:_\((\d+)_seconds\))?\"(?: on \"(\w+)\")?$`
)

var (
	clientConnectedRegex           = regexp.MustCompile(clientConnectedPattern)
	clientDisconnectedRegex        = regexp.MustCompile(clientDisconnectedPattern)
	clientSwitchedAffiliationRegex = regexp.MustCompile(clientSwitchedAffiliationPattern)
	gameOverRegex                  = regexp.MustCompile(gameOverPattern)
	loadingMapRegex                = regexp.MustCompile(loadingMapPattern)
	playerSayRegex                 = regexp.MustCompile(playerSayPattern)
	teamScoredRegex                = regexp.MustCompile(teamScoredPattern)
	teamSetSideRegex               = regexp.MustCompile(teamSetSidePattern)
	teamTriggeredRegex             = regexp.MustCompile(teamTriggeredPattern)
	worldTriggeredRegex            = regexp.MustCompile(worldTriggeredPattern)
)

type ClientDisconnected struct {
	client srcds.Client
	reason string
}

type ClientSwitchedAffiliation struct {
	client srcds.Client
	from   string
	to     string
}

// PlayerSaid message is sent whenever a player says something using their console
type PlayerSaid struct {
	player  Player
	channel PlayerSaidChannel
	msg     string
}

type PlayerSaidChannel int

const (
	GlobalChannel PlayerSaidChannel = iota + 1
	TeamChannel
)

// TeamScored message is sent whenever a team wins a round
type TeamScored struct {
	teamAffiliation string
	teamScore       int
	teamPlayerCount int
}

// TeamSideSet message is sent whenever team information is set
type TeamSideSet struct {
	teamAffiliation string
	teamName        string
}

// TeamTriggered message is sent whenever a team wins a rounds
type TeamTriggered struct {
	teamAffiliation string
	triggered       string
	ctScore         int
	terroristScore  int
}

// WorldTriggered message is sent whenever the world experience a change
type WorldTriggered struct {
	trigger WorldTrigger
	eta     time.Time
}

func parseClientConnected(le srcds.LogEntry) (srcds.Client, error) {
	result := clientConnectedRegex.FindStringSubmatch(le.Message)

	if len(result) != 2 {
		return srcds.Client{}, errors.New("Could not client " + le.Message)
	}

	cl, err := srcds.ParseClient(result[1])

	if err != nil {
		return srcds.Client{}, errors.New("Could not parse client in " + le.Message)
	}

	return cl, nil
}

func parseClientDisconnected(le srcds.LogEntry) (ClientDisconnected, error) {
	r := clientDisconnectedRegex.FindStringSubmatch(le.Message)

	if len(r) < 1 {
		return ClientDisconnected{}, errors.New("Could not parse " + le.Message)
	}

	cl, err := srcds.ParseClient(r[1])

	if err != nil {
		return ClientDisconnected{}, errors.New("Could not parse player in " + le.Message)
	}

	return ClientDisconnected{
		client: cl,
		reason: r[2],
	}, nil
}

func parseClientSwitchedAffiliation(le srcds.LogEntry) (ClientSwitchedAffiliation, error) {
	r := clientSwitchedAffiliationRegex.FindStringSubmatch(le.Message)

	if len(r) != 4 {
		return ClientSwitchedAffiliation{}, errors.New("Could not parse " + le.Message)
	}

	cl, err := srcds.ParseClient(r[1])

	if err != nil {
		return ClientSwitchedAffiliation{}, errors.New("Could not parse player in " + le.Message)
	}

	return ClientSwitchedAffiliation{
		client: cl,
		from:   r[2],
		to:     r[3],
	}, nil
}

func parseLoadingMap(le srcds.LogEntry) (string, error) {
	r := loadingMapRegex.FindStringSubmatch(le.Message)

	if len(r) != 2 {
		return "", errors.New("Could not parse:" + le.Message)
	}

	return r[1], nil
}

func parsePlayerSay(le srcds.LogEntry) (PlayerSaid, error) {
	sayTokens := playerSayRegex.FindStringSubmatch(le.Message)

	if len(sayTokens) != 4 {
		return PlayerSaid{}, errors.New("Could not parse " + le.Message)
	}

	cl, err := srcds.ParseClient(sayTokens[1])

	if err != nil {
		return PlayerSaid{}, errors.New("Could not parse player in " + le.Message)
	}

	r := PlayerSaid{
		player: playerFromSrcdsClient(cl),
		msg:    sayTokens[3],
	}

	switch strings.ToUpper(sayTokens[2]) {
	case "SAY_TEAM":
		r.channel = TeamChannel
	default:
		r.channel = GlobalChannel
	}

	return r, nil
}

func parseTeamScored(le srcds.LogEntry) (TeamScored, error) {
	result := teamScoredRegex.FindStringSubmatch(le.Message)

	if len(result) != 4 {
		return TeamScored{}, errors.New("Could not parse " + le.Message)
	}

	var err error
	r := TeamScored{
		teamAffiliation: strings.ToUpper(result[1]),
	}

	r.teamScore, err = strconv.Atoi(result[2])
	r.teamPlayerCount, err = strconv.Atoi(result[3])

	if err != nil {
		return TeamScored{}, err
	}

	return r, nil
}

func parseTeamSetSide(le srcds.LogEntry) (TeamSideSet, error) {
	result := teamSetSideRegex.FindStringSubmatch(le.Message)

	if len(result) != 3 {
		return TeamSideSet{}, errors.New("Could not parse " + le.Message)
	}

	return TeamSideSet{
		teamAffiliation: strings.ToUpper(result[1]),
		teamName:        result[2],
	}, nil
}

func parseTeamTriggered(le srcds.LogEntry) (TeamTriggered, error) {
	result := teamTriggeredRegex.FindStringSubmatch(le.Message)

	if len(result) != 5 {
		return TeamTriggered{}, errors.New("Could not parse " + le.Message)
	}

	var err error
	r := TeamTriggered{
		teamAffiliation: result[1],
		triggered:       result[2],
	}

	r.ctScore, err = strconv.Atoi(result[3])
	r.terroristScore, err = strconv.Atoi(result[4])

	if err != nil {
		return TeamTriggered{}, err
	}

	return r, nil
}

func parseWorldTriggered(le srcds.LogEntry) (WorldTriggered, error) {
	result := worldTriggeredRegex.FindStringSubmatch(le.Message)

	if len(result) != 4 {
		return WorldTriggered{}, errors.New("Could not parse " + le.Message)
	}

	r := WorldTriggered{
		trigger: UnknownTrigger,
		eta:     time.Now(),
	}

	switch strings.ToUpper(result[1]) {
	case "ROUND_START":
		r.trigger = RoundStart
	case "ROUND_END":
		r.trigger = RoundEnd
	case "RESTART_ROUND":
		r.trigger = RoundRestarting

		if i, err := strconv.Atoi(result[2]); err == nil && i > 0 {
			r.eta = time.Now().Add(time.Duration(i) * time.Second)
		}
	case "MATCH_START":
		r.trigger = MatchStart
	case "GAME_COMMENCING":
		r.trigger = GameCommencing
	default:
		return WorldTriggered{}, errors.New("Unknown world trigger: " + result[1])
	}

	return r, nil
}
