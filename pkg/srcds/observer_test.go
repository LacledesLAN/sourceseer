package srcds

import (
	"bufio"
	"os"
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

func printableEOL(s string) string {
	switch s {
	case eolUnix:
		return "Unix EOL"
	case eolWindows:
		return "Windows EOL"
	case "\n\r":
		return "RISC"
	default:
		return "Unknown EOL"
	}
}

func Test_NewObserver(t *testing.T) {
	tests := []struct {
		filename      string
		expectedStats observerStatistics
		expectedEOL   string
	}{
		{
			filename:      `./testdata/simple.log`,
			expectedStats: observerStatistics{totalLines: 13, blankLines: 3, logLines: 3},
			expectedEOL:   eolUnix,
		},
		{
			filename:      `./testdata/simple.win.log`,
			expectedStats: observerStatistics{totalLines: 13, blankLines: 3, logLines: 3},
			expectedEOL:   eolWindows,
		},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			fh, err := os.Open(test.filename)
			defer fh.Close()
			if err != nil {
				t.Fatalf("Could not open file: %q.", err)
			}

			r := bufio.NewReader(fh)

			sut := NewObserver()
			sut.Read(r)
			sut.Wait()

			if sut.statistics != test.expectedStats {
				t.Errorf("Statistics did not meet expectations; expected %+v but got %+v", test.expectedStats, sut.statistics)
			}

			if sut.EndOfLine != test.expectedEOL {
				t.Errorf("Expected end of line %q but got %q", printableEOL(test.expectedEOL), printableEOL(sut.EndOfLine))
			}
		})
	}
}
