package srcds

import "regexp"

const (
	extractClientPattern  = `^"(.{1,32})<(\d{0,1})><([\w:]*)><{0,1}([a-zA-Z0-9]*?)>{0,1}"`
	extractClientsPattern = `"(.{1,32})<(\d{0,1})><([a-zA-Z0-9:_]*)><{0,1}([a-zA-Z0-9]*?)>{0,1}" ([^[\-?\d+ -?\d+ -?\d+\]]?)`
	serverCvarEchoPattern = `^"([^\s]{3,})" = "([^\s"]{0,})"`
	serverCvarSetPattern  = `^server_cvar: "([^\s]{3,})" "(.*)"$`
)

var (
	extractClientRegex  = regexp.MustCompile(extractClientPattern)
	extractClientsRegex = regexp.MustCompile(extractClientsPattern) //group 1: name; group 2: slot; group 3: uid; group 4: team (if any)
	serverCvarEchoRegex = regexp.MustCompile(serverCvarEchoPattern)
	serverCvarSetRegex  = regexp.MustCompile(serverCvarSetPattern)
)
