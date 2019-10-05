package srcds

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

// Cvar represents a watched SRCDS console variable
type Cvar struct {
	LastUpdated time.Time
	Value       string
	seededValue bool
}

// Cvars represents a collection of watched console variables
type Cvars struct {
	v   map[string]Cvar
	mux sync.Mutex
}

func (c *Cvars) addWatcher(names ...string) {
	if len(names) < 0 {
		return
	}

	c.mux.Lock()
	if c.v == nil {
		c.v = make(map[string]Cvar)
	}

	for _, name := range names {
		name = strings.TrimSpace(name)

		if len(name) == 0 {
			continue
		}

		if _, found := c.v[name]; !found {
			c.v[name] = Cvar{}
		}
	}
	c.mux.Unlock()
}

func (c *Cvars) getNames() []string {
	c.mux.Lock()
	defer c.mux.Unlock()

	r := make([]string, 0, len(c.v))

	for name := range c.v {
		r = append(r, name)
	}

	return r
}

func (c *Cvars) seedWatcher(name, value string) {
	name = strings.TrimSpace(name)

	if len(name) == 0 {
		return
	}

	if c.v == nil {
		c.v = make(map[string]Cvar)
	}

	c.mux.Lock()
	if cvar, found := c.v[name]; !found {
		c.v[name] = Cvar{Value: strings.TrimSpace(value), seededValue: true}
	} else if cvar.LastUpdated.IsZero() {
		// seed the value only if it hasn't been naturally found
		c.v[name] = Cvar{Value: strings.TrimSpace(value), seededValue: true}
	}
	c.mux.Unlock()
}

func (c *Cvars) setIfWatched(name, value string, asOf time.Time) {
	name = strings.TrimSpace(name)

	if len(name) == 0 {
		return
	}

	value = strings.TrimSpace(value)

	if asOf.IsZero() {
		asOf = time.Now()
	}

	c.mux.Lock()
	if _, found := c.v[name]; found {
		c.v[name] = Cvar{
			LastUpdated: asOf,
			Value:       strings.TrimSpace(value),
			seededValue: false,
		}
	}
	c.mux.Unlock()
}

func (c *Cvars) tryFloat(name string, fallback float32) (value float32, nonFallback bool) {
	str, nonFallback := c.tryString(name, "")

	if !nonFallback {
		return fallback, nonFallback
	}

	f, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return fallback, false
	}

	return float32(f), true
}

func (c *Cvars) tryInt(name string, fallback int) (value int, nonFallback bool) {
	str, nonFallback := c.tryString(name, "")

	if !nonFallback {
		return fallback, nonFallback
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		return fallback, false
	}

	return i, true
}

func (c *Cvars) tryString(name, fallback string) (value string, nonFallback bool) {
	if c == nil {
		return fallback, false
	}

	c.mux.Lock()
	cvar, found := c.v[name]
	c.mux.Unlock()

	if found && (cvar.seededValue || !cvar.LastUpdated.IsZero()) {
		return cvar.Value, true
	}

	return fallback, false
}
