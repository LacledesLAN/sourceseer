package csgo

//
// Counter-Strike: Global Offensive's logging doesn't include much information about actions taken from game rules.
//
//	• If/when team sides are switched at half-time
//	• If/when overtime begins
//	• If/when the current match is over (and the scoreboard is being displayed)
//
// Since we need this kind of information is needed to make time-sensitive reaction-decisions runtime calculations based
// on assumptions must be made. All of these functions need to be rock-solid and should make bests guesses whenever required
// to keep running. Since the current game state cannot be guaranteed calculations are usually done using past-tense
// information that is known.
//

const (
	defaultMpDoWarmupPeriod    int = 1
	defaultMpHalftime          int = 1
	defaultMpMaxrounds         int = 30
	defaultMpMatchRestartDelay int = 15
	defaultMpOvertimeEnabled   int = 1
	defaultMpOvertimeMaxrounds int = 6
	defaultMpWarmupPausetimer  int = 0
	defaultMpWarmuptime        int = 300
	defaultSvPausable          int = 0
)

func calculateLastRoundWinThreshold(mpMaxRounds, mpOvertimeMaxRounds int, lastCompletedRound lastInt) lastInt {
	if mpMaxRounds < 1 {
		mpMaxRounds = defaultMpMaxrounds
	}

	if notClinchable := mpMaxRounds%2 == 0; notClinchable {
		if otRoundsCompleted := int(lastCompletedRound) - mpMaxRounds; otRoundsCompleted > 0 {
			if mpOvertimeMaxRounds < 1 {
				mpOvertimeMaxRounds = defaultMpOvertimeMaxrounds
			}

			if otNotClinchable := mpOvertimeMaxRounds%2 == 0; otNotClinchable {
				otPeriodsCompleted := calcOvertimePeriodNumber(mpMaxRounds, mpOvertimeMaxRounds, lastCompletedRound) - 1

				if otRoundsCompleted%mpOvertimeMaxRounds == 0 {
					otPeriodsCompleted--
				}

				return lastInt(mpMaxRounds/2 + (mpOvertimeMaxRounds / 2 * (otPeriodsCompleted + 1)) + 1)
			}

			return lastInt(mpMaxRounds/2 + mpOvertimeMaxRounds/2 + 1)
		}
	}

	return lastInt(mpMaxRounds/2 + 1)
}

func calcOvertimePeriodNumber(mpMaxRounds, mpOvertimeMaxRounds int, lastCompletedRound lastInt) int {
	if mpMaxRounds < 1 {
		mpMaxRounds = defaultMpMaxrounds
	}

	if int(lastCompletedRound)-mpMaxRounds < 0 {
		return 0
	}

	if mpOvertimeMaxRounds < 1 {
		mpOvertimeMaxRounds = defaultMpOvertimeMaxrounds
	}

	if otNotClinchable := mpOvertimeMaxRounds%2 == 0; otNotClinchable {
		return ((int(lastCompletedRound) - mpMaxRounds) / mpOvertimeMaxRounds) + 1
	}

	return 1
}

func calculateSidesAreCurrentlySwitched(mpHalftime, mpMaxRounds, mpOvertimeMaxRounds int, lastCompletedRound lastInt) bool {
	if mpHalftime < 0 || mpHalftime > 1 {
		mpHalftime = defaultMpHalftime
	}

	if mpMaxRounds < 1 {
		mpMaxRounds = defaultMpMaxrounds
	}

	currentRound := int(lastCompletedRound) + 1

	if currentRound > mpMaxRounds/2 && mpHalftime == 1 {
		if mpOvertimeMaxRounds < 1 {
			mpOvertimeMaxRounds = defaultMpOvertimeMaxrounds
		}

		if otPeriod := calcOvertimePeriodNumber(mpMaxRounds, mpOvertimeMaxRounds, lastCompletedRound) - 1; otPeriod > 0 {
			if otNotClinchable := mpOvertimeMaxRounds%2 == 0; otNotClinchable {
				otRoundsCompleted := currentRound - mpMaxRounds

				if isFirstHalf := otRoundsCompleted-(otPeriod*mpOvertimeMaxRounds) <= mpOvertimeMaxRounds/2; isFirstHalf {
					return otPeriod%2 == 0
				}

				return otPeriod%2 != 0
			}
		}

		return currentRound <= mpMaxRounds+(mpOvertimeMaxRounds/2)
	}

	return false
}
