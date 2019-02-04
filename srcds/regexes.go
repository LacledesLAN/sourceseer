package srcds

import "regexp"

const (
	extractPlayerInfoPattern = `"(.{1,32})<(\d{0,1})><([a-zA-Z0-9:_]*)><{0,1}([a-zA-Z0-9]*?)>{0,1}" ([^[\-?\d+ -?\d+ -?\d+\]]?)`
	serverCvarEchoPattern    = `^"([^\s]{3,})" = "(.*)"$`
	serverCvarSetPattern     = `^server_cvar: "([^\s]{3,})" "(.*)"$`
)

var (
	extractPlayerInfoRegex = regexp.MustCompile(extractPlayerInfoPattern) //group 1: name; group 2: slot; group 3: uid; group 4: team (if any)
	serverCvarEchoRegex    = regexp.MustCompile(serverCvarEchoPattern)
	serverCvarSetRegex     = regexp.MustCompile(serverCvarSetPattern)
)
