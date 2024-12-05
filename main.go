package main

import (
	"net/http"
	"strconv"

	"github.com/felipezoltowski/go-webserver/diceroll"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Struct to format json return
type Response struct {
	CriticalFailures  int `json:"critical_failures"`
	Failures          int `json:"failures"`
	Successes         int `json:"successes"`
	CriticalSuccesses int `json:"critical_successes"`
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

		criticalFailures, failures, successes, criticalSuccesses := diceroll.DiceRollOdds(modifierInt, dcInt)

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
