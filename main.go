package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Define struct to format json return as requested
type Response struct {
	CriticalFailures  int `json:"critical_failures"`
	Failures          int `json:"failures"`
	Successes         int `json:"successes"`
	CriticalSuccesses int `json:"critical_successes"`
}

// Handle the logic for Natural One
func HandleNaturalOne(
	naturalOne, criticalFailureThreshold, failureThreshold, successThreshold, maxDiceValue int,
	criticalFailures, failures, successes, criticalSuccesses *int,
) {
	if naturalOne > criticalFailureThreshold {
		switch {
		case naturalOne >= successThreshold:
			if *criticalSuccesses == maxDiceValue {
				println("entrou")
				*criticalSuccesses--
			} else {
				println("entrou2")
				*criticalSuccesses = max(0, *criticalSuccesses-1)
			}
			*successes++
		case naturalOne >= failureThreshold:
			println("entrou3")
			*successes = max(0, *successes-1)
			*failures++
		default:
			println("entrou4")
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
				println("entrou5")
				*criticalFailures--
				*failures++
			} else {
				*successes++
			}
		case naturalTwenty < successThreshold:
			*successes = max(0, *successes-1)
			println("entrou6")
			*criticalSuccesses++
		}
	} else {
		// Natural Twenty is a guaranteed success, so it's promoted to critical success.
		// Impossible to have 20 critical successes
		if *criticalSuccesses == maxDiceValue {
			println("entrou7")
			*criticalSuccesses--
		} else {
			println("entrou8")
			*criticalSuccesses = min(19, *criticalSuccesses+1)
		}
	}
}

func DiceRollCalc(modifier, dc int) (criticalFailures, failures, successes, criticalSuccesses int) {

	var minDiceValue int = 1
	var maxDiceValue int = 20

	criticalFailureThreshold := dc - 10
	failureThreshold := dc
	successThreshold := dc + 10

	// Count range of occurrences for each DiceValue;
	// modifier -1 is to make an inclusive for the range.
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
	fmt.Printf("naturalOne=%d \n", naturalOne)
	HandleNaturalOne(naturalOne, criticalFailureThreshold, failureThreshold, successThreshold, maxDiceValue,
		&criticalFailures, &failures, &successes, &criticalSuccesses)

	fmt.Printf("naturalTwenty=%d \n", naturalTwenty)
	HandleNaturalTwenty(naturalTwenty, successThreshold, failureThreshold, maxDiceValue,
		&criticalFailures, &failures, &successes, &criticalSuccesses)

	return criticalFailures, failures, successes, criticalSuccesses
}

func main() {
	app := gin.Default()

	app.Use(cors.Default())

	// Define um endpoint
	app.GET("/api/pathfinder2e/v1/distribution", func(c *gin.Context) {

		modifier := c.Query("modifier")
		dc := c.Query("dc")

		modifierInt, err := strconv.Atoi(modifier)
		if err != nil {
			c.JSON(400, gin.H{
				"erro": "invalid modifier",
			})
		}

		dcInt, err := strconv.Atoi(dc)
		if err != nil {
			c.JSON(400, gin.H{
				"erro": "invalid modifier",
			})
		}

		criticalFailures, failures, successes, criticalSuccesses := DiceRollCalc(modifierInt, dcInt)

		response := Response{
			CriticalFailures:  criticalFailures,
			Failures:          failures,
			Successes:         successes,
			CriticalSuccesses: criticalSuccesses,
		}

		c.JSON(http.StatusOK, response)

	})

	// Start server on port 8080
	app.Run(":8080")
}
