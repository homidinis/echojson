package users

//todo: update, delete
import (
	"database/sql"

	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}
type User struct {
	ID         string `json:"id"`
	Age        string `json:"age"`
	First_name string `json:"firstname"`
	Last_name  string `json:"lastname"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Group      string `json:"group"`
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
func GetUsers(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	statement, err := db.Prepare(`SELECT id, age, first_name, last_name, email, username, "group" FROM public.users;`)
	if err != nil {
		fmt.Println(err)
	}
	rows, err := statement.Query()
	if err != nil {
		fmt.Println(err)
	}
	var userContainer []User
	for rows.Next() {
		var users User
		err := rows.Scan(&users.ID, &users.Age, &users.First_name, &users.Last_name, &users.Email, &users.Username, &users.Group)
		if err != nil {
			fmt.Println(err)
		}
		userContainer = append(userContainer, users)
	}
	res := Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  userContainer,
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

func AddUsers(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name //returns dummyuser instead of Dummy A?
	var users []User    // declare "user" as new User struct
	if err := c.Bind(&users); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	vals := []interface{}{}
	sqlStr := `INSERT INTO public.users(id, age, first_name, last_name, email, username, password, "group") VALUES `

	for _, row := range users { //index,name_of_
		sqlStr += " (?, ?, ?, ?, ?, ?, ?, ?),"
		vals = append(vals, row.ID, row.Age, row.First_name, row.Last_name, row.Email, row.Username, row.Password, row.Group)
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
		Result:  users,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
================================

# UPDATE PRODUCTS

====================================
*/
func UpdateUsers(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")

	var users User // declare "users" as new User struct
	if err := c.Bind(&users); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	statement, err := db.Prepare(`UPDATE public.users SET age=$1, first_name=$2, last_name=$3, email=$4, username=$5, password=$6, "group"=$7 WHERE id=$5 RETURNING id;`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return err
	}
	var items Item
	err = statement.QueryRow(&users.Age, &users.First_name, &users.Last_name, &users.Email, &users.Username, &users.Password, &users.Group).Scan(&items.Product_id)
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
func DeleteUsers(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")

	var userContainer User // declare "user" as new User struct
	if err := c.Bind(&userContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	statement, err := db.Prepare(`DELETE FROM public.users WHERE id=$1 RETURNING id`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return err
	}
	var users User
	err = statement.QueryRow(&userContainer.ID).Scan(&users.ID)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return err
	}

	response := Response{
		Message: "SUCCESS updated by" + name,
		Status:  "SUCCESS",
		Result:  users,
	}

	return c.JSON(http.StatusCreated, response)
}
