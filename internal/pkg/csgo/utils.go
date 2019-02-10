package csgo

import (
	"errors"
	"math"
	"regexp"
	"strings"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

const (
	// List of valid stock map names
	validMaps = "/ar_baggage/ar_dizzy/ar_monastery/ar_shoots/cs_agency/cs_assault/cs_italy/cs_militia/cs_office/de_austria/de_bank/de_biome/de_cache/de_canals/de_cbble/de_dust2/de_inferno/de_tinyorange/de_lake/de_mirage/de_nuke/de_overpass/de_safehouse/de_shortnuke/de_stmarc/de_subzero/de_sugarcane/de_train/"
)

var (
	srcdsSafeChars = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
)

// CalculateWinThreshold determines how many rounds a team needs to win to win a map given how many rounds have been completed so far
func CalculateWinThreshold(mpMaxRounds, mpOvertimeMaxRounds, lastCompletedRound int) int {
	if mpMaxRounds < 1 {
		mpMaxRounds = 30
	}

	if mpOvertimeMaxRounds < 1 {
		mpMaxRounds = 6
	}

	winThreshold := mpMaxRounds/2 + 1

	if totalOTRoundsCompleted := lastCompletedRound - mpMaxRounds; totalOTRoundsCompleted > 0 {
		otPeriodsCompleted := (totalOTRoundsCompleted - 1) / mpOvertimeMaxRounds
		return winThreshold + ((mpOvertimeMaxRounds/2 + 1) * (otPeriodsCompleted + 1))
	}

	return winThreshold
}

// HostnameFromTeamNames generates a hostname for srcds from two teamnames
func HostnameFromTeamNames(mpTeamname1 string, mpTeamname2 string) string {
	mpTeamname1 = SanitizeTeamName(mpTeamname1)
	mpTeamname2 = SanitizeTeamName(mpTeamname2)

	if len(mpTeamname1) == 0 {
		if len(mpTeamname2) == 0 {
			return "CSGO Tourney Server"
		}

		mpTeamname1 = "Unspecified"
	}

	if len(mpTeamname2) == 0 {
		mpTeamname2 = "Unspecified"
	}

	glue := "-vs-"
	hostname := mpTeamname1 + glue + mpTeamname2

	if len(hostname) <= srcds.MaxHostnameLength {
		return hostname
	}

	maxTeamNameLength := int(math.Floor(float64(srcds.MaxHostnameLength-len(glue)) / 2))

	if len(mpTeamname1) > maxTeamNameLength && len(mpTeamname2) > maxTeamNameLength {
		return mpTeamname1[:maxTeamNameLength] + glue + mpTeamname2[:maxTeamNameLength]
	}

	if len(mpTeamname1) <= maxTeamNameLength {
		hostname = mpTeamname1 + glue
		remainingLen := srcds.MaxHostnameLength - len(hostname)

		return hostname + mpTeamname2[:remainingLen]
	}

	hostname = glue + mpTeamname2
	remainingLen := srcds.MaxHostnameLength - len(hostname)

	return mpTeamname1[:remainingLen] + hostname
}

// SanitizeTeamName for safe use in SRCDS
func SanitizeTeamName(s string) string {
	s = strings.Join(strings.Fields(strings.TrimSpace(s)), "_")
	s = srcdsSafeChars.ReplaceAllString(s, "")

	if len(s) > 32 {
		return s[:28] + "..." + s[len(s)-1:]
	}

	return s
}

// validateStockMapName test if the provide map name is a valid stock map
func validateStockMapName(mapName string) error {
	if len(strings.Trim(mapName, "")) < 1 {
		return errors.New("invalid csgo map; cannot be empty or whitespace string")
	}

	if mapName != strings.ToLower(mapName) {
		return errors.New("invalid csgo map; must be all lowercase")
	}

	if strings.Index(validMaps, "/"+mapName+"/") == -1 {
		return errors.New("\"" + mapName + "\" is not a valid stock map")
	}

	return nil
}

// validateStockMapNames tests if the provide map names are all valid stock maps
func validateStockMapNames(maps []string) error {
	if len(maps) < 1 {
		return errors.New("invalid csgo map list, list was empty")
	}

	for _, m := range maps {
		if err := validateStockMapName(m); err != nil {
			return err
		}
	}

	return nil
}
