package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"database/sql"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kilo-health-tracker/kilo-database/utils"
	"github.com/kilo-health-tracker/kilo-database/models"
)

type GetProgramRow struct {
	Name         string        `json:"name"`
	Name_2       string        `json:"name2"`
	GroupID      int16         `json:"groupID"`
	ExerciseName string        `json:"exerciseName"`
	Weight       sql.NullInt16 `json:"weight"`
	Sets         int16         `json:"sets"`
	Reps         int16         `json:"reps"`
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
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	name := c.Param("name")

	programRecords, err := queries.GetProgram(ctx, name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to get Program: %s", err)})
	}
	log.Println(programRecords)

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

	return c.JSON(http.StatusOK, program)
}

// Get Program names
func GetProgramNames(c echo.Context) error {
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

	composition, err := queries.GetProgramNames(ctx, int32(limit))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to get Program names: %s", err)})
	}

	return c.JSON(http.StatusOK, composition)
}

// Delete Program
func DeleteProgram(c echo.Context) error {
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	ctx := context.Background()
	name := c.Param("name")

	error := queries.DeleteProgram(ctx, name)
	if error != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to delete Program: %s", error)})
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted Program: %s", name)})
}

func CreateProgram(c echo.Context) error {
	var requestBody Program
	queries, err := utils.GetQueryInterface()
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to establish connection to postgres: %s", err)})
	}

	// bind request body to variable given
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to bind request body to composition type: %s", err)})
	}

	ctx := context.Background()

	log.Println(requestBody)
	programName := requestBody.Name
	fmt.Printf("Inserting program: %s\n", programName)

	response, err := queries.SubmitProgram(ctx, programName)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to submit program: %s", err)})
	}
	log.Println(response)

	for _, workout := range requestBody.Workouts {
		workoutParams := models.SubmitWorkoutParams{
			Name:        workout.Name,
			ProgramName: programName,
		}
		log.Println(workoutParams)

		log.Println("Submitting Workout...")
		_, err := queries.SubmitWorkout(ctx, workoutParams)
		if err != nil {
			log.Fatal(err)
			return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to submit program: %s", err)})
		}

		programWorkoutLink := models.SubmitProgramDetailsParams{
			ProgramName: programName,
			WorkoutName: workout.Name,
		}

		log.Println("Submitting program details...")
		_, err2 := queries.SubmitProgramDetails(ctx, programWorkoutLink)
		if err2 != nil {
			log.Fatal(err2)
			return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to submit program: %s", err2)})
		}

		log.Println("Submitting workout details...")
		for _, exercise := range workout.Exercises {
			details := models.SubmitWorkoutDetailsParams {
				WorkoutName:  workout.Name,
				GroupID:      exercise.GroupId,
				ExerciseName: exercise.Name,
				Sets:         exercise.Sets,
				Reps:         exercise.Reps,
				Weight:       sql.NullInt16{Int16: exercise.Weight, Valid: true},
			}
			log.Println(details)
			_, err := queries.SubmitWorkoutDetails(ctx, details)
			if err != nil {
				log.Fatal(err)
				return c.JSON(http.StatusInternalServerError, GenericResponse{fmt.Sprintf("Failed to submit program: %s", err)})
			}
		}
	}

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully created Program: %s", programName)})
}