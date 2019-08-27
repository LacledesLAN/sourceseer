package srcds

import (
	"fmt"
	"testing"
)

func Test_ClientFlagRegistry_Register(t *testing.T) {
	t.Parallel()

	const (
		testFlagAlpha ClientFlag = 1 << iota
		testFlagBeta
	)

	sut := ClientFlagRegistry{}

	sut.Register(testFlagAlpha, "", "   ", "\n", "\r\n", "\t", "\t\n\t")
	if len(sut.registry) != 0 {
		t.Error("Should not be able to add empty trigger")
	}

	sut.Register(testFlagAlpha, "a", "acquire", "", "analyze", "  ", "azotize", "\n")
	if len(sut.registry) != 4 {
		t.Errorf("Expected 4 flags to be registered but got %d.", len(sut.registry))
	}

	sut.Register(testFlagAlpha, "bravo", "BRAVO", " brAvo  ")
	sut.Register(testFlagBeta, "Bravo", "bRaVo", "BRAVO\n")
	if len(sut.registry) != 5 {
		t.Error("Should not be able to register the same flag twice.")
	}
}

func Test_ClientFlagRegistry_Find(t *testing.T) {
	t.Parallel()

	const (
		testFlagAlpha ClientFlag = 1 << iota
		testFlagBeta
	)

	sut := ClientFlagRegistry{}
	sut.Register(testFlagAlpha, "a", "alpha", "")
	sut.Register(testFlagBeta, "b", "BRAVO", "\n")

	if _, found := sut.Find(""); found {
		t.Errorf("An empty string should never be found")
	}

	validCases := map[ClientFlag][]string{
		testFlagAlpha: []string{"a", "A", "alpha", "ALPHA", "aLPHA", "Alpha"},
		testFlagBeta:  []string{"b", "B", "Bravo ", " BRAVO"},
	}

	for f, tests := range validCases {
		t.Run(fmt.Sprintf("%018b", f), func(t *testing.T) {
			for _, test := range tests {
				if actual, matched := sut.Find(test); !matched {
					t.Errorf("String %q should have returned an associated flag", test)
				} else if actual != f {
					t.Errorf("Expected flag %018b but got %018b.", f, actual)
				}
			}
		})
	}
}
