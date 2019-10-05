package csgo

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

const (
	alphaNumericChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lacledesMaps      = "/de_lltest/de_tinyorange/poolday/"
	stockMaps         = "/ar_baggage/ar_dizzy/ar_monastery/ar_shoots/cs_agency/cs_assault/cs_italy/cs_militia/cs_office/de_austria/de_bank/de_biome/de_cache/de_canals/de_cbble/de_dust2/de_inferno/de_lake/de_mirage/de_nuke/de_overpass/de_safehouse/de_shortnuke/de_stmarc/de_subzero/de_sugarcane/de_train/"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TourneyServer(mpTeamname1, mpTeamname2, pass, rconPass, tvPass string, maps ...string) error {
	mpTeamname1 = strings.TrimSpace(mpTeamname1)
	if len(mpTeamname1) == 0 {
		return errors.New("mpTeamname1 cannot be empty")
	}

	mpTeamname2 = strings.TrimSpace(mpTeamname2)
	if len(mpTeamname2) == 0 {
		return errors.New("mpTeamname2 cannot be empty")
	}

	if strings.ToLower(mpTeamname1) == strings.ToLower(mpTeamname2) {
		return errors.New("mpTeamname1 and mpTeamname2 cannot match")
	}

	pass = strings.TrimSpace(pass)
	rconPass = strings.TrimSpace(rconPass)

	tvPass = strings.TrimSpace(tvPass)
	if len(tvPass) < 1 {
		b := make([]byte, 16)
		for i := range b {
			b[i] = alphaNumericChars[rand.Intn(len(alphaNumericChars))]
		}
		tvPass = string(b)
	}

	return nil
}
