package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kilo-health-tracker/kilo-web-server/utils"
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

// Delete program
func DeleteExercise(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	if err := queries.DeleteExercise(context.Background(), c.Param("name")); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Exercise: %s", c.Param("name"))})
}
