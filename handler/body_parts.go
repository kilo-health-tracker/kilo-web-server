package handler

import (
	"context"
	//"fmt"
	"net/http"
	//"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/kilo-health-tracker/kilo-web-server/utils"
	"github.com/kilo-health-tracker/kilo-web-server/models"
)

// Get BodyPart details
// func GetBodyPart(c echo.Context) error {
// 	queries, err := utils.GetQueryInterface()
// 	if err != nil {
// 		return err
// 	}

// 	BodyPart, err := queries.GetBodyPart(context.Background(), c.Param("name"))
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(http.StatusOK, BodyPart)
// }

// Submits an BodyPart definition
func SubmitBodyPart(c echo.Context) error {
	var requestBody models.SubmitBodyPartParams
	queries, err := utils.GetQueryInterface()
	if err != nil {
		return err
	}

	// bind request body to variable given
	if err := c.Bind(&requestBody); err != nil {
		return err
	}

	ctx := context.Background()

	entry := models.SubmitBodyPartParams{
		Name:         requestBody.Name,
		Region:       requestBody.Region,
		UpperOrLower: requestBody.UpperOrLower,
	}

	bodyPart, err := queries.SubmitBodyPart(ctx, entry)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, bodyPart)
}

// Delete program
// func DeleteBodyPart(c echo.Context) error {
// 	queries, err := utils.GetQueryInterface()
// 	if err != nil {
// 		return err
// 	}

// 	if err := queries.DeleteBodyPart(context.Background(), c.Param("name")); err != nil {
// 		return err
// 	}

// 	return c.JSON(http.StatusOK, GenericResponse{fmt.Sprintf("Successfully deleted BodyPart: %s", c.Param("name"))})
// }
