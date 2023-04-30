package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kilo-health-tracker/kilo-web-server/models"
	"github.com/kilo-health-tracker/kilo-web-server/utils"
)

// Submit nutrition entry
func SubmitNutrition(c echo.Context) error {
	var requestBody models.SubmitNutritionParams
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	// bind request body to variable given
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	nutrition, err := queries.SubmitNutrition(context.Background(), requestBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to post Nutrition entry: %s", err)})
	}

	return c.JSON(http.StatusOK, nutrition)
}

// Get nutrition entry details
func GetNutrition(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	nutrition, err := queries.GetNutrition(context.Background(), c.Param("date"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nutrition)
}

// Delete nutrition entry
func DeleteNutrition(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	if err := queries.DeleteNutrition(context.Background(), c.Param("date")); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Nutrition entry submitted on: %s", c.Param("date"))})
}
