package main

//todo: update, delete
import (
	"echojson/models"
	"echojson/usecase"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

//	type UserFull struct {
//		Username   gstring `json:"username"`
//		Password   string `json:"password"`
//	}

/*
============================

# MAIN

=============================
*/
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	config := echojwt.Config{ //configures "restricted" to restrict requests without an Authorization token
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(models.JwtCustomClaims) //use custom claims jwtCustomClaims
		},
		SigningKey: []byte("secret"),
	}

	r := e.Group("/restricted") //anything that uses r is a group
	r.Use(echojwt.WithConfig(config))
	e.POST("/login", usecase.Login)
	e.GET("/showProducts", usecase.GetProducts)
	e.GET("/showUsers", usecase.GetUsers)
	e.GET("/showTransactions", usecase.GetTransactions)
	e.GET("/showCarts", usecase.GetCart)
	e.POST("/addUsers", usecase.InsertUsers)
	e.POST("/checkout", usecase.Checkout)
	// /restricted/*
	r.POST("/addCarts", usecase.InsertCart)
	r.POST("/updateCarts", usecase.UpdateCart)
	r.GET("/updateProducts", usecase.UpdateProducts)
	r.POST("/addProducts", usecase.InsertProducts)
	r.DELETE("/deleteProducts", usecase.DeleteProducts)

	r.POST("/updateUsers", usecase.UpdateUsers)

	r.DELETE("/deleteUsers", usecase.DeleteUsers)

	r.POST("/addTransactions", usecase.InsertTransactions)
	r.DELETE("/deleteTransactions", usecase.DeleteTransactions)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		report, ok := err.(*echo.HTTPError)
		if !ok {
			report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if castedObject, ok := err.(validator.ValidationErrors); ok {
			for _, err := range castedObject {
				switch err.Tag() {
				case "required":
					report.Message = fmt.Sprintf("%s is required", err.Field())

				case "email":
					report.Message = fmt.Sprintf("%s is not a valid email", err.Field())

				case "gte":
					report.Message = fmt.Sprintf("%s value must be greater than %s", err.Field(), err.Param())

				case "lte":
					report.Message = fmt.Sprintf("%s value must be lower than %s", err.Field(), err.Param())
				}
				break
			}
		}

		c.Logger().Error(report)
		c.JSON(report.Code, report)
	}

	e.Logger.Fatal(e.Start(":9000"))

}
