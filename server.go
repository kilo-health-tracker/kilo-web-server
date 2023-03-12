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
	server.GET("/exercise", handler.GetExerciseNames)
	server.GET("/exercise/:name", handler.GetExercise)
	server.DELETE("/exercise/:name", handler.DeleteExercise)

	//// Workout
	server.GET("/workout", handler.GetWorkoutNames)
	server.DELETE("/workout/:name", handler.DeleteWorkout)
	server.GET("/workout/:name", handler.GetWorkout)
	server.GET("/workout/:name/:date", handler.GetWorkoutPerformed)
	server.DELETE("/workout/:name/:date", handler.DeleteWorkoutPerformed)

	//// Program
	server.POST("/program", handler.CreateProgram)
	server.GET("/program", handler.GetProgramNames)
	server.GET("/program/:name", handler.GetProgram)

	//// Composition
	server.POST("/composition", handler.SubmitComposition)
	server.GET("/composition/:date", handler.GetComposition)
	server.DELETE("/composition/:date ", handler.DeleteComposition)

	//// Nutrition
	server.POST("/nutrition", handler.SubmitNutrition)
	server.GET("/nutrition/:date", handler.GetNutrition)
	server.DELETE("/nutrition/:date", handler.DeleteNutrition)

	// Start server
	server.Logger.Fatal(server.Start(":1323"))
}
