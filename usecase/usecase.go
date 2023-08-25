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

func INSERTProducts(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	//extract user from token
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
	result := controller.AddProducts(items)
	response := models.Response{
		Message: "OK, added by " + user,
		Status:  "OK",
		Result:  result,
	}

	return c.JSON(http.StatusOK, response)

}

func UPDATEProducts(c echo.Context) error {
	var itemContainer models.Item

	if err := c.Bind(&itemContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}

	result := controller.UpdateProducts(itemContainer) //updateproducts takes itemContainer (models.Item struct) and outputs struct
	return c.JSON(http.StatusOK, result)
}
