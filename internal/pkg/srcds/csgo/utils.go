package csgo

import (
	"math"
	"regexp"
	"strings"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

const (
	defaultMpHalftime          int = 1
	defaultMpMaxrounds         int = 30
	defaultMpMatchRestartDelay int = 15
	defaultMpOvertimeMaxrounds int = 6
	defaultSvPausable          int = 0
)

const (
	maxTeamNameLength = 31
)

var (
	srcdsSafeChars = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
)

// calculateSidesAreSwitched determines if sides should currently be swapped (mp_team1 is affiliated Terrorist)
func calculateSidesAreSwitched(mpHalftime, mpMaxRounds, mpOvertimeMaxRounds, lastCompletedRound int) bool {
	// TODO - this function needs PROPER unit tests after in-game confirmation are done
	if mpHalftime < 0 || mpHalftime > 1 {
		mpHalftime = defaultMpHalftime
	}

	if mpMaxRounds < 1 {
		mpMaxRounds = defaultMpMaxrounds
	}

	if mpOvertimeMaxRounds < 1 {
		mpOvertimeMaxRounds = defaultMpOvertimeMaxrounds
	}

	currentRound := lastCompletedRound + 1

	if mpHalftime == 1 && currentRound > mpMaxRounds/2 {
		if currentRound <= mpMaxRounds+(mpOvertimeMaxRounds/2) {
			return true
		}

		if otNotClinchable := mpOvertimeMaxRounds%2 == 0; otNotClinchable {
			//// TODO
		}
	}

	return false
}

// calculateWinThreshold determines the minimum number of rounds a team needs to win to win a map given how many rounds have been completed so far
func calculateWinThreshold(mpMaxRounds, mpOvertimeMaxRounds, lastCompletedRound int) int {
	if mpMaxRounds < 1 {
		mpMaxRounds = defaultMpMaxrounds
	}

	if mpOvertimeMaxRounds < 1 {
		mpOvertimeMaxRounds = defaultMpOvertimeMaxrounds
	}

	if notClinchable := mpMaxRounds%2 == 0; notClinchable {
		if otRoundsCompleted := lastCompletedRound - mpMaxRounds; otRoundsCompleted > 0 {
			if otNotClinchable := mpOvertimeMaxRounds%2 == 0; otNotClinchable {
				otPeriodsCompleted := otRoundsCompleted / mpOvertimeMaxRounds

				if otRoundsCompleted%mpOvertimeMaxRounds == 0 {
					otPeriodsCompleted = otPeriodsCompleted - 1
				}

				return mpMaxRounds/2 + (mpOvertimeMaxRounds / 2 * (otPeriodsCompleted + 1)) + 1
			}

			return mpMaxRounds/2 + mpOvertimeMaxRounds/2 + 1
		}
	}

	return mpMaxRounds/2 + 1
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

	if len(s) > maxTeamNameLength {
		return s[:maxTeamNameLength-3] + "..." + s[len(s)-1:]
	}

	return s
}
