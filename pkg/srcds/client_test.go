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

var allFlags = [16]ClientFlag{
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
			},
			"TF2": {
				{c0: `"Betabot<2><[U:1:7609438]><>"`, c1: `"Betabot<26><[U:1:7609438]><Unassigned>"`},
				{c0: `"Betabot<2><[U:1:7609438]><Unassigned>"`, c1: `"Betabot<26><[U:1:7609438]><Unassigned>"`},
				{c0: `"Betabot<2><[U:1:7609438]><Unassigned>"`, c1: `"Betabot<26><[U:1:7609438]><Red>"`},
				{c0: `"Betabot<2><[U:1:7609438]><Unassigned>"`, c1: `"Betabot<26><[U:1:7609438]><Blue>"`},
				{c0: `"Betabot<2><[BOT]><>"`, c1: `"Betabot<26><[BOT]><Unassigned>"`},
				{c0: `"Betabot<2><[BOT]><Unassigned>"`, c1: `"Betabot<26><[BOT]><Unassigned>"`},
				{c0: `"Betabot<2><[BOT]><Unassigned>"`, c1: `"Betabot<26><[BOT]><Red>"`},
				{c0: `"Betabot<2><[BOT]><Unassigned>"`, c1: `"Betabot<26><[BOT]><Blue>"`},
				//TODO: NEED TF2 TV EXAMPLE
				//TODO: NEED TF2 CONSOLE EXAMPLE
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
				`"george pay attention to the call<0><STEAM_1:0:53845815><TERRORIST>"`,
			},
			"TF2": []string{
				`"The Masked Unit<2><[U:1:7609438]><>"`,
				`"Betabot<2><[U:1:7609438]><Unassigned>"`,
				`"Nurse Ratchet<2><[U:1:7609438]><Red>"`,
				`"Whalers on the Moon<3><[U:1:7609438]><Blue>"`,
				`"CreditToTeam<3><BOT><>"`,
				`"CreditToTeam<3><BOT><Unassigned>"`,
				`"CreditToTeam<3><BOT><Red>"`,
				`"CreditToTeam<3><BOT><Blue>"`,
				//TODO: NEED TF2 TV EXAMPLE
				//TODO: NEED TF2 CONSOLE EXAMPLE
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
		t.Run(fmt.Sprintf("%018b", f), func(t *testing.T) {
			// set & verify
			sut.EnableFlag(f)
			if !sut.HasFlag(f) {
				t.Errorf("Flag %018b should be set (client flags value was %018b).", f, sut.flags)
			}

			// verify no other flags got set
			for _, f2 := range allFlags {
				if f2 != f && sut.HasFlag(f2) {
					t.Errorf("Flag %018b should not be set (client flags value was %018b).", f2, sut.flags)
				}
			}

			// unset & verify
			sut.RemoveFlag(f)
			if sut.HasFlag(f) {
				t.Errorf("Flag %018b should NOT be set (client flags value was %018b).", f, sut.flags)
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
			"TF2": {
				`"Betabot<2><[BOT]><>"`,
				`"Betabot<2><[BOT]><Unassigned>"`,
				`"Betabot<26><[BOT]><Red>"`,
				`"Betabot<26><[BOT]><Blue>"`,
			},
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
			},
			"TF2": {
				`"The Masked Unit<2><[U:1:7609438]><>"`,
				`"Betabot<2><[U:1:7609438]><Unassigned>"`,
				`"Nurse Ratchet<2><[U:1:7609438]><Red>"`,
				`"Whalers on the Moon<3><[U:1:7609438]><Blue>"`,
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
			},
			"TF2": {
				`"Betabot<2><[U:1:7609438]><>"`,
				`"Betabot<2><[BOT]><>"`,
				`"Betabot<2><[BOT]><Unassigned>"`,
				`"Betabot<26><[BOT]><Red>"`,
				`"Betabot<26><[BOT]><Blue>"`,
				//TODO: NEED TF2 TV EXAMPLE
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

func Test_Client_RefreshEquivalentClient(t *testing.T) {
	sut := Clients{
		{Affiliation: "Blue", ServerSlot: 22, SteamID: "STEAM_1:0:1699142", Username: "Countess de la Roca"},
		{Affiliation: "ORIGINAL", ServerSlot: 1614, SteamID: "STEAM_1:0:9876543", Username: "Original Username"},
		{Affiliation: "Red", ServerSlot: 24, SteamID: "STEAM_1:0:1699143", Username: "Hedonism Bot"},
		{Affiliation: "CT", ServerSlot: 23, SteamID: "STEAM_1:0:1699144", Username: "Malfunctioning Eddie"},
	}

	mock := Client{Affiliation: "TERRORIST", ServerSlot: 123, SteamID: "STEAM_1:0:9876543", Username: "New Username"}

	sut.RefreshEquivalentClient(mock)

	i := sut.clientIndex(mock)

	if (sut)[i].Affiliation != mock.Affiliation {
		t.Errorf("Expected Affiliation %q but got %q.", mock.Affiliation, (sut)[i].Affiliation)
	}

	if (sut)[i].ServerSlot != mock.ServerSlot {
		t.Errorf("Expected ServerSlot %d but got %d.", mock.ServerSlot, (sut)[i].ServerSlot)
	}

	if (sut)[i].Username != mock.Username {
		t.Errorf("Expected Username %q but got %q.", mock.Username, (sut)[i].Username)
	}
}

func Test_Client_RemoveAllFlags(t *testing.T) {
	sut := Client{}

	// Turn on all the flags
	for _, f := range allFlags {
		sut.EnableFlag(f)
		if !sut.HasFlag(f) {
			t.Errorf("Was unable to enable flag %018b (client flags value was %018b).", f, sut.flags)
		}
	}

	sut.RemoveAllFlags()
	for _, f := range allFlags {
		if sut.HasFlag(f) {
			t.Errorf("Flag %018b should have been reset (client flags value was %018b).", f, sut.flags)
		}
	}
}

func Test_Client_ToggleFlags(t *testing.T) {
	sut := Client{}
	for _, f := range allFlags {
		t.Run(fmt.Sprintf("%018b", f), func(t *testing.T) {
			// toggle on
			sut.ToggleFlag(f)
			if !sut.HasFlag(f) {
				t.Errorf("Flag %018b should be set (client flags value was %018b).", f, sut.flags)
			}

			// verify no other flags got toggled on
			for _, f2 := range allFlags {
				if f2 != f && sut.HasFlag(f2) {
					t.Errorf("Flag %018b should not be set (client flags value was %018b).", f2, sut.flags)
				}
			}

			// toggle off
			sut.ToggleFlag(f)
			if sut.HasFlag(f) {
				t.Errorf("Flag %018b should NOT be set (client flags value was %018b).", f, sut.flags)
			}
		})
	}
}

func Test_Clients(t *testing.T) {
	clients := [4]Client{
		Client{Affiliation: "GROUP 1", ServerSlot: 1, SteamID: "STEAM_1:0:1699142", Username: "Countess de la Roca"},
		Client{Affiliation: "GROUP 2", ServerSlot: 2, SteamID: "STEAM_1:0:1699143", Username: "Hedonism Bot"},
		Client{Affiliation: "GROUP 1", ServerSlot: 3, SteamID: "STEAM_1:0:1699144", Username: "Malfunctioning Eddie"},
		Client{Affiliation: "GROUP 2", ServerSlot: 4, SteamID: "STEAM_1:0:1699145", Username: "Hair Robot"},
	}
	sut := Clients{}

	for i, c := range clients {
		if len(sut) != i {
			t.Errorf("Before client %q joined expected sut to have %d clients NOT %d.", c.Username, i, len(sut))
		}

		sut.ClientJoined(c)
		sut.ClientJoined(c) // Verify client doesn't get added twice
		if len(sut) != i+1 {
			t.Errorf("After client %q joined expected sut to have %d clients NOT %d.", c.Username, i, len(sut))
		}

		if !sut.HasClient(c) {
			t.Errorf("Client %q not found after being added.", c.Username)
		}
	}

	for i := len(sut) - 1; i >= 0; i-- {
		c := clients[i]
		sut.ClientDropped(c)
		if sut.HasClient(c) {
			t.Errorf("Client %q found after being dropped.", c.Username)
		}

		sut.ClientDropped(c) // Verify a second client doesn't get dropped
		if len(sut) != i {
			t.Errorf("After client %q dropped expected sut to have %d clients NOT %d.", c.Username, i, len(sut))
		}
	}
}

func Test_Clients_Flags(t *testing.T) {
	sut := Clients{
		Client{Affiliation: "EVEN", ServerSlot: 1, SteamID: "STEAM_1:0:1699142", Username: "Countess de la Roca"},
		Client{Affiliation: "ODD", ServerSlot: 2, SteamID: "STEAM_1:0:1699143", Username: "Hedonism Bot"},
		Client{Affiliation: "EVEN", ServerSlot: 3, SteamID: "STEAM_1:0:1699144", Username: "Malfunctioning Eddie"},
		Client{Affiliation: "ODD", ServerSlot: 4, SteamID: "STEAM_1:0:1699145", Username: "Hair Robot"},
	}

	t.Run("Before Enabling", func(t *testing.T) {
		l := len(sut.WithFlags(clientFlagAlpha, clientFlagBravo))
		if l != 0 {
			t.Errorf("Expected NO clients to show with the flags %018b and %018b but got %d.", clientFlagAlpha, clientFlagBravo, l)
		}

		l = len(sut.WithoutFlags(clientFlagAlpha, clientFlagNovember))
		if l != len(sut) {
			t.Errorf("Expected ALL clients to show without the flags %018b and %018b but got %d.", clientFlagAlpha, clientFlagNovember, l)
		}
	})

	t.Run("EnableFlag", func(t *testing.T) {
		for i := range sut {
			sut.EnableFlag(sut[i], clientFlagAlpha, clientFlagCharlie, clientFlagEcho, clientFlagHotel, clientFlagJuliett)

			if !sut[i].HasFlag(clientFlagCharlie) {
				t.Errorf("1) After enabling flag %018b for %q it should have returned true for HasFlag()", clientFlagAlpha, sut[i].Username)
			}
		}

		l := len(sut.WithFlags(clientFlagAlpha, clientFlagBravo))
		if l != 0 {
			t.Errorf("4) Expected NO clients to show with the flags %018b and %018b but got %d.", clientFlagAlpha, clientFlagBravo, l)
		}

		l = len(sut.WithoutFlags(clientFlagBravo, clientFlagDelta))
		if l != len(sut) {
			t.Errorf("6) Expected ALL clients to show without the flags %018b and %018b but got %d", clientFlagBravo, clientFlagDelta, l)
		}
	})

	t.Run("RemoveFlag", func(t *testing.T) {
		for i := range sut {
			sut.RemoveFlag(sut[i], clientFlagAlpha, clientFlagBravo)

			if sut[i].HasFlag(clientFlagAlpha) {
				t.Errorf("Expected client %q to NOT have flag %018b enabled.", sut[i].Username, clientFlagAlpha)
			}
		}

		l := len(sut.WithFlags(clientFlagAlpha))
		if l != 0 {
			t.Errorf("Expected NO clients to show WITH the flag %018b but got %d.", clientFlagAlpha, l)
		}

		l = len(sut.WithoutFlags(clientFlagAlpha))
		if l != len(sut) {
			t.Errorf("Expected ALL clients to show WITHOUT the flag %018b but got %d.", clientFlagAlpha, l)
		}
	})

	t.Run("RemoveFlags", func(t *testing.T) {
		sut.RemoveFlags(clientFlagCharlie, clientFlagHotel)

		l := len(sut.WithFlags(clientFlagCharlie, clientFlagHotel))
		if l != 0 {
			t.Errorf("Expected NO clients to show WITH the flags %018b and %018b but got %d.", clientFlagCharlie, clientFlagHotel, l)
		}

		l = len(sut.WithoutFlags(clientFlagCharlie, clientFlagHotel))
		if l != len(sut) {
			t.Errorf("Expected ALL clients to show WITHOUT the flags %018b and %018b but got %d.", clientFlagCharlie, clientFlagHotel, l)
		}
	})

	t.Run("RemoveAllFlags", func(t *testing.T) {
		sut.RemoveAllFlags()

		l := len(sut.WithFlags(clientFlagEcho, clientFlagJuliett))
		if l != 0 {
			t.Errorf("Expected NO clients to show WITH the flags %018b and %018b but got %d.", clientFlagEcho, clientFlagJuliett, l)
		}

		l = len(sut.WithoutFlags(clientFlagEcho, clientFlagJuliett))
		if l != len(sut) {
			t.Errorf("Expected ALL clients to show WITHOUT the flags %018b and %018b but got %d.", clientFlagEcho, clientFlagJuliett, l)
		}
	})
}
