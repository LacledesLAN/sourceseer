package csgo

import (
	"testing"

	"github.com/lacledeslan/sourceseer/srcds"
)

func TestValidateStockMapName(t *testing.T) {
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

func TestValidateStockMapNames(t *testing.T) {
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
