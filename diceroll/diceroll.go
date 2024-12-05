package diceroll

// Handle the logic for Natural One
func HandleNaturalOne(
	naturalOne, criticalFailureThreshold, failureThreshold, successThreshold, maxDiceValue int,
	criticalFailures, failures, successes, criticalSuccesses *int,
) {
	if naturalOne > criticalFailureThreshold {
		switch {
		case naturalOne >= successThreshold:
			if *criticalSuccesses == maxDiceValue {
				*criticalSuccesses--
			}
			*successes++
		case naturalOne >= failureThreshold:
			*successes = max(0, *successes-1)
			*failures++
		default:
			*failures = max(0, *failures-1)
			*criticalFailures++
		}
	}
}

func HandleNaturalTwenty(
	naturalTwenty, successThreshold, failureThreshold, maxDiceValue int,
	criticalFailures, failures, successes, criticalSuccesses *int,
) {
	if naturalTwenty < successThreshold {
		switch {
		case naturalTwenty < failureThreshold:
			*failures = max(0, *failures-1)
			// For nat20 it always has at least one failure
			if *criticalFailures == maxDiceValue {
				*criticalFailures--
				*failures++
			} else {
				*successes++
			}
		case naturalTwenty < successThreshold:
			*successes = max(0, *successes-1)
			*criticalSuccesses++
		}
	} else {
		// Natural Twenty is a guaranteed success, so it's promoted to critical success.
		// Impossible to have 20 critical successes
		*criticalSuccesses = min(19, *criticalSuccesses+1)
	}
}

func DiceRollOdds(modifier, dc int) (criticalFailures, failures, successes, criticalSuccesses int) {

	var minDiceValue int = 1
	var maxDiceValue int = 20

	criticalFailureThreshold := dc - 10
	failureThreshold := dc
	successThreshold := dc + 10

	// Count range of occurrences for each DiceValue;
	// modifier -1 garantees that superior thresholds are exclusive when calculating fails and successes
	criticalFailures = max(0, min(maxDiceValue, criticalFailureThreshold-modifier))
	failures = max(0, min(maxDiceValue, failureThreshold-modifier-1)) - criticalFailures
	successes = max(0, min(maxDiceValue, successThreshold-modifier-1)) - failures - criticalFailures
	criticalSuccesses = max(0, min(20, maxDiceValue-(successThreshold-modifier)))

	// Check that success and criticalSuccesses are 0 when impossible
	if modifier+maxDiceValue < failureThreshold {
		//When highest value(20+modifier) < dc, we cant have success
		successes = 0
		criticalSuccesses = 0
	}
	// only critical failures are possible
	if modifier+maxDiceValue < criticalFailureThreshold {
		failures = 0
	}

	naturalOne := minDiceValue + modifier
	naturalTwenty := maxDiceValue + modifier
	// Handle Natural One and Natural Twenty
	HandleNaturalOne(naturalOne, criticalFailureThreshold, failureThreshold, successThreshold, maxDiceValue,
		&criticalFailures, &failures, &successes, &criticalSuccesses)

	HandleNaturalTwenty(naturalTwenty, successThreshold, failureThreshold, maxDiceValue,
		&criticalFailures, &failures, &successes, &criticalSuccesses)

	return criticalFailures, failures, successes, criticalSuccesses
}
