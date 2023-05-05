package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kilo-health-tracker/kilo-web-server/utils"
	"github.com/kilo-health-tracker/kilo-web-server/models"
)

// Get exercise details
func GetExercise(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	exercise, err := queries.GetExercise(context.Background(), c.Param("name"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, exercise)
}

// Get exercise names
func GetExerciseNames(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return err
	}

	exerciseNames, err := queries.GetExercises(context.Background(), int32(limit))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, exerciseNames)
}

type ExerciseInput struct {
	Name      string   `json:"name"`
	Primary   []string `json:"primary"`
	Secondary []string `json:"secondary"`
	Tertiary  []string `json:"tertiary"`
	Type      string   `json:"type"`
	Variation string   `json:"variation"`
}

// Submits an exercise definition
func SubmitExercise(c echo.Context) error {
	var requestBody ExerciseInput
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	// bind request body to variable given
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	ctx := context.Background()

	record := models.SubmitExerciseParams{
		Name:      requestBody.Name,
		Type:      sql.NullString{String: requestBody.Type, Valid: true},
		Variation: sql.NullString{String: requestBody.Variation, Valid: true},
	}

	exercise, err := queries.SubmitExercise(ctx, record)
	if err != nil {
		return err
	}

	for _, bodyPart := range requestBody.Primary {
		entry := models.SubmitExerciseDetailsParams{
			ExerciseName: requestBody.Name,
			BodyPart:     bodyPart,
			Level:        "primary",
		}
		_, err := queries.SubmitExerciseDetails(ctx, entry)
		if err != nil {
			return err
		}
	}
	for _, bodyPart := range requestBody.Secondary {
		entry := models.SubmitExerciseDetailsParams{
			ExerciseName: requestBody.Name,
			BodyPart:     bodyPart,
			Level:        "secondary",
		}
		_, err := queries.SubmitExerciseDetails(ctx, entry)
		if err != nil {
			return err
		}
	}
	for _, bodyPart := range requestBody.Tertiary {
		entry := models.SubmitExerciseDetailsParams{
			ExerciseName: requestBody.Name,
			BodyPart:     bodyPart,
			Level:        "tertiary",
		}
		_, err := queries.SubmitExerciseDetails(ctx, entry)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, exercise)
}

// Delete program
func DeleteExercise(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	if err := queries.DeleteExercise(context.Background(), c.Param("name")); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted the body part definition for: %s", c.Param("name"))})
}
