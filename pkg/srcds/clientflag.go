package srcds

import (
	"strings"
	"sync"
)

// ClientFlag represents a flag set for a client
type ClientFlag uint16

// ClientFlagRegistry associates strings with ClientFlags
type ClientFlagRegistry struct {
	mux      sync.Mutex
	registry map[string]ClientFlag
}

// Find determines if a case-insensitive string is associated with a ClientFlag
func (r *ClientFlagRegistry) Find(s string) (flag ClientFlag, isAssociated bool) {
	s = strings.TrimSpace(s)

	if len(s) < 1 {
		return 0, false
	}

	r.mux.Lock()
	f, found := r.registry[strings.ToLower(s)]
	r.mux.Unlock()
	return f, found
}

// Register case-insensitive string associations with the specified ClientFlag
func (r *ClientFlagRegistry) Register(f ClientFlag, s string, ss ...string) {
	r.mux.Lock()

	if len(r.registry) == 0 {
		r.registry = make(map[string]ClientFlag, 1+len(ss))
	}

	s = strings.TrimSpace(s)
	if len(s) > 0 {
		r.registry[strings.ToLower(s)] = f
	}

	for _, s := range ss {
		s = strings.TrimSpace(s)

		if len(s) < 1 {
			continue
		}

		r.registry[strings.ToLower(s)] = f
	}

	r.mux.Unlock()
}
