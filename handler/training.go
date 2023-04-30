package handler

import (
	"net/http"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kilo-health-tracker/kilo-web-server/utils"
)

type TrainingWeight struct {
	Percentage float64 `json:"percentage"`
	TrainingWeight float64 `json:"trainingWeight"`
	TrainingMax float64 `json:"trainingMax"`
}

// Calculate training weight given reps, rir, and weight
func GetTrainingWeight(c echo.Context) error {
	reps, err := strconv.Atoi(c.QueryParam("reps"))
	if err != nil {
		return err
	}
	rir, err := strconv.ParseFloat(c.QueryParam("rir"), 64)
	if err != nil {
		return err
	}
	weight, err := strconv.ParseFloat(c.QueryParam("weight"), 64)
	if err != nil {
		return err
	}

	rirMapping := utils.GetRIRMapping()
	weightTable := utils.GetWeightTable()
	percentage := weightTable[reps-1][rirMapping[rir]]
	trainingWeight := fmt.Sprintf("%.2f", weight*percentage)
	trainingMax := fmt.Sprintf("%.2f", weight/percentage)

	fmt.Printf("Percentage: %v", percentage*100)
	fmt.Printf("\nTraining Weight: %v", trainingWeight)
	fmt.Printf("\nTraining Max: %v", trainingMax)

	response := TrainingWeight {
		percentage,
		weight*percentage,
		weight/percentage,
	}

	return c.JSON(http.StatusOK, response)
}