package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/kilo-health-tracker/kilo-web-server/handler"
)

func main() {
	// Echo instance
	server := echo.New()

	// Middleware
	server.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.RequestID(),
	)

	server.HTTPErrorHandler = func(err error, c echo.Context) {
		// Take required information from error and context and send it to a service like New Relic
		fmt.Println(c.Path(), c.QueryParams(), err.Error())

		// Call the default handler to return the HTTP response
		server.DefaultHTTPErrorHandler(err, c)
	}

	// Route => handler
	//// Exercise
	api := server.Group("/api")
	api.GET("/exercise", handler.GetExerciseNames)
	api.GET("/exercise/:name", handler.GetExercise)
	api.DELETE("/exercise/:name", handler.DeleteExercise)

	//// Workout
	api.GET("/workout", handler.GetWorkoutNames)
	api.DELETE("/workout/:name", handler.DeleteWorkout)
	api.GET("/workout/:name", handler.GetWorkout)
	api.POST("/workout", handler.SubmitWorkoutPerformed)
	api.GET("/workout/:name/:date", handler.GetWorkoutPerformed)
	api.DELETE("/workout/:name/:date", handler.DeleteWorkoutPerformed)

	//// Program
	api.POST("/program", handler.CreateProgram)
	api.GET("/program", handler.GetProgramNames)
	api.GET("/program/:name", handler.GetProgram)

	//// Composition
	api.POST("/composition", handler.SubmitComposition)
	api.GET("/composition/:date", handler.GetComposition)
	api.DELETE("/composition/:date ", handler.DeleteComposition)

	//// Nutrition
	api.POST("/nutrition", handler.SubmitNutrition)
	api.GET("/nutrition/:date", handler.GetNutrition)
	api.DELETE("/nutrition/:date", handler.DeleteNutrition)

	//// Training Calculations Table
	api.GET("/training/weight", handler.GetTrainingWeight)
	api.GET("/training/max", handler.GetTrainingMax)

	// Start server
	server.Logger.Fatal(server.Start(":1323"))
}
