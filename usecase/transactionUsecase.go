package usecase

import (
	"echojson/models"
	"echojson/repository"
	"echojson/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetTransactions(c echo.Context) error {
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

	var transactions models.Transaction
	if err := c.Bind(&transactions); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	result, err := repository.GetTransaction(transactions.Transaction_id, user)
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

func InsertTransactions(c echo.Context) error {
	var transactions models.RequestTransaction
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, isAdmin, err := utils.ExtractAccessClaims(tokenStr)
	if !isAdmin {
		response := models.Response{
			Message: "isAdmin not true!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  "Unauthorized",
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	// var transactions = new(models.RequestTransaction)
	err = utils.BindValidateStruct(c, &transactions)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			// Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	result, err := repository.AddTransaction(transactions, user)
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

func DeleteTransactions(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, isAdmin, err := utils.ExtractAccessClaims(tokenStr)
	if !isAdmin {
		response := models.Response{
			Message: "isAdmin not true!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  "Unauthorized",
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	var transactions models.Transaction
	err = utils.BindValidateStruct(c, &transactions)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result, err := repository.DeleteTransaction(transactions, user)
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
