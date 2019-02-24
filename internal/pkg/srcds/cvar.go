package srcds

import "time"

type Cvar struct {
	LastUpdate time.Time
	Value      string
}
