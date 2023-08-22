package main

//todo: update, delete
import (
	"database/sql"
	"echojson/jwtTest"
	"echojson/products"
	"echojson/transactions"
	"echojson/users"
	"fmt"
	"net/http"
	"time"

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

type User struct {
	First_name string `json:"firstname"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}
type Response struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Result  interface{} `json:"result"`
}
type Item struct {
	Product_id  int    `json:"product_id"`
	Name        string `json:"item"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

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
			return new(jwtTest.JwtCustomClaims) //use custom claims jwtCustomClaims
		},
		SigningKey: []byte("secret"),
	}

	r := e.Group("/restricted") //anything that uses r is a group
	r.Use(echojwt.WithConfig(config))
	e.POST("/login", login)
	e.GET("/showProducts", products.GetProducts)
	e.GET("/showUsers", users.GetUsers)
	e.GET("/showTransactions", transactions.GetTransaction)
	// /restricted/*

	r.GET("/updateProducts", products.UpdateProducts)
	r.POST("/addProducts", products.AddProducts)
	r.DELETE("/deleteProducts", products.DeleteProducts)

	r.POST("/updateUsers", users.UpdateUsers)
	r.POST("/addUsers", users.AddUsers)
	r.DELETE("/deleteUsers", users.DeleteUsers)

	r.POST("/updateTransactions", transactions.UpdateTransaction)
	r.POST("/addTransactions", transactions.AddTransaction)
	r.DELETE("/deleteTransactions", transactions.DeleteTransaction)

	e.Logger.Fatal(e.Start(":9000"))

}

/*
==========================

# GENERATE ACCESS TOKEN

===========================
*/

const (
	refreshTokenCookieName = "refresh-token"
	accessTokenCookieName  = "access-token"
	jwtSecretKey           = "secret"
	jwtRefreshSecretKey    = "secret"
)

func GetJWTSecret() string {
	return jwtSecretKey
}
func GetJWTRefresh() string {
	return jwtRefreshSecretKey
}

func GenerateAccessToken(user User) (string, error) {
	return GenerateToken(user, []byte(GetJWTSecret()))
}

/*
============================

# GENERATE TOKEN

==============================
*/
func GenerateToken(user User, secret []byte) (string, error) {

	claims := &jwtTest.JwtCustomClaims{ //need to put the struct in a common file exportable by main AND Products or it will complaim
		Name:  user.First_name,
		Admin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTErrorChecker(err error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("login"))
}

/*
==================================

# LOGIN

==================================
*/
var container User

func login(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")
	if err != nil {
		fmt.Println("DB Connection error: ", err)
		return err
	}
	var user User
	if err := c.Bind(&user); err != nil {
		fmt.Println("Bind error: ", err)
		return err
	}
	statement, err := db.Prepare("SELECT username, password, first_name FROM users WHERE username=$1 AND password=$2")
	if err != nil {
		fmt.Println("Prepare err :", err)
		return err
	}
	err = statement.QueryRow(user.Username, user.Password).Scan(&container.Username, &container.Password, &container.First_name) //fetch name of who just logged in
	if err != nil {
		fmt.Println("Query err :", err)
		return echo.ErrUnauthorized
	}

	if user.Username != container.Username || user.Password != container.Password { //if input doesn't equal what's in the database
		fmt.Println("typed username: ", user.Username) //debug
		fmt.Println("typed password: ", user.Password)
		fmt.Println("container username: ", container.Username)
		fmt.Println("container password: ", container.Password)
		return echo.ErrUnauthorized //get mad
	}
	//convert container to User instance so it can be passed into generate access token (gen access token needs User struct for "user.First_name")
	userInstance := User{
		First_name: container.First_name,
		Username:   container.Username,
		Password:   container.Password,
	}
	token, err := GenerateAccessToken(userInstance)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Something happened with token generation")
	}

	res := Response{
		Message: "OK, welcome " + container.First_name, //returns firstname
		Status:  "OK",
		Result:  token,
	}
	return c.JSON(http.StatusOK, res)
}
