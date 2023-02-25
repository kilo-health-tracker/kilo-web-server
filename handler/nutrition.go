package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kilo-health-tracker/kilo-database/models"
	"github.com/kilo-health-tracker/kilo-database/utils"
)

// Submit nutrition entry
func SubmitNutrition(c echo.Context) error {
	var requestBody models.SubmitNutritionParams
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()

	// bind request body to variable given
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to bind request body to Nutrition type: %s", err)})
	}

	nutrition, err := queries.SubmitNutrition(ctx, requestBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to post Nutrition entry: %s", err)})
	}

	return c.JSON(http.StatusOK, nutrition)
}

// Get nutrition entry details
func GetNutrition(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	date := c.Param("date")

	nutrition, err := queries.GetNutrition(ctx, date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to get Nutrition entry: %s", err)})
	}

	return c.JSON(http.StatusOK, nutrition)
}

// Delete nutrition entry
func DeleteNutrition(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	date := c.Param("date")

	error := queries.DeleteNutrition(ctx, date)
	if error != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to delete Nutrition entry: %s", error)})
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Nutrition entry submitted on: %s", date)})
}
