package srcds

import (
	"fmt"
	"testing"
)

const (
	clientFlagAlpha ClientFlag = 1 << iota
	clientFlagBravo
	clientFlagCharlie
	clientFlagDelta
	clientFlagEcho
	clientFlagGolf
	clientFlagHotel
	clientFlagIndia
	clientFlagJuliett
	clientFlagKilo
	clientFlagLima
	clientFlagMike
	clientFlagNovember
	clientFlagOscar
	clientFlagPapa
	clientFlagQuedec
)

var allFlags = []ClientFlag{
	clientFlagAlpha, clientFlagBravo, clientFlagCharlie, clientFlagDelta, clientFlagEcho, clientFlagGolf, clientFlagHotel, clientFlagIndia,
	clientFlagJuliett, clientFlagKilo, clientFlagLima, clientFlagMike, clientFlagNovember, clientFlagOscar, clientFlagPapa, clientFlagQuedec,
}

func Test_ClientsAreEquivalent(t *testing.T) {
	t.Run("Should Be Equivalent", func(t *testing.T) {
		equivalentTests := map[string][]struct {
			c0 string
			c1 string
		}{
			"CSGO": {
				{c0: `"GOTV<2><BOT><>"`, c1: `"GOTV<3><BOT><>"`},
				{c0: `"Console<0><Console><Console>"`, c1: `"Console<0><Console><Console>"`},
				{c0: `"The Masked Unit<2><[STEAM_1:0:53045815]><>"`, c1: `"The Masked Unit<2><[STEAM_1:0:53045815]><>"`},
				{c0: `"The Masked Unit<2><[STEAM_1:0:53045815]><>"`, c1: `"The Masked Unit<99><STEAM_1:0:53045815]><>"`},
				{c0: `"The Masked Unit<2><[STEAM_1:0:53045815]><>"`, c1: `"The Renamed Unit<2><STEAM_1:0:53045815]><>"`},
				{c0: `"The Masked Unit<2><[STEAM_1:0:53045815]><>"`, c1: `"The Renamed Unit<2><STEAM_1:0:53045815]><Unassigned>"`},
				{c0: `"The Masked Unit<2><[STEAM_1:0:53045815]><>"`, c1: `"The Renamed Unit<2><STEAM_1:0:53045815]><CT>"`},
				{c0: `"The Masked Unit<2><[STEAM_1:0:53045815]><>"`, c1: `"The Renamed Unit<2><STEAM_1:0:53045815]><TERRORIST>"`},
				{c0: `"The Masked Unit<2><[STEAM_1:0:53045815]><CT>"`, c1: `"The Renamed Unit<2><[STEAM_1:0:53045815]><TERRORIST>"`},
				{c0: `"The Masked Unit<2><[STEAM_1:0:53045815]><CT>"`, c1: `"The Renamed Unit<2><[STEAM_1:0:53045815]><Unassigned>"`},
				{c0: `"The Masked Unit<2><[STEAM_1:0:53045815]><TERRORIST>"`, c1: `"The Renamed Unit<2><[STEAM_1:0:53045815]><Unassigned>"`},
				//TODO: NEED SERVER OFFLINE / CLIENT OFFLINE EXAMPLES
			},
			"TF2": {
				{c0: `"Betabot<2><[U:1:7609438]><>"`, c1: `"Betabot<26><[U:1:7609438]><Unassigned>"`},
				{c0: `"Betabot<2><[U:1:7609438]><Unassigned>"`, c1: `"Betabot<26><[U:1:7609438]><Unassigned>"`},
				{c0: `"Betabot<2><[U:1:7609438]><Unassigned>"`, c1: `"Betabot<26><[U:1:7609438]><Red>"`},
				{c0: `"Betabot<2><[U:1:7609438]><Unassigned>"`, c1: `"Betabot<26><[U:1:7609438]><Blue>"`},
				//TODO: NEED BOT EXAMPLES
				//TODO: NEED CONSOLE EXAMPLES
				//TODO: NEED TV EXAMPLES
				//TODO: NEED SERVER OFFLINE / CLIENT OFFLINE EXAMPLES
			},
		}

		for name, tests := range equivalentTests {
			t.Run(name, func(t *testing.T) {
				for _, test := range tests {
					if c0, ok := ParseClient(test.c0); !ok {
						t.Fatalf("Could not parse client from %q", test.c0)
					} else {
						if c1, ok := ParseClient(test.c1); !ok {
							t.Fatalf("Could not parse client from %q", test.c1)
						} else if !ClientsAreEquivalent(c0, c1) {
							t.Errorf("Clients %v and %v should have been considered equivalent.", test.c0, test.c1)
						}
					}
				}
			})
		}
	})

	t.Run("Should NOT be Equivalent", func(t *testing.T) {
		notEquivalentTests := map[string][]struct {
			c0 string
			c1 string
		}{
			"CSGO": {
				{c0: `"Nameless<13><STEAM_1:0:00000000><>"`, c1: `"Nameless<12><STEAM_1:0:99999999><>"`},
				{c0: `"Jack<2><BOT><>"`, c1: `"Jill<2><BOT><>"`},
			},
			"TF2": {
				{c0: `"Betabot<2><[U:1:0000000]><Unassigned>"`, c1: `"Betabot<2><[U:1:9999999]><Unassigned>"`},
			},
		}

		for name, tests := range notEquivalentTests {
			t.Run(name, func(t *testing.T) {
				for _, test := range tests {
					if c0, ok := ParseClient(test.c0); !ok {
						t.Fatalf("Could not parse client from %q", test.c0)
					} else {
						if c1, ok := ParseClient(test.c1); !ok {
							t.Fatalf("Could not parse client from %q", test.c1)
						} else if ClientsAreEquivalent(c0, c1) {
							t.Errorf("Clients %v and %v should NOT have been considered equivalent.", test.c0, test.c1)
						}
					}
				}
			})
		}
	})

}

func Test_ClientUnidentifiable(t *testing.T) {
	t.Run("Should be identifiable", func(t *testing.T) {
		tests := map[string][]string{
			"CSGO": []string{
				`"Console<0><Console><Console>"`,
				`"GOTV<2><BOT><>"`,
				`"GOTV<2><BOT><Unassigned>"`,
				`"Boxy Robot<13><STEAM_1:0:53045815><>"`,
				`"Countess de la Roca<6><STEAM_1:0:53045815><CT>"`,
				`"doku pay attention to the call<0><STEAM_1:0:53045815><TERRORIST>"`,
				//TODO: NEED SERVER OFFLINE / CLIENT OFFLINE EXAMPLES
			},
			"TF2": []string{
				`"The Masked Unit<2><[U:1:7609438]><>"`,
				`"Betabot<2><[U:1:7609438]><Unassigned>"`,
				`"Nurse Ratchet<2><[U:1:7609438]><Red>"`,
				`"Whalers on the Moon<3><[U:1:7609438]><Blue>"`,
				//TODO: NEED TF2 TV EXAMPLE
				//TODO: NEED TF2 BOT EXAMPLES
				//TODO: NEED TF2 CONSOLE EXAMPLE
				//TODO: NEED SERVER OFFLINE / CLIENT OFFLINE EXAMPLES
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				for _, str := range test {
					if c, ok := ParseClient(str); !ok {
						t.Fatalf("Could not parse client from %q", str)
					} else {
						if ClientUnidentifiable(c) {
							t.Errorf("Client %+v should have been identifiable.", c)
						}
					}
				}
			})
		}
	})
}

func Test_Client_EnableFlag(t *testing.T) {
	sut := Client{}
	for _, f := range allFlags {
		t.Run(fmt.Sprintf("%016b", f), func(t *testing.T) {
			// set & verify
			sut.EnableFlag(f)
			if !sut.HasFlag(f) {
				t.Errorf("Flag %016b should be set (client flags value was %016b).", f, sut.flags)
			}

			// verify no other flags got set
			for _, f2 := range allFlags {
				if f2 != f && sut.HasFlag(f2) {
					t.Errorf("Flag %016b should not be set (client flags value was %016b).", f2, sut.flags)
				}
			}

			// unset & verify
			sut.RemoveFlag(f)
			if sut.HasFlag(f) {
				t.Errorf("Flag %016b should NOT be set (client flags value was %016b).", f, sut.flags)
			}
		})
	}
}

func Test_Client_IsBot(t *testing.T) {
	t.Run("Valid Cases", func(t *testing.T) {
		validCases := map[string][]string{
			"CSGO": {
				`"GOTV<3><BOT><>"`,
				`"John<6><BOT><>"`,
				`"Jim<2><BOT><CT>"`,
				`"Joe<9><BOT><TERRORIST>"`,
			},
			"TF2": {},
			//TODO: NEED TF2 BOT EXAMPLES
		}

		for name, test := range validCases {
			t.Run(name, func(t *testing.T) {
				for _, str := range test {
					if c, ok := ParseClient(str); !ok {
						t.Fatalf("Could not parse client from %q", str)
					} else if !c.IsBot() {
						t.Errorf("%v should be a bot", str)
					}
				}
			})
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := map[string][]string{
			"CSGO": {
				`"Console<0><Console><Console>"`,
				`"Boxy Robot<13><STEAM_1:0:53045815><>"`,
				`"Boxy Robot<13><STEAM_1:0:53045815><CT>"`,
				`"Boxy Robot<13><STEAM_1:0:53045815><TERRORIST>"`,
				`"Boxy Robot<13><STEAM_1:0:53045815><UNASSIGNED>"`,
				//TODO: NEED SERVER OFFLINE / CLIENT OFFLINE EXAMPLES
			},
			"TF2": {
				`"The Masked Unit<2><[U:1:7609438]><>"`,
				`"Betabot<2><[U:1:7609438]><Unassigned>"`,
				`"Nurse Ratchet<2><[U:1:7609438]><Red>"`,
				`"Whalers on the Moon<3><[U:1:7609438]><Blue>"`,
				//TODO: NEED SERVER OFFLINE / CLIENT OFFLINE EXAMPLES
				//TODO: NEED TF2 CONSOLE EXAMPLE
			},
		}

		for name, test := range invalidCases {
			t.Run(name, func(t *testing.T) {
				for _, str := range test {
					if c, ok := ParseClient(str); !ok {
						t.Fatalf("Could not parse client from %q", str)
					} else if c.IsBot() {
						t.Errorf("%v should NOT be a bot", str)
					}
				}
			})
		}
	})
}

func Test_Client_IsConsole(t *testing.T) {
	t.Run("Valid Cases", func(t *testing.T) {
		validCases := map[string][]string{
			"CSGO": {
				`"Console<0><Console><Console>"`,
			},
			"TF2": {},
			//TODO: NEED TF2 CONSOLE EXAMPLES
		}

		for name, test := range validCases {
			t.Run(name, func(t *testing.T) {
				for _, str := range test {
					if c, ok := ParseClient(str); !ok {
						t.Fatalf("Could not parse client from %q", str)
					} else if !c.IsConsole() {
						t.Errorf("%v should be a console", str)
					}
				}
			})
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidCases := map[string][]string{
			"CSGO": {
				`"GOTV<3><BOT><>"`,
				`"The Masked Unit<2><[STEAM_1:0:53045815]><>"`,
				`"Tony<11><BOT><>"`,
				//TODO: NEED SERVER OFFLINE / CLIENT OFFLINE EXAMPLES
			},
			"TF2": {
				`"Betabot<2><[U:1:7609438]><>"`,
				//TODO: NEED TF2 TV EXAMPLE
				//TODO: NEED TF2 BOT EXAMPLES
				//TODO: NEED SERVER OFFLINE / CLIENT OFFLINE EXAMPLES
			},
		}

		for name, test := range invalidCases {
			t.Run(name, func(t *testing.T) {
				for _, str := range test {
					if c, ok := ParseClient(str); !ok {
						t.Fatalf("Could not parse client from %q", str)
					} else if c.IsConsole() {
						t.Errorf("%v should NOT be a console", str)
					}
				}
			})
		}
	})
}

func Test_Client_RemoveAllFlags(t *testing.T) {
	sut := Client{}

	// Turn on all the flags
	for _, f := range allFlags {
		sut.EnableFlag(f)
		if !sut.HasFlag(f) {
			t.Errorf("Was unable to enable flag %016b (client flags value was %016b).", f, sut.flags)
		}
	}

	sut.RemoveAllFlags()
	for _, f := range allFlags {
		if sut.HasFlag(f) {
			t.Errorf("Flag %016b should have been reset (client flags value was %016b).", f, sut.flags)
		}
	}
}

func Test_Client_ToggleFlags(t *testing.T) {
	sut := Client{}
	for _, f := range allFlags {
		t.Run(fmt.Sprintf("%016b", f), func(t *testing.T) {
			// toggle on
			sut.ToggleFlag(f)
			if !sut.HasFlag(f) {
				t.Errorf("Flag %016b should be set (client flags value was %016b).", f, sut.flags)
			}

			// verify no other flags got toggled on
			for _, f2 := range allFlags {
				if f2 != f && sut.HasFlag(f2) {
					t.Errorf("Flag %016b should not be set (client flags value was %016b).", f2, sut.flags)
				}
			}

			// toggle off
			sut.ToggleFlag(f)
			if sut.HasFlag(f) {
				t.Errorf("Flag %016b should NOT be set (client flags value was %016b).", f, sut.flags)
			}
		})
	}
}

func Test_Clients(t *testing.T) {
	//client joined
	//has client
	//client dropped
	//has client
}

func Test_Client_RemoveFlags(t *testing.T) {
	//TODO!
}

func Test_Clients_WithFlags(t *testing.T) {
	//TODO!
}

func Test_Clients_WithoutFlags(t *testing.T) {
	//TODO!
}
