package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kilo-health-tracker/kilo-database/utils"
)

// Get exercise details
func GetExercise(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	name := c.Param("name")

	exercise, err := queries.GetExercise(ctx, name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to get Exercise: %s", err)})
	}

	return c.JSON(http.StatusOK, exercise)
}

// Get exercise names
func GetExerciseNames(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	limitString := c.QueryParam("limit")

	limit, err := strconv.Atoi(limitString)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to convert limit to integer: %s", err)})
	}

	exerciseNames, err := queries.GetExercises(ctx, int32(limit))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to get Exercise names: %s", err)})
	}

	return c.JSON(http.StatusOK, exerciseNames)
}

// Delete program
func DeleteExercise(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	name := c.Param("name")

	error := queries.DeleteExercise(ctx, name)
	if error != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to delete Exercise: %s", error)})
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Exercise: %s", name)})
}
