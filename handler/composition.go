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

type GenericResponse struct {
	Message string `json:"message"`
}

// Submit composition entry
func SubmitComposition(c echo.Context) error {
	var requestBody models.SubmitCompositionParams
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
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to bind request body to composition type: %s", err)})
	}

	composition, err := queries.SubmitComposition(ctx, requestBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to post composition details: %s", err)})
	}

	return c.JSON(http.StatusOK, composition)
}

// Get composition entry details
func GetComposition(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	date := c.Param("date")

	composition, err := queries.GetComposition(ctx, date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to get composition details: %s", err)})
	}

	return c.JSON(http.StatusOK, composition)
}

// Delete composition entry
func DeleteComposition(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	date := c.Param("date")

	error := queries.DeleteComposition(ctx, date)
	if error != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to delete composition: %s", error)})
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Composition submitted on: %s", date)})
}
