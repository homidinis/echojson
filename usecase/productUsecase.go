package usecase

import (
	"echojson/models"
	"echojson/repository"
	"echojson/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

/*===============================

GET

=================================*/

func GetProducts(c echo.Context) error { //bisa pass value kesini?

	var items models.Item // declare "user" as new User struct
	if err := c.Bind(&items); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	var respon models.Response

	usera, err := repository.GetProducts(items.Product_id)
	fmt.Println("product id: ")
	userBytes, _ := json.MarshalIndent(usera, "", "\t")
	fmt.Println(string(userBytes))

	if err != nil {
		log.Fatal(err)
		respon.Message = "Error"
		respon.Status = "Failed to acquire product data"
		respon.Result = err.Error()
	}
	respon.Message = "SUCCESS"
	respon.Status = "Acquired product data"
	respon.Result = usera
	return c.JSON(http.StatusCreated, respon)

}

/*===============================

INSERT

=================================*/

func InsertProducts(c echo.Context) error {
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
	user, isAdmin, err := utils.ExtractAccessClaims(tokenStr)
	fmt.Println(tokenStr)
	if err != nil {
		fmt.Println("something happened with ExtractAccessClaims")
		fmt.Println(err)
	}
	if !isAdmin {
		response := models.Response{
			Message: "isAdmin not true!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  "Unauthorized",
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	var items models.RequestItem
	//bind json input to items
	err = utils.BindValidateStruct(c, &items)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	//result is when we pass items into AddProducts

	result, err := repository.AddProducts(items)
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
		Message: "OK",
		UserID:  user,
		Status:  "OK",
		Result:  result,
		Errors:  nil,
	}

	return c.JSON(http.StatusOK, response)

}

/*===============================

UPDATE

=================================*/

func UpdateProducts(c echo.Context) error {
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
	var itemContainer models.Item

	err = utils.BindValidateStruct(c, &itemContainer)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result, err := repository.UpdateProducts(itemContainer, user)
	if err != nil {
		response := models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  result,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	if err := c.Validate(itemContainer); err != nil {
		return err
	}
	response := models.Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  result,
		Errors:  nil,
	}
	return c.JSON(http.StatusOK, response)
}

/*===============================

DELETE

=================================*/

func DeleteProducts(c echo.Context) error { //wrapper for DeleteProducts
	//extract user frm tokens
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
	//container for the json request
	var itemContainer models.Item
	err = utils.BindValidateStruct(c, &itemContainer)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result, err := repository.DeleteProducts(itemContainer, user)
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
	return c.JSON(http.StatusOK, response) //outputs a response struct
}
