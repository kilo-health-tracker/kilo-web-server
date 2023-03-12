package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"database/sql"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kilo-health-tracker/kilo-database/utils"
	"github.com/kilo-health-tracker/kilo-database/models"
)

// Formats the records returned from GetProgram functions into a Program struct.
func formatProgramResponse(programRecords []models.GetProgramRow) Program {
	program := Program {
		Name: programRecords[0].Name,
		Workouts: []Workout {},
	}

	for _, record := range programRecords {
		workout := Workout {
			Name: record.Name_2,
			Exercises: []Exercise {
				{
					Name: record.ExerciseName,
					GroupId: record.GroupID,
					Sets: record.Sets,
					Reps: record.Reps,
					Weight: record.Weight.Int16,
				},
			},
		}

		program.Workouts = append(program.Workouts, workout)
	}

	return program
}

type Exercise struct {
	Name string `json:"name"`
	GroupId int16 `json:"group_id"`
	Sets int16 `json:"sets"`
	Reps int16 `json:"reps"`
	Weight int16 `json:"weight"`
}

type Workout struct {
	Name string `json:"name"`
	Exercises []Exercise `json:"exercises"`
}

type Program struct {
	Name string `json:"name"`
	Workouts []Workout `json:"workouts"`
}

// Get Program
func GetProgram(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	ctx := context.Background()
	name := c.Param("name")

	programRecords, err := queries.GetProgram(ctx, name)
	if err != nil {
		return err
	}
	program := formatProgramResponse(programRecords)

	return c.JSON(http.StatusOK, program)
}

// Gets a list of program names.
func GetProgramNames(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return err
	}

	ctx := context.Background()
	composition, err := queries.GetProgramNames(ctx, int32(limit))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, composition)
}

// Delete Program
func DeleteProgram(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	ctx := context.Background()
	name := c.Param("name")

	if err := queries.DeleteProgram(ctx, name); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Program: %s", name)})
}

// Submits a list of exercises tied to a given workout.
func submitExercises(workoutName string, exercises []Exercise, ctx context.Context, queries *models.Queries) error {
	for _, exercise := range exercises {
		details := models.SubmitWorkoutDetailsParams {
			WorkoutName:  workoutName,
			GroupID:      exercise.GroupId,
			ExerciseName: exercise.Name,
			Sets:         exercise.Sets,
			Reps:         exercise.Reps,
			Weight:       sql.NullInt16{Int16: exercise.Weight, Valid: true},
		}

		if _, err := queries.SubmitWorkoutDetails(ctx, details); err != nil {
			return err
		}
	}

	return nil
}

// Submits a list of workouts tied to the given program.
func submitWorkouts(ctx context.Context, requestBody Program, queries *models.Queries) error {
	for _, workout := range requestBody.Workouts {
		workoutParams := models.SubmitWorkoutParams{
			Name:        workout.Name,
			ProgramName: requestBody.Name,
		}

		if _, err := queries.SubmitWorkout(ctx, workoutParams); err != nil {
			return err
		}

		programWorkoutLink := models.SubmitProgramDetailsParams{
			ProgramName: requestBody.Name,
			WorkoutName: workout.Name,
		}

		if _, err := queries.SubmitProgramDetails(ctx, programWorkoutLink); err != nil {
			return err
		}

		if err := submitExercises(workout.Name, workout.Exercises, ctx, queries); err != nil {
			return err
		}
	}

	return nil
}

// Submits a program.
func submitProgram(requestBody Program, queries *models.Queries) error {
	ctx := context.Background()

	if _, err := queries.SubmitProgram(ctx, requestBody.Name); err != nil {
		return err
	}

	if err := submitWorkouts(ctx, requestBody, queries); err != nil {
		return err
	}

	return nil
}

func CreateProgram(c echo.Context) error {
	var requestBody Program

	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	// bind request body to variable given
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	if err := submitProgram(requestBody, queries); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully created Program: %s", requestBody.Name)})
}