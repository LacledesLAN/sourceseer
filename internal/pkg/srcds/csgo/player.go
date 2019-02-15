package csgo

import (
	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

// Player is a srcds client that is actively playing csgo
type Player struct {
	srcds.Client
	IsReady bool
}

func playerFromSrcdsClient(c srcds.Client) Player {
	return Player{
		Client:  c,
		IsReady: false}
}
