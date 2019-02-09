package csgo

import "regexp"

const (
	gameOverPattern                   = `^Game Over: (\w+)[ ]+(\w+) score (\d+):(\d+) after (\d+) min$`
	matchStartPattern                 = `^World triggered "Match_Start" on "(\w+)"$`
	teamScoredPattern                 = `^Team "(CT|TERRORIST)" scored "(\d+)" with "(\d+)" players$`
	teamSetSidePattern                = `^Team playing "(.{1,})": (.{1,})$`
	worldTriggeredMatchStartPattern   = `^World triggered "Match_Start" on "(\w+)"$`
	worldTriggeredRoundRestartPattern = `^World triggered "Restart_Round_\((\d+)_second\)`
)

var (
	gameOverRegex                   = regexp.MustCompile(gameOverPattern)
	matchStartRegex                 = regexp.MustCompile(matchStartPattern)
	teamScoredRegex                 = regexp.MustCompile(teamScoredPattern)
	teamSetSideRegex                = regexp.MustCompile(teamSetSidePattern)
	worldTriggeredMatchStartRegex   = regexp.MustCompile(worldTriggeredMatchStartPattern)
	worldTriggeredRoundRestartRegex = regexp.MustCompile(worldTriggeredRoundRestartPattern)
)
