package usecase

import (
	"echojson/controller"
	"echojson/models"
	"echojson/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

//                                       usecase is where all the bind functions go (user input is handled here, along with any related errors)
//													It passes arguments to functions defined in controller

func Login(c echo.Context) error {
	var users models.User // declare "user" as new User struct
	if err := c.Bind(&users); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	container, err := controller.GETLogin(users.Username)
	//container = login(user.password)
	containerSlice := container[0]                                                                                 //refers to the first user instance (array stores user instances) array looks like: container[{Username,Password,Firstname}]
	if err := bcrypt.CompareHashAndPassword([]byte(containerSlice.Password), []byte(users.Password)); err != nil { //move to usecase
		// If the two passwords don't match, return a 401 status.

		// return echo.NewHTTPError(http.StatusUnauthorized, "Password is incorrect")
		fmt.Println("typed username: ", users.Username) //debug
		fmt.Println("typed password: ", users.Password)
		fmt.Println("container username: ", containerSlice.Username)
		fmt.Println("container password: ", containerSlice.Password)
		return echo.ErrUnauthorized
	}

	//convert container to User instance so it can be passed into generate access token (gen access token needs User struct for "user.First_name")
	userInstance := models.User{
		First_name: containerSlice.First_name,
		Username:   containerSlice.Username,
		Password:   containerSlice.Password,
	}
	token, err := utils.GenerateAccessToken(userInstance)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Something happened with token generation")
	}
	firstname := string(container[0].First_name)
	res := models.Response{
		Message: "OK, welcome " + firstname, //returns firstname
		Status:  "OK",
		Result:  token,
	}
	return c.JSON(http.StatusOK, res)
}

/*===============================


GET


=================================*/

func GETDataProducts(c echo.Context) error { //bisa pass value kesini?

	var items models.Item // declare "user" as new User struct
	if err := c.Bind(&items); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	var respon models.Response

	usera, err := controller.GetProducts(items.Product_id)
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

func GETTransactions(c echo.Context) error {
	var transactions models.Transaction
	if err := c.Bind(&transactions); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	result, err := controller.GetTransaction(transactions.Transaction_id)
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
		Errors:  err.Error(),
	}
	return c.JSON(http.StatusOK, response)
}
func GETUsers(c echo.Context) error {

	var users models.User
	if err := c.Bind(&users); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	result := models.Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  controller.GetUsers(users.ID),
	}
	return c.JSON(http.StatusOK, result)
}

/*===============================


INSERT


=================================*/

func INSERTProducts(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	user, _ := utils.ExtractAccessClaims(tokenStr)
	fmt.Println(tokenStr)
	if err != nil {
		fmt.Println("something happened with ExtractAccessClaims")
		fmt.Println(err)
	}
	//bind json input to items
	var items []models.Item
	if err := c.Bind(&items); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	//result is when we pass items into AddProducts
	response := models.Response{
		Message: "OK, added by " + user,
		Status:  "OK",
		Result:  controller.AddProducts(items),
	}

	return c.JSON(http.StatusOK, response)

}

func INSERTTransactions(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, _ := utils.ExtractAccessClaims(tokenStr)

	var transactions []models.Transaction // declare "user" as new User struct

	if err := c.Bind(&transactions); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result := models.Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  controller.AddTransaction(transactions, user),
	}
	return c.JSON(http.StatusOK, result)
}
func INSERTUsers(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, _ := utils.ExtractAccessClaims(tokenStr)

	var users []models.User // declare "user" as new User struct

	if err := c.Bind(&users); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result, err := controller.AddUsers(users, user)
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
		Errors:  err.Error(),
	}
	return c.JSON(http.StatusOK, response)
}

/*===============================


UPDATE


=================================*/

func UPDATEProducts(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, _ := utils.ExtractAccessClaims(tokenStr)
	var itemContainer models.Item

	if err := c.Bind(&itemContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	result, err := controller.UpdateProducts(itemContainer, user)
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
		Errors:  err.Error(),
	}
	return c.JSON(http.StatusOK, response)
}

func UPDATETransactions(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, _ := utils.ExtractAccessClaims(tokenStr)

	var transaction models.Transaction // declare "transaction" as new Transaction struct, this is to contain the JSON input
	if err := c.Bind(&transaction); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	result, err := controller.UpdateTransaction(transaction, user)
	if err != nil {
		//TODO HANDLE ERROR
	}
	return c.JSON(http.StatusOK, result)
}

func UPDATEUsers(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, _ := utils.ExtractAccessClaims(tokenStr)

	var userContainer models.User // declare "users" as new User struct for binding
	if err := c.Bind(&userContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	result, err := controller.UpdateUsers(userContainer, user)
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
		Errors:  err.Error(),
	}
	return c.JSON(http.StatusOK, response)
}

/*===============================


DELETE


=================================*/

func DELETEProducts(c echo.Context) error { //wrapper for DeleteProducts
	//extract user frm tokens
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, _ := utils.ExtractAccessClaims(tokenStr)
	//container for the json request
	var itemContainer models.Item
	if err := c.Bind(&itemContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user. if err is not nil, print out the response struct
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result, err := controller.DeleteProducts(itemContainer, user)
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
		Errors:  err.Error(),
	}
	return c.JSON(http.StatusOK, response) //outputs a response struct
}

func DELETETransactions(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, _ := utils.ExtractAccessClaims(tokenStr)

	var transContainer models.Transaction // declare "user" as new User struct
	if err := c.Bind(&transContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	result, err := controller.DeleteTransaction(transContainer, user)
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
		Errors:  err.Error(),
	}

	return c.JSON(http.StatusOK, response)
}

func DELETEUsers(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		return err
	}
	user, _ := utils.ExtractAccessClaims(tokenStr)

	var userContainer models.User
	if err := c.Bind(&userContainer); err != nil {
		fmt.Println("Bind Error:", err)
	}

	result, err := controller.DeleteUsers(userContainer, user)
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
		Errors:  err.Error(),
	}
	return c.JSON(http.StatusOK, response)
}
