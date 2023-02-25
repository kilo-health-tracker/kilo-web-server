package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kilo-health-tracker/kilo-database/utils"
)

// Get workout
func GetWorkout(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	name := c.Param("name")

	workout, err := queries.GetWorkout(ctx, name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to get Workout: %s", err)})
	}

	return c.JSON(http.StatusOK, workout)
}

// Get workout performed
func GetWorkoutPerformed(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	dateString := c.Param("date")
	date, err := time.Parse("YYYY-MM-DD", dateString)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to parse string to date: %s", err)})
	}

	workoutPerformed, err := queries.GetWorkoutPerformed(ctx, date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to get Workout performed: %s", err)})
	}
	return c.JSON(http.StatusOK, workoutPerformed)
}

// Get workout names
func GetWorkoutNames(c echo.Context) error {
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

	workoutNames, err := queries.GetWorkoutNames(ctx, int32(limit))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to get Workout names: %s", err)})
	}

	return c.JSON(http.StatusOK, workoutNames)
}

// Delete workout
func DeleteWorkout(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	name := c.Param("name")

	error := queries.DeleteWorkout(ctx, name)
	if error != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to delete Workout: %s", error)})
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Workout: %s", name)})
}

// Delete workout
func DeleteWorkoutPerformed(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	dateString := c.Param("date")
	date, err := time.Parse("YYYY-MM-DD", dateString)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to parse string to date: %s", err)})
	}

	error := queries.DeleteWorkoutPerformed(ctx, date)
	if error != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to delete Workout performed: %s", error)})
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Workout performed on: %s", date)})
}
