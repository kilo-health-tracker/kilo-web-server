package handler

import (
	"context"
	"fmt"
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
		return err
	}

	name := c.Param("name")
	workout, err := queries.GetWorkout(context.Background(), name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, workout)
}

// Get workout performed
func GetWorkoutPerformed(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	date, _ := time.Parse("YYYY-MM-DD", c.Param("date"))
	if err != nil {
		return err
	}

	workoutPerformed, err := queries.GetWorkoutPerformed(context.Background(), date)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, workoutPerformed)
}

// Get workout names
func GetWorkoutNames(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return err
	}

	workoutNames, err := queries.GetWorkoutNames(context.Background(), int32(limit))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, workoutNames)
}

// Delete workout
func DeleteWorkout(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	if err := queries.DeleteWorkout(context.Background(), c.Param("name")); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Workout: %s", c.Param("name"))})
}

// Delete workout performed
func DeleteWorkoutPerformed(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	date, err := time.Parse("YYYY-MM-DD", c.Param("date"))
	if err != nil {
		return err
	}

	if err := queries.DeleteWorkoutPerformed(context.Background(), date); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Workout performed on: %s", date)})
}
