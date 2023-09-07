package usecase

import (
	"echojson/models"
	"echojson/repository"
	"echojson/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetCart(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, err := utils.ExtractAccessClaims(tokenStr)
	if err != nil {
		response := models.Response{
			Message: "MISSING TOKEN!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	var cart models.Cart
	if err := c.Bind(&cart); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	result, err := repository.GetCart(cart.ID, user)
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
		UserID:  user,
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  result,
		Errors:  nil,
	}
	return c.JSON(http.StatusOK, response)
}

func InsertCart(c echo.Context) error {
	var cart models.RequestCart
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, _ := utils.ExtractAccessClaims(tokenStr)

	// var cart = new(models.Requestcart)
	err = utils.BindValidateStruct(c, &cart)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	result, err := repository.AddCart(cart, user)
	if err != nil {
		response := models.Response{
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

func DeleteCart(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, _ := utils.ExtractAccessClaims(tokenStr)

	var cart models.Cart
	err = utils.BindValidateStruct(c, &cart)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result, err := repository.DeleteCart(cart, user)
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
