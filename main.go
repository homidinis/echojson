package main

//todo: update, delete
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "echojson/auth"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

//	type UserFull struct {
//		Username   string `json:"username"`
//		Password   string `json:"password"`
//	}
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}
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

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	config := echojwt.Config{ //configures "restricted" to restrict requests without an Authorization token
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims) //use custom claims jwtCustomClaims
		},
		SigningKey: []byte("secret"),
	}

	r := e.Group("/restricted") //anything that uses r is a group
	r.Use(echojwt.WithConfig(config))
	e.POST("/login", login)
	e.GET("/showProducts", getProducts)
	// /restricted/*
	r.GET("/updateProducts", updateProducts)
	r.POST("/addProducts", addProducts)
	r.POST("/deleteProducts", deleteProducts)

	e.Logger.Fatal(e.Start(":9000"))

}

/*
==================================

# LOGIN

==================================
*/
var container User

func login(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")

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

	err = auth.GenerateAccessToken(container.Username, c)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Token is incorrect")
	}
	res := Response{
		Message: "OK, welcome " + container.First_name, //returns firstname
		Status:  "OK",
		Result:  t,
	}
	return c.JSON(http.StatusOK, res)
}

/*
=================================

GET PRODUCTS (SELECT)

==================================
*/
func getProducts(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	statement, err := db.Prepare("SELECT * FROM products")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := statement.Query()
	if err != nil {
		fmt.Println(err)
	}
	var products []Item
	for rows.Next() {
		var product Item
		err := rows.Scan(&product.Name, &product.Description, &product.Price, &product.Product_id)
		if err != nil {
			fmt.Println(err)
		}
		products = append(products, product)
	}
	res := Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  products,
	}
	return c.JSON(http.StatusOK, res)
}

/*
========================================

# REPLACE SQL

========================================
*/
func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

/*========================================

ADD PRODUCTS

========================================*/

func addProducts(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name //returns dummyuser instead of Dummy A?
	var items []Item    // declare "user" as new User struct
	if err := c.Bind(&items); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	vals := []interface{}{}
	sqlStr := `INSERT INTO products (name, description, price) VALUES `

	for _, row := range items { //index,name_of_
		sqlStr += "(?, ?, ?),"
		vals = append(vals, row.Name, row.Description, row.Price)
	}
	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	// replacing ? with $n for postgres
	sqlStr = ReplaceSQL(sqlStr, "?")
	//prepare the statement
	statement, _ := db.Prepare(sqlStr)

	//format all vals at once
	_, err = statement.Exec(vals...)
	if err != nil {
		fmt.Println("Bind Error:", err)
		return err
	}

	response := Response{
		Message: "SUCCESS added by " + name,
		Status:  "SUCCESS",
		Result:  items,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
================================

# UPDATE PRODUCTS

====================================
*/
func updateProducts(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")

	var itemContainer Item // declare "user" as new User struct
	if err := c.Bind(&itemContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	statement, err := db.Prepare(`UPDATE public.products SET name=$1, description=$2, price=$3 WHERE product_id=$5 RETURNING product_id;`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return err
	}
	var items Item
	err = statement.QueryRow(&itemContainer.Name, &itemContainer.Description, &itemContainer.Price, &itemContainer.Product_id).Scan(&items.Product_id)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return err
	}

	response := Response{
		Message: "SUCCESS updated by" + name,
		Status:  "SUCCESS",
		Result:  items,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
================================

# DELETE PRODUCTS

====================================
*/
func deleteProducts(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")

	var itemContainer Item // declare "user" as new User struct
	if err := c.Bind(&itemContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	statement, err := db.Prepare(`UPDATE public.products SET name=$1, description=$2, price=$3 WHERE product_id=$5 RETURNING product_id;`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return err
	}
	var items Item
	err = statement.QueryRow(&itemContainer.Name, &itemContainer.Description, &itemContainer.Price, &itemContainer.Product_id).Scan(&items.Product_id)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return err
	}

	response := Response{
		Message: "SUCCESS updated by" + name,
		Status:  "SUCCESS",
		Result:  items,
	}

	return c.JSON(http.StatusCreated, response)
}
