package diceroll_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/felipezoltowski/go-webserver/diceroll"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Estrutura de resposta esperada
type Response struct {
	CriticalFailures  int `json:"critical_failures"`
	Failures          int `json:"failures"`
	Successes         int `json:"successes"`
	CriticalSuccesses int `json:"critical_successes"`
}

func TestDistributionEndpoint(t *testing.T) {
	// Router setup
	router := gin.Default()
	router.GET("/api/pathfinder2e/v1/distribution", func(c *gin.Context) {
		modifier := c.Query("modifier")
		dc := c.Query("dc")

		modifierInt, err := strconv.Atoi(modifier)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "invalid modifier",
			})
			return
		}

		dcInt, err := strconv.Atoi(dc)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "invalid dc",
			})
			return
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

	t.Run("Valid Input", func(t *testing.T) {
		// Simular requisição HTTP
		req, _ := http.NewRequest("GET", "/api/pathfinder2e/v1/distribution?modifier=15&dc=20", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificar resposta
		assert.Equal(t, http.StatusOK, w.Code)

		var resp Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// Verificar os valores calculados (valores esperados do exemplo fornecido)
		assert.Equal(t, 1, resp.CriticalFailures)
		assert.Equal(t, 3, resp.Failures)
		assert.Equal(t, 10, resp.Successes)
		assert.Equal(t, 6, resp.CriticalSuccesses)
	})

	t.Run("Valid Input v2", func(t *testing.T) {
		// Simular requisição HTTP
		req, _ := http.NewRequest("GET", "/api/pathfinder2e/v1/distribution?modifier=1&dc=50", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificar resposta
		assert.Equal(t, http.StatusOK, w.Code)

		var resp Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// Verificar os valores calculados (valores esperados do exemplo fornecido)
		assert.Equal(t, 19, resp.CriticalFailures)
		assert.Equal(t, 1, resp.Failures)
		assert.Equal(t, 0, resp.Successes)
		assert.Equal(t, 0, resp.CriticalSuccesses)
	})

	t.Run("Valid Input v3", func(t *testing.T) {
		// Simular requisição HTTP
		req, _ := http.NewRequest("GET", "/api/pathfinder2e/v1/distribution?modifier=50&dc=1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificar resposta
		assert.Equal(t, http.StatusOK, w.Code)

		var resp Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// Verificar os valores calculados (valores esperados do exemplo fornecido)
		assert.Equal(t, 0, resp.CriticalFailures)
		assert.Equal(t, 0, resp.Failures)
		assert.Equal(t, 1, resp.Successes)
		assert.Equal(t, 19, resp.CriticalSuccesses)
	})

	t.Run("Invalid Modifier", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/pathfinder2e/v1/distribution?modifier=abc&dc=20", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificar resposta
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "invalid modifier", resp["erro"])
	})

	t.Run("Invalid DC", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/pathfinder2e/v1/distribution?modifier=15&dc=xyz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verificar resposta
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "invalid dc", resp["erro"])
	})
}
