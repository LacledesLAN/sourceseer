package srcds

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_Cvars(t *testing.T) {
	const fallbackFloat = float32(10.77)
	const fallbackInt = 1077
	const fallbackString = "No more implants. I don't want to end up a cold, emotionless machine like you."

	type checkFunc func(cvars *Cvars, name string) error
	check := func(funcs ...checkFunc) []checkFunc { return funcs }

	expectFoundInGetNames := func() checkFunc {
		return func(c *Cvars, name string) error {
			for _, k := range c.getNames() {
				if k == name {
					return nil
				}
			}
			return fmt.Errorf("Cvar %q should have been found in GetNames", name)
		}
	}

	expectNotFoundInGetNames := func() checkFunc {
		return func(c *Cvars, name string) error {
			for _, k := range c.getNames() {
				if k == name {
					return fmt.Errorf("Cvar %q should NOT have been found in GetNames", name)
				}
			}
			return nil
		}
	}

	expectFloat := func(expected float32) checkFunc {
		return func(c *Cvars, name string) error {
			if actual, nonFallback := c.tryFloat(name, fallbackFloat); actual != expected {
				return fmt.Errorf("Cvar %q as float should have returned %f not %f", name, expected, actual)
			} else if !nonFallback {
				return fmt.Errorf("Cvar %q was successfully returned float - should have also returned true not false", name)
			}
			return nil
		}
	}

	expectFallbackFloat := func() checkFunc {
		return func(c *Cvars, name string) error {
			if actual, nonFallback := c.tryFloat(name, fallbackFloat); actual != fallbackFloat {
				return fmt.Errorf("Cvar %q as float should have returned the fallback value not %f", name, actual)
			} else if nonFallback {
				return fmt.Errorf("Cvar %q did not successfully return a float - should have also returned false not true", name)
			}
			return nil
		}
	}

	expectInt := func(expected int) checkFunc {
		return func(c *Cvars, name string) error {
			if actual, nonFallback := c.tryInt(name, fallbackInt); actual != expected {
				return fmt.Errorf("Cvar %q as int should have returned %d not %d", name, expected, actual)
			} else if !nonFallback {
				return fmt.Errorf("Cvar %q was successfully returned int - should have also returned true not false", name)
			}
			return nil
		}
	}

	expectFallbackInt := func() checkFunc {
		return func(c *Cvars, name string) error {
			if actual, nonFallback := c.tryInt(name, fallbackInt); actual != fallbackInt {
				return fmt.Errorf("Cvar %q as int should have returned the fallback value not %d", name, actual)
			} else if nonFallback {
				return fmt.Errorf("Cvar %q did not successfully return a int - should have also returned false not true", name)
			}
			return nil
		}
	}

	expectString := func(expected string) checkFunc {
		return func(c *Cvars, name string) error {
			if actual, nonFallback := c.tryString(name, fallbackString); actual != expected {
				return fmt.Errorf("Cvar %q as string should have returned %q not %q", name, expected, actual)
			} else if !nonFallback {
				return fmt.Errorf("Cvar %q was successfully returned string - should have also returned true not false", name)
			}
			return nil
		}
	}

	expectFallbackString := func() checkFunc {
		return func(c *Cvars, name string) error {
			if actual, nonFallback := c.tryString(name, fallbackString); actual != fallbackString {
				return fmt.Errorf("Cvar %q as string should have returned the fallback value not %q", name, actual)
			} else if nonFallback {
				return fmt.Errorf("Cvar %q did not successfully return a string - should have also returned false not true", name)
			}
			return nil
		}
	}

	tests := map[string]struct {
		v      string
		checks []checkFunc
	}{
		"": {
			v:      "I'm an empty string",
			checks: check(expectNotFoundInGetNames(), expectFallbackFloat(), expectFallbackInt(), expectFallbackString()),
		},
		"a_int": {
			v:      "123",
			checks: check(expectFoundInGetNames(), expectFloat(float32(123)), expectInt(123), expectString("123")),
		},
		"a_float": {
			v:      "78.09",
			checks: check(expectFoundInGetNames(), expectFloat(float32(78.09)), expectFallbackInt(), expectString("78.09")),
		},
		"a_string": {
			v:      "four five six",
			checks: check(expectFoundInGetNames(), expectFallbackFloat(), expectFallbackInt(), expectString("four five six")),
		},
	}

	t.Run("Add Watcher", func(t *testing.T) {
		i := 0
		if rand.Float32() < 0.5 {
			i++
		}

		for name, test := range tests {
			i++
			t.Run(name, func(t *testing.T) {
				mockTime := time.Now()
				if i%2 == 0 {
					mockTime = time.Time{}
				}

				sut := &Cvars{}
				sut.addWatcher(name)
				sut.setIfWatched(name, test.v, mockTime)

				for _, check := range test.checks {
					if err := check(sut, name); err != nil {
						t.Error(err)
					}
				}
			})
		}
	})

	t.Run("Seed Watcher", func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				sut := &Cvars{}
				sut.seedWatcher(name, test.v)

				for _, check := range test.checks {
					if err := check(sut, name); err != nil {
						t.Error(err)
					}
				}
			})
		}
	})

	t.Run("Seeding, Setting, and Retrieval", func(t *testing.T) {
		sut := &Cvars{}
		sut.addWatcher("alpha")

		if _, nonFallback := sut.tryString("alpha", "fallback value"); nonFallback {
			t.Error("A cvar added for watching, that never had a value set, should return the fallback value on retrieval.")
		}

		sut.seedWatcher("alpha", "beta")
		if cvar, _ := sut.tryString("alpha", "fallback"); cvar != "beta" {
			t.Error("A cvar added for watching, that never had a value set but was seeded, should return the seeded value on retrieval.")
		}

		sut.seedWatcher("alpha", "charlie")
		if cvar, _ := sut.tryString("alpha", "fallback"); cvar != "charlie" {
			t.Error("Seeding should overwrite any previously-seeded values.")
		}

		if len(sut.getNames()) > 1 {
			t.Error("Should not have duplicate cvars with the same name.")
		}

		sut.setIfWatched("alpha", "delta", time.Now())
		if cvar, _ := sut.tryString("alpha", "fallback"); cvar != "delta" {
			t.Error("Setting a value on a watched cvar should overwrite any seeded value.")
		}

		sut.seedWatcher("alpha", "echo")
		if cvar, _ := sut.tryString("alpha", "fallback"); cvar != "delta" {
			t.Errorf("Seeding a set cvar should not overwrite any set value.")
		}
	})
}
