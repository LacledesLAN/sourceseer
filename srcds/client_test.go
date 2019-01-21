package srcds

import (
	"testing"
)

func Test_ClientsAreEquivalent(t *testing.T) {
	testDatum := []struct {
		name           string
		client0        *Client
		client1        *Client
		expectedResult bool
	}{
		{"`nil` clients should not match", nil, nil, false},
		{"Defaulted clients should not match", &Client{}, &Client{}, false},
		{"Empty client should not match", &Client{Username: "", SteamID: ""}, &Client{Username: "", SteamID: ""}, false},
		{"`nil` should not match with an actual client", &Client{Username: "Mark 7-G"}, nil, false},
		{"Defaulted client should not match with an actual client", &Client{}, &Client{Username: "Mark 7-G"}, false},
		{"Empty client should not match with an actual client", &Client{Username: "", SteamID: ""}, &Client{Username: "Mark 7-G"}, false},
		{"Same Username and no SteamID should match", &Client{Username: "John Quincy Adding Machine"}, &Client{Username: "John Quincy Adding Machine"}, true},
		{"Different Username should not match", &Client{Username: "Macaulay Culkon"}, &Client{Username: "Dr. Widnar"}, false},
		{"Matching Username and SteamId should match", &Client{Username: "iZac", SteamID: "ph1l l4m4rr"}, &Client{Username: "iZac", SteamID: "ph1l l4m4rr"}, true},
		{"Matching Username but different SteamID should not match", &Client{Username: "iZac", SteamID: "l4m4rr"}, &Client{Username: "iZac", SteamID: "ph1l"}, false},
		{"Different Username but matching SteamID should not match", &Client{Username: "ABC", SteamID: "123"}, &Client{Username: "DEF", SteamID: "123"}, false},
	}

	for _, testData := range testDatum {
		t.Run(testData.name, func(t *testing.T) {
			actualResult := ClientsAreEquivalent(testData.client0, testData.client1)

			if actualResult != testData.expectedResult {
				t.Errorf("Test %q failed; expected '%t' but got '%t'.", testData.name, testData.expectedResult, actualResult)
			}
		})
	}
}
