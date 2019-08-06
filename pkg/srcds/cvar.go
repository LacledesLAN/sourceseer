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
}

// Cvars represents a collection of watched console variables
type Cvars struct {
	v   map[string]Cvar
	mux sync.Mutex
}

func (c *Cvars) addWatcher(names ...string) {
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
	r := make([]string, 0, len(c.v))

	for name := range c.v {
		r = append(r, name)
	}

	c.mux.Unlock()
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
		c.v[name] = Cvar{Value: strings.TrimSpace(value)}
	} else if cvar.LastUpdated.IsZero() {
		// seed the value only if it hasn't been naturally found
		c.v[name] = Cvar{Value: strings.TrimSpace(value)}
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
		}
	}
	c.mux.Unlock()
}

func (c *Cvars) tryFloat(name string, fallback float32) (value float32, nonFallback bool) {
	c.mux.Lock()
	cvar, found := c.v[name]
	c.mux.Unlock()

	if !found {
		return fallback, false
	}

	f, err := strconv.ParseFloat(cvar.Value, 32)
	if err != nil {
		return fallback, false
	}

	return float32(f), true
}

func (c *Cvars) tryInt(name string, fallback int) (value int, nonFallback bool) {
	c.mux.Lock()
	cvar, found := c.v[name]
	c.mux.Unlock()

	if !found {
		return fallback, false
	}

	i, err := strconv.Atoi(cvar.Value)

	if err != nil {
		return fallback, false
	}

	return i, true
}

func (c *Cvars) tryString(name, fallback string) (value string, nonFallback bool) {
	c.mux.Lock()
	cvar, found := c.v[name]
	c.mux.Unlock()

	if !found {
		return fallback, false
	}

	return cvar.Value, true
}
