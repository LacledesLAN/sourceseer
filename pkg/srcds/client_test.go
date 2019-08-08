package srcds

import (
	"fmt"
	"testing"
)

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
				//TODO: NEED TV EXAMPLES
				//TODO: NEED BOT EXAMPLES
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

func Test_Client_Flags(t *testing.T) {
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

	allFlags := []ClientFlag{clientFlagAlpha, clientFlagBravo, clientFlagCharlie, clientFlagDelta, clientFlagEcho, clientFlagGolf, clientFlagHotel, clientFlagIndia,
		clientFlagJuliett, clientFlagKilo, clientFlagLima, clientFlagMike, clientFlagNovember, clientFlagOscar, clientFlagPapa, clientFlagQuedec}

	t.Run("Enable / Remove", func(t *testing.T) {
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
	})

	t.Run("Toggle", func(t *testing.T) {
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
	})

	t.Run("Reset", func(t *testing.T) {
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
	})
}
