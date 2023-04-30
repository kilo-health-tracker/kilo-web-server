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

	"github.com/kilo-health-tracker/kilo-web-server/models"
	"github.com/kilo-health-tracker/kilo-web-server/utils"
)

// Formats the records returned from GetWorkoutPerformed functions into a WorkoutPerformed struct.
func formatGetWorkoutPerformedResponse(workoutPerformedRecords []models.GetWorkoutPerformedRow) WorkoutPerformed {
	workoutPerformed := WorkoutPerformed {
		Name: workoutPerformedRecords[0].WorkoutName,
		Date: workoutPerformedRecords[0].SubmittedOn.String(),
		Groups: []Group{},
	}

	groupMap := make(map[int][]ExercisePerformed)

	for _, record := range workoutPerformedRecords {
		exercisePerformed := ExercisePerformed {
			Name: record.ExerciseName,
			Weight: record.Weight,
			Reps: record.Reps,
			RIR: record.RepsInReserve.String,
		}
		
		groupMap[int(record.GroupID)] = append(groupMap[int(record.GroupID)], exercisePerformed)
		fmt.Println(groupMap[int(record.GroupID)])
		
	}

	for id, _ := range groupMap {
		group := Group {
			ID: int16(id),
			SetsPerformed: [][]ExercisePerformed{},
		}
		group.SetsPerformed = append(group.SetsPerformed, groupMap[int(id)])
		workoutPerformed.Groups = append(workoutPerformed.Groups, group)
		fmt.Println(workoutPerformed)
	}

	return workoutPerformed
}

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

	date, err := time.Parse("2006-01-02", c.Param("date"))
	if err != nil {
		return err
	}

	workoutPerformedRecords, err := queries.GetWorkoutPerformed(context.Background(), date)
	if err != nil {
		return err
	}

	workoutPerformed := formatGetWorkoutPerformedResponse(workoutPerformedRecords)

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

	date, err := time.Parse("2006-01-02", c.Param("date"))
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
	RIR    string `json:"rir"`
}

type Group struct {
	ID            int16                 `json:"id"`
	SetsPerformed [][]ExercisePerformed `json:"sets"`
}

type WorkoutPerformed struct {
	Name   string  `json:"name"`
	Date   string  `json:"date"`
	Groups []Group `json:"groups"`
}

// Submits a list of workouts tied to the given program.
func createSetsPerformed(ctx context.Context, groups []Group, queries *models.Queries, workoutID int32) error {
	for _, group := range groups {
		for index, set := range group.SetsPerformed {
			setParams := models.SubmitSetPerformedParams{
				WorkoutID: workoutID,
				GroupID:   group.ID,
				SetNumber: int16(index+1),
			}

			response, err := queries.SubmitSetPerformed(ctx, setParams)
			if err != nil {
				return err
			}

			for _, exercise := range set {
				details := models.SubmitExercisePerformedParams{
					SetID:         response.ID,
					ExerciseName:  exercise.Name,
					Reps:          exercise.Reps,
					Weight:        exercise.Weight,
					RepsInReserve: sql.NullString{String: exercise.RIR, Valid: true},
				}

				if _, err := queries.SubmitExercisePerformed(ctx, details); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func createWorkoutPerformed(requestBody WorkoutPerformed, queries *models.Queries) error {
	ctx := context.Background()

	date, err := time.Parse("2006-01-02", requestBody.Date)
	if err != nil {
		return err
	}

	workoutPerformedDetails := models.SubmitWorkoutPerformedParams{
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
func SubmitWorkoutPerformed(c echo.Context) error {
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

	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully submitted %s workout performed on %s", requestBody.Name, requestBody.Date)})

}
