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

type GenericResponse struct {
	Message string `json:"message"`
}

// Submit composition entry
func SubmitComposition(c echo.Context) error {
	var requestBody models.SubmitCompositionParams
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	// bind request body to variable given
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	composition, err := queries.SubmitComposition(context.Background(), requestBody)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, composition)
}

// Get composition entry details
func GetComposition(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	composition, err := queries.GetComposition(context.Background(), c.Param("date"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, composition)
}

// Delete composition entry
func DeleteComposition(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	if err := queries.DeleteComposition(context.Background(), c.Param("date")); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Composition submitted on: %s", c.Param("date"))})
}
