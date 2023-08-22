package transactions

//todo: update, delete
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
type Transaction struct {
	ID               int    `json:"id"`
	Transaction_id   string `json:"transaction_id"`
	Product_id       string `json:"product_id"`
	Transaction_date string `json:"transaction_date"`
	User_id          string `json:"user_id"`
	Payment_method   string `json:"payment_method"`
	Quantity         string `json:"quantity"`
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
func GetTransaction(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	statement, err := db.Prepare(`SELECT id, transaction_id, product_id, transaction_date, user_id, payment_method, quantity FROM public.transaction_history;`)
	if err != nil {
		fmt.Println(err)
	}
	rows, err := statement.Query()
	if err != nil {
		fmt.Println(err)
	}
	var trsContainer []Transaction

	for rows.Next() {
		var trans Transaction
		err := rows.Scan(&trans.ID, &trans.Transaction_id, &trans.Product_id, &trans.Transaction_date, &trans.User_id, &trans.Payment_method, &trans.Quantity)
		if err != nil {
			fmt.Println(err)
		}
		trsContainer = append(trsContainer, trans)
	}
	res := Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  trsContainer,
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

func AddTransaction(c echo.Context) error {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtTest.JwtCustomClaims)
	name := claims.Name            //returns dummyuser instead of Dummy A?
	var transactions []Transaction // declare "user" as new User struct
	if err := c.Bind(&transactions); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	vals := []interface{}{}
	sqlStr := `INSERT INTO public.transaction_history(transaction_id, product_id, transaction_date, user_id, payment_method, quantity) VALUES`

	for _, row := range transactions { //index,name_of_ ; for every data inputted in, run loop
		sqlStr += " (?, ?, ?, ?, ?, ?),"
		vals = append(vals, row.Transaction_id, row.Product_id, row.Transaction_date, row.User_id, row.Payment_method, row.Quantity)
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
		Result:  vals,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
================================

# UPDATE PRODUCTS

====================================
*/
func UpdateTransaction(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtTest.JwtCustomClaims)
	name := claims.Name
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")

	var transaction Transaction // declare "transaction" as new Transaction struct
	if err := c.Bind(&transaction); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	statement, err := db.Prepare(`UPDATE public.transaction_history SET id=$1, transaction_id=$2, product_id=$3, transaction_date=$4, user_id=$5, payment_method=$6, quantity=$7 WHERE id=$8; RETURNING id;`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return err
	}
	var transactions Transaction
	err = statement.QueryRow().Scan(&transactions.ID)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return err
	}

	response := Response{
		Message: "SUCCESS updated by" + name,
		Status:  "SUCCESS",
		Result:  transactions,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
================================

# DELETE PRODUCTS

====================================
*/
func DeleteTransaction(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtTest.JwtCustomClaims)
	name := claims.Name
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mkp_demo sslmode=disable")

	var transContainer Transaction // declare "user" as new User struct
	if err := c.Bind(&transContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	statement, err := db.Prepare(`DELETE FROM public.transaction_history WHERE transaction_id=$1 RETURNING transaction_id`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return err
	}
	var transactions Transaction
	err = statement.QueryRow(&transContainer.ID).Scan(&transactions.Transaction_id)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return err
	}

	response := Response{
		Message: "SUCCESS updated by" + name,
		Status:  "SUCCESS",
		Result:  transactions,
	}

	return c.JSON(http.StatusCreated, response)
}
