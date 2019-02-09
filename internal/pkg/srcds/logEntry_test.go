package srcds

import (
	"testing"
	"time"
)

func Test_ExtractLogEntry(t *testing.T) {
	testDatum := []struct {
		actualRaw    string
		expectedMsg  string
		expectedTime time.Time
	}{
		{"Server will auto-restart if there is a crash.", "", time.Time{}},
		{"L 1/2/2000 - 03:04:00: Sweet llamas of the Bahamas!", "Sweet llamas of the Bahamas!", time.Unix(946803840, 0)},
		{"L 01/2/2000 - 03:04:00: Excuse my language but I have had it with you ruffling my petticoats!", "Excuse my language but I have had it with you ruffling my petticoats!", time.Unix(946803840, 0)},
		{"L 1/02/2000 - 03:04:00: Your music is bad & you should feel bad!", "Your music is bad & you should feel bad!", time.Unix(946803840, 0)},
		{"L 01/02/2000 - 03:04:00: Did everything just taste purple for a second?", "Did everything just taste purple for a second?", time.Unix(946803840, 0)},
		{"L 01/02/2000 - 3:04:00: When you look this good, you don’t have to know anything!", "When you look this good, you don’t have to know anything!", time.Unix(946803840, 0)},
	}

	for _, testData := range testDatum {
		result := ExtractLogEntry(testData.actualRaw)

		if result.Message != testData.expectedMsg {
			t.Errorf("Expected message of %q but got %q", testData.expectedMsg, result.Message)
		}

		if result.Raw != testData.actualRaw {
			t.Errorf("Expected raw message of %q but got %q", testData.actualRaw, result.Raw)
		}

		if result.Timestamp != testData.expectedTime {
			t.Errorf("Expected timestamp of %q but got %q", testData.expectedTime, result.Timestamp)
		}
	}
}
