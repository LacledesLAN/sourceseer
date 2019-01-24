package csgo

import (
	"github.com/lacledeslan/sourceseer/srcds"
)

type csgoClient struct {
	srcds.Client
	isReady bool
}
