package usecase

import (
	"echojson/models"
	"echojson/repository"
	"echojson/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

//                                       usecase is where all the bind functions go (user input is handled here, along with any related errors)
//													It passes arguments to functions defined in repository

func Login(c echo.Context) error {
	var users models.User // declare "user" as new User struct
	if err := c.Bind(&users); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	container, err := repository.GetLogin(users.Username)
	if err != nil {
		fmt.Println("getlogin error")
		return err
	}

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
	if containerSlice.Admin {
		userInstance := models.User{
			ID:       containerSlice.ID,
			Username: containerSlice.Username,
			Password: containerSlice.Password,
			Admin:    true,
		}
		token, err := utils.GenerateAccessTokenAdmin(userInstance)
		fmt.Println(userInstance.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Something happened with token generation")
		}
		firstname := string(container[0].First_name)
		res := models.Response{
			UserID:  userInstance.ID,
			Message: "Admin token generated OK, welcome " + firstname, //returns firstname
			Status:  "OK",
			Result:  token,
		}
		return c.JSON(http.StatusOK, res)
	} else {
		userInstance := models.User{
			ID:       containerSlice.ID,
			Username: containerSlice.Username,
			Password: containerSlice.Password,
			Admin:    false,
		}
		token, err := utils.GenerateAccessTokenUser(userInstance)
		fmt.Println(userInstance.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Something happened with token generation")
		}
		firstname := string(container[0].First_name)
		res := models.Response{
			UserID:  userInstance.ID,
			Message: "User token generated OK, welcome " + firstname, //returns firstname
			Status:  "OK",
			Result:  token,
		}
		return c.JSON(http.StatusOK, res)
	}
}
