package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kilo-health-tracker/kilo-database/models"
	"github.com/kilo-health-tracker/kilo-database/utils"
)

// Get workout
func GetWorkout(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	workout, err := queries.GetWorkout(context.Background(), c.Param("name"))
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

	date, err := time.Parse("YYYY-MM-DD", c.Param("date"))
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



type ExercisePerformed struct {
	Name   string `json:"name"`
	Weight int16  `json:"weight"`
	Reps   int16  `json:"reps"`
	RIR    int16  `json:"rir"`
}

type SetPerformed struct {
	ExercisesPerformed []ExercisePerformed
}

type Group struct {
	ID                 int16               `json:"id"`
	SetsPerformed []SetPerformed `json:"sets"`
}

type WorkoutPerformed struct {
	Name   string  `json:"name"`
	Date   string  `json:"date"`
	Groups []Group `json:"groups"`
}

// Submits a list of exercises tied to a given workout.
func createExercisesPerformed(setID int32, sets []SetPerformed, ctx context.Context, queries *models.Queries) error {
	for index, set := range sets {
		for _, exercise := range set.ExercisesPerformed {
			details := models.SubmitExercisePerformedParams {
				SetID: setID,
				ExerciseName: exercise.Name,
				Reps: exercise.Reps,
				Weight: exercise.Weight,
				RepsInReserve: sql.NullInt16{Int16: exercise.Weight, Valid: true},
			}
	
			if _, err := queries.SubmitExercisePerformed(ctx, details); err != nil {
				return err
			}
		}
	}

	return nil
}


// Submits a list of workouts tied to the given program.
func createSetsPerformed(ctx context.Context, groups []Group, queries *models.Queries, workoutID int32) error {
	for _, group := range groups {
		setParams := models.SubmitSetPerformedParams{
			WorkoutID:	workoutID,
			GroupID: 	group.ID,
		}

		response, err := queries.SubmitSetPerformed(ctx, setParams)
		if err != nil {
			return err
		}

		if err := createExercisesPerformed(response.ID, group.SetsPerformed, ctx, queries); err != nil {
			return err
		}
	}

	return nil
}

func createWorkoutPerformed(requestBody WorkoutPerformed, queries *models.Queries) error {
	ctx := context.Background()

	date, err := time.Parse("YYYY-MM-DD", requestBody.Date)
	if err != nil {
		return err
	}

	workoutPerformedDetails := models.SubmitWorkoutPerformedParams {
		SubmittedOn: date,
		WorkoutName: requestBody.Name,
	}
	response, err := queries.SubmitWorkoutPerformed(ctx, workoutPerformedDetails)
	if err != nil {
		return err
	}

	if err := createSetsPerformed(ctx, requestBody.Groups, queries, response.ID); err != nil {
		return err
	}

	return nil
}

// Submit workout performed
func SubmitWorkoutPerfomed(c echo.Context) error {
	var requestBody WorkoutPerformed
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if err := createWorkoutPerformed(requestBody, queries); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully created Program: %s", requestBody.Name)})

}
