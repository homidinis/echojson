package usecase

import (
	"echojson/models"
	"echojson/repository"
	"echojson/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {

	var users models.User
	if err := c.Bind(&users); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result := models.Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  repository.GetUser(users.ID),
		Errors:  nil,
	}
	return c.JSON(http.StatusOK, result)
}

func InsertUsers(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		response := models.Response{
			Message: "ExtractToken MISSING TOKEN!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	user, _, err := utils.ExtractAccessClaims(tokenStr)

	var users = new(models.RequestUser)

	err = utils.BindValidateStruct(c, &users)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result, err := repository.AddUser(*users, user)
	if err != nil {
		response := models.Response{
			UserID:  user,
			Message: "ERROR in AddUser calling",
			Status:  "ERROR",
			Result:  result,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	response := models.Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  result,
		Errors:  nil,
	}
	return c.JSON(http.StatusOK, response)
}

func UpdateUsers(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		response := models.Response{
			Message: "ExtractToken MISSING TOKEN!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	user, _, err := utils.ExtractAccessClaims(tokenStr)

	var userContainer models.User // declare "users" as new User struct for binding
	err = utils.BindValidateStruct(c, &userContainer)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result, err := repository.UpdateUser(userContainer, user)
	if err != nil {
		fmt.Println("Exec Error:", err)
		response := models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  result,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	response := models.Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  result,
		Errors:  nil,
	}
	return c.JSON(http.StatusOK, response)
}

func DeleteUsers(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		response := models.Response{
			Message: "ExtractToken MISSING TOKEN!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	user, _, err := utils.ExtractAccessClaims(tokenStr)

	var userContainer models.User
	err = utils.BindValidateStruct(c, &userContainer)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	result, err := repository.DeleteUser(userContainer, user)
	if err != nil {
		response := models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  result,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	response := models.Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  result,
		Errors:  nil,
	}
	return c.JSON(http.StatusOK, response)
}
