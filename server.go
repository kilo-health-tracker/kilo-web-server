package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/kilo-health-tracker/kilo-web-server/handler"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	//// Exercise
	e.GET("/exercise", handler.GetExerciseNames)
	e.GET("/exercise/:name", handler.GetExercise)
	e.DELETE("/exercise/:name", handler.DeleteExercise)

	//// Workout
	e.GET("/workout", handler.GetWorkoutNames)
	e.DELETE("/workout/:name", handler.DeleteWorkout)
	e.GET("/workout/:name", handler.GetWorkout)
	e.GET("/workout/:name/:date", handler.GetWorkoutPerformed)
	e.DELETE("/workout/:name/:date", handler.DeleteWorkoutPerformed)

	//// Program
	e.GET("/exercise", handler.GetProgramNames)
	e.GET("/exercise/:name", handler.GetProgram)

	//// Composition
	e.POST("/composition", handler.SubmitComposition)
	e.GET("/composition/:date", handler.GetComposition)
	e.DELETE("/composition/:date ", handler.DeleteComposition)

	//// Nutrition
	e.POST("/nutrition", handler.SubmitNutrition)
	e.GET("/nutrition/:date", handler.GetNutrition)
	e.DELETE("/nutrition/:date", handler.DeleteNutrition)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
