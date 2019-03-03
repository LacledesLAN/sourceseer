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
	clientSwitchedAffiliationRegex = regexp.MustCompile(clientSwitchedAffiliationPattern)
	gameOverRegex                  = regexp.MustCompile(gameOverPattern)
	loadingMapRegex                = regexp.MustCompile(loadingMapPattern)
	playerSayRegex                 = regexp.MustCompile(playerSayPattern)
	teamScoredRegex                = regexp.MustCompile(teamScoredPattern)
	teamSetSideRegex               = regexp.MustCompile(teamSetSidePattern)
	teamTriggeredRegex             = regexp.MustCompile(teamTriggeredPattern)
	worldTriggeredRegex            = regexp.MustCompile(worldTriggeredPattern)
)

//ClientSwitchedAffiliation message is sent whenever a client switches affiliation (aka teams)
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

//PlayerSaidChannel represents the channel in which a player sends a message
type PlayerSaidChannel int

const (
	//ChannelGlobal is seen by everyone in the csgo server
	ChannelGlobal PlayerSaidChannel = iota + 1
	//ChannelAffiliation is seen by anyone in the relevant team
	ChannelAffiliation
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
		r.channel = ChannelAffiliation
	default:
		r.channel = ChannelGlobal
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
