package diceroll

const (
	minDiceValue int = 1
	maxDiceValue int = 20
)

func calcDegreesThreshold(dc, offset int) int {
	return dc + offset
}

func calcOddsRange(maxDiceValue, threshold, modifier, offset int) int {
	return max(0, min(maxDiceValue, threshold-modifier+offset))
}

func HandleNaturalOne(
	naturalOne, criticalFailureThreshold, failureThreshold, successThreshold, maxDiceValue int,
	criticalFailures, failures, successes, criticalSuccesses int,
) (int, int, int, int) {
	if naturalOne > criticalFailureThreshold {
		switch {
		case naturalOne >= successThreshold:
			if criticalSuccesses == maxDiceValue {
				criticalSuccesses--
			}
			successes++
		case naturalOne >= failureThreshold:
			successes = max(0, successes-1)
			failures++
		default:
			failures = max(0, failures-1)
			criticalFailures++
		}
	}
	return criticalFailures, failures, successes, criticalSuccesses
}

func HandleNaturalTwenty(
	naturalTwenty, successThreshold, failureThreshold, maxDiceValue int,
	criticalFailures, failures, successes, criticalSuccesses int,
) (int, int, int, int) {
	if naturalTwenty < successThreshold {
		switch {
		case naturalTwenty < failureThreshold:
			failures = max(0, failures-1)
			// For nat20 it always has at least one failure
			if criticalFailures == maxDiceValue {
				criticalFailures--
				failures++
			} else {
				successes++
			}
		case naturalTwenty < successThreshold:
			successes = max(0, successes-1)
			criticalSuccesses++
		}
	} else {
		// Natural Twenty is a guaranteed success, so it's promoted to critical success.
		// Impossible to have 20 critical successes
		criticalSuccesses = min(19, criticalSuccesses+1)
	}
	return criticalFailures, failures, successes, criticalSuccesses
}

func DiceRollOdds(modifier, dc int) (criticalFailures, failures, successes, criticalSuccesses int) {

	criticalFailureThreshold := calcDegreesThreshold(dc, -10)
	failureThreshold := calcDegreesThreshold(dc, 0)
	successThreshold := calcDegreesThreshold(dc, 10)

	// modifier -1 garantees that superior thresholds are exclusive when calculating fails and successes
	criticalFailures = calcOddsRange(maxDiceValue, criticalFailureThreshold, modifier, 0)
	failures = calcOddsRange(maxDiceValue, failureThreshold, modifier, -1) - criticalFailures
	successes = calcOddsRange(maxDiceValue, successThreshold, modifier, -1) - failures - criticalFailures
	criticalSuccesses = max(0, min(20, maxDiceValue-(successThreshold-modifier)))

	// Check that success and criticalSuccesses are 0 when impossible
	if modifier+maxDiceValue < failureThreshold {
		//When highest dice value(20+modifier) < dc, we cant have success
		successes = 0
		criticalSuccesses = 0
	}

	naturalOne := minDiceValue + modifier
	naturalTwenty := maxDiceValue + modifier

	criticalFailures, failures, successes, criticalSuccesses = HandleNaturalOne(naturalOne, criticalFailureThreshold, failureThreshold, successThreshold, maxDiceValue,
		criticalFailures, failures, successes, criticalSuccesses)

	criticalFailures, failures, successes, criticalSuccesses = HandleNaturalTwenty(naturalTwenty, successThreshold, failureThreshold, maxDiceValue,
		criticalFailures, failures, successes, criticalSuccesses)

	return criticalFailures, failures, successes, criticalSuccesses
}
