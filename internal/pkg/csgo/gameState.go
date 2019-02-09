package csgo

import (
	"regexp"
)

var (
	regexBetweenQuotes = regexp.MustCompile("(?:\")([0-9A-Za-z_]*)(\")")
)
