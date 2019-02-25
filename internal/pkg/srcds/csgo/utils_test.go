package csgo

import (
	"testing"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

func Test_CalculateWinThreshold(t *testing.T) {
	testDatum := []struct {
		mpMaxRounds         int
		mpOvertimeMaxRounds int
		lastCompletedRound  int
		expectedResult      int
	}{
		// normal rounds clinchable
		{1, 2, -1, 1}, // impossible condition
		{1, 2, 0, 1}, {1, 2, 1, 1},
		{1, 2, 2, 1}, // impossible condition

		// normal rounds not clinchable; OT rounds clinchable
		{2, 3, -1, 2}, // impossible condition
		{2, 3, 0, 2}, {2, 3, 1, 2}, {2, 3, 2, 2},
		{2, 3, 3, 3}, {2, 3, 4, 3},
		{2, 3, 5, 3}, // impossible condition

		// normal rounds not clinchable; OT rounds not clinchable
		{2, 2, -1, 2}, // impossible condition
		{2, 2, 0, 2}, {2, 2, 1, 2}, {2, 2, 2, 2},
		{2, 2, 3, 3}, {2, 2, 4, 3}, // OT 1
		{2, 2, 5, 4}, {2, 2, 6, 4}, // OT 2
		{2, 2, 7, 5}, {2, 2, 8, 5}, // OT 3
		{2, 2, 115, 59}, {2, 2, 116, 59}, // OT 56

		// "Hasty" server settings
		{4, 3, 0, 3}, {4, 3, 1, 3}, {4, 3, 2, 3}, {4, 3, 3, 3}, {4, 3, 4, 3},
		{4, 3, 5, 4}, {4, 3, 6, 4}, {4, 3, 5, 4}, // OT

		// default server settings
		{30, 6, -878, 16}, {30, 6, -1, 16}, //7possible conditions
		{30, 6, 0, 16}, {30, 6, 1, 16}, {30, 6, 2, 16}, {30, 6, 3, 16}, {30, 6, 4, 16}, {30, 6, 5, 16}, {30, 6, 6, 16}, {30, 6, 7, 16}, {30, 6, 8, 16}, {30, 6, 9, 16},
		{30, 6, 10, 16}, {30, 6, 11, 16}, {30, 6, 12, 16}, {30, 6, 13, 16}, {30, 6, 14, 16}, {30, 6, 15, 16}, {30, 6, 16, 16}, {30, 6, 17, 16}, {30, 6, 18, 16},
		{30, 6, 19, 16}, {30, 6, 20, 16}, {30, 6, 21, 16}, {30, 6, 22, 16}, {30, 6, 23, 16}, {30, 6, 24, 16}, {30, 6, 25, 16}, {30, 6, 26, 16}, {30, 6, 27, 16},
		{30, 6, 28, 16}, {30, 6, 29, 16}, {30, 6, 30, 16},
		{30, 6, 31, 19}, {30, 6, 32, 19}, {30, 6, 33, 19}, {30, 6, 34, 19}, {30, 6, 35, 19}, {30, 6, 36, 19}, // OT 1
		{30, 6, 37, 22}, {30, 6, 38, 22}, {30, 6, 39, 22}, {30, 6, 40, 22}, {30, 6, 41, 22}, {30, 6, 42, 22}, // OT 2
		{30, 6, 49, 28}, {30, 6, 55, 31}, {30, 6, 61, 34}, {30, 6, 67, 37}, {30, 6, 73, 40}, {30, 6, 79, 43},
		{30, 6, 91, 49}, {30, 6, 92, 49}, {30, 6, 93, 49}, {30, 6, 94, 49}, {30, 6, 95, 49}, {30, 6, 96, 49},

		// prevailing community tourney settings
		{30, 7, -878, 16}, {30, 7, -1, 16}, // impossible conditions
		{30, 7, 0, 16}, {30, 7, 1, 16}, {30, 7, 2, 16}, {30, 7, 3, 16}, {30, 7, 4, 16}, {30, 7, 5, 16}, {30, 7, 6, 16}, {30, 7, 7, 16}, {30, 7, 8, 16}, {30, 7, 9, 16},
		{30, 7, 10, 16}, {30, 7, 11, 16}, {30, 7, 12, 16}, {30, 7, 13, 16}, {30, 7, 14, 16}, {30, 7, 15, 16}, {30, 7, 16, 16}, {30, 7, 17, 16}, {30, 7, 18, 16},
		{30, 7, 19, 16}, {30, 7, 20, 16}, {30, 7, 21, 16}, {30, 7, 22, 16}, {30, 7, 23, 16}, {30, 7, 24, 16}, {30, 7, 25, 16}, {30, 7, 26, 16}, {30, 7, 27, 16},
		{30, 7, 28, 16}, {30, 7, 29, 16}, {30, 7, 30, 16},
		{30, 7, 31, 19}, {30, 7, 32, 19}, {30, 7, 33, 19}, {30, 7, 34, 19}, {30, 7, 35, 19}, {30, 7, 36, 19}, {30, 7, 37, 19}, // OT
		{30, 7, 38, 19}, {30, 7, 540, 19}, // impossible conditions
	}

	for _, d := range testDatum {
		actual := calculateWinThreshold(d.mpMaxRounds, d.mpOvertimeMaxRounds, d.lastCompletedRound)

		if actual != d.expectedResult {
			t.Errorf("With `mp_maxrounds` = `%d` and `mp_overtime_maxrounds` = `%d` and the last completed round being `%d` a calculated win threshold should be %d not %d!", d.mpMaxRounds, d.mpOvertimeMaxRounds, d.lastCompletedRound, d.expectedResult, actual)
		}
	}
}

func Test_ValidateStockMapName(t *testing.T) {
	testDatum := []struct {
		mapName     string
		expectError bool
	}{
		{"ar_baggage", false}, {"ar_dizzy", false}, {"ar_monastery", false}, {"ar_shoots", false}, {"cs_agency", false}, {"cs_assault", false}, {"cs_militia", false}, {"cs_italy", false}, {"cs_office", false}, {"de_austria", false},
		{"de_bank", false}, {"de_biome", false}, {"de_cache", false}, {"de_canals", false}, {"de_cbble", false}, {"de_dust2", false}, {"de_inferno", false}, {"de_lake", false}, {"de_mirage", false}, {"de_nuke", false},
		{"de_overpass", false}, {"de_safehouse", false}, {"de_shortnuke", false}, {"de_stmarc", false}, {"de_subzero", false}, {"de_train", false}, {"", true}, {"   ", true}, {"\t", true}, {"\n", true}, {"garbage", true},
	}

	for _, testData := range testDatum {
		errResult := validateStockMapName(testData.mapName)

		if testData.expectError && errResult == nil {
			t.Errorf("Map '%q' should have resulted in an error", testData.mapName)
		} else if !testData.expectError && errResult != nil {
			t.Errorf("Map '%q' should have not have resulted in an error; but got error %q.", testData.mapName, errResult)
		}
	}
}

func Test_ValidateStockMapNames(t *testing.T) {
	testDatum := []struct {
		maps        []string
		expectError bool
	}{
		{[]string{}, true},
		{[]string{"ar_baggage", "de_safehouse"}, false},
		{[]string{"ar_baggage", "de_nope"}, true},
		{[]string{"ar_baggage", "", "de_train"}, true},
		{[]string{"\n", "ar_baggage", "de_train"}, true},
	}

	for _, testData := range testDatum {
		errResult := validateStockMapNames(testData.maps)

		if testData.expectError && errResult == nil {
			t.Errorf("Maplist '%q' should have resulted in an error", testData.maps)
		} else if !testData.expectError && errResult != nil {
			t.Errorf("Maplist '%q' should have not have resulted in an error; but got error %q.", testData.maps, errResult)
		}
	}
}

func Test_HostnameFromTeamNames(t *testing.T) {
	testDatum := []struct {
		mpTeamname1    string
		mpTeamname2    string
		expectedResult string
	}{
		// both team names in range
		{"a", "b", "a-vs-b"},
		{"aaa", "bbb", "aaa-vs-bbb"},
		{"aaaaaaaaaaaa", "bbbbbbbbbbbb", "aaaaaaaaaaaa-vs-bbbbbbbbbbbb"},
		// both teams out of range
		{"aaaaaaaaaaaaAAAA", "bbbbbbbbbbbbBBBB", "aaaaaaaaaaaa-vs-bbbbbbbbbbbb"},
		// team 1 out of range
		{"aaaaaaaaaaaaA", "bbbbbbbbbbbb", "aaaaaaaaaaaa-vs-bbbbbbbbbbbb"},
		{"aaaaaaaaaaaaAAAAAAAAAaaa", "bbb", "aaaaaaaaaaaaAAAAAAAAA-vs-bbb"},
		{"aaaaaaaaaaaaAAAAAAAAAAAaaa", "b", "aaaaaaaaaaaaAAAAAAAAAAA-vs-b"},
		// team 2 out of range
		{"bbbbbbbbbbbb", "aaaaaaaaaaaaA", "bbbbbbbbbbbb-vs-aaaaaaaaaaaa"},
		{"bbb", "aaaaaaaaaaaaAAAAAAAAAaaa", "bbb-vs-aaaaaaaaaaaaAAAAAAAAA"},
		{"b", "aaaaaaaaaaaaAAAAAAAAAAAaaa", "b-vs-aaaaaaaaaaaaAAAAAAAAAAA"},
		// whitespace clean up
		{"  aaa ", " bbb  ", "aaa-vs-bbb"},
		{" a a ", " b b ", "a_a-vs-b_b"},
		{" a  a ", " b  b ", "a_a-vs-b_b"},
		{" a\ta ", " b\nb ", "a_a-vs-b_b"},
	}

	for _, testData := range testDatum {
		result := HostnameFromTeamNames(testData.mpTeamname1, testData.mpTeamname2)

		if len(result) > srcds.MaxHostnameLength {
			t.Errorf("TOO LONG! - Teams %q and %q should expected to be %q but got %q", testData.mpTeamname1, testData.mpTeamname2, testData.expectedResult, result)
		} else if testData.expectedResult != result {
			t.Errorf("Teams %q and %q should expected to be %q but got %q", testData.mpTeamname1, testData.mpTeamname2, testData.expectedResult, result)
		}
	}
}

func Test_SanitizeTeamName(t *testing.T) {
	testDatum := []struct {
		original string
		expected string
	}{
		{"", ""},
		{"   ", ""},
		{" \n \n \t \r\n   ", ""},
		{"_ ", "_"},
		{" -", "-"},
		{" a b c ", "a_b_c"},
		{"@abc", "abc"},
	}

	for _, testData := range testDatum {
		result := SanitizeTeamName(testData.original)

		if result != testData.expected {
			t.Errorf("String %q should have been sanitized to %q but instead became %q.", testData.original, testData.expected, result)
		}
	}
}
