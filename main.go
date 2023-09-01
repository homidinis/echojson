package main

//todo: update, delete
import (
	"echojson/models"
	"echojson/usecase"

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
func main() {
	e := echo.New()

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
	e.GET("/showProducts", usecase.GETDataProducts)
	e.GET("/showUsers", usecase.GETUsers)
	e.GET("/showTransactions", usecase.GETTransactions)
	// /restricted/*

	r.GET("/updateProducts", usecase.UPDATEProducts)
	r.POST("/addProducts", usecase.INSERTProducts)
	r.DELETE("/deleteProducts", usecase.DELETEProducts)

	r.POST("/updateUsers", usecase.UPDATEUsers)
	r.POST("/addUsers", usecase.INSERTUsers)
	r.DELETE("/deleteUsers", usecase.DELETEUsers)

	r.POST("/updateTransactions", usecase.UPDATETransactions)
	r.POST("/addTransactions", usecase.INSERTTransactions)
	r.DELETE("/deleteTransactions", usecase.DELETETransactions)

	e.Logger.Fatal(e.Start(":9000"))

}
