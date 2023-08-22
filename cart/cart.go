package cart

import (
	"database/sql"
	"echojson/jwtTest"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

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

/*
=================================

GET PRODUCTS (SELECT)

==================================
*/
func GetProducts(c echo.Context) error {
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

func AddProducts(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtTest.JwtCustomClaims)
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
func UpdateProducts(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtTest.JwtCustomClaims)
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
func DeleteProducts(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtTest.JwtCustomClaims)
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
		Message: "SUCCESS deleted by" + name,
		Status:  "SUCCESS",
		Result:  items,
	}

	return c.JSON(http.StatusCreated, response)
}
