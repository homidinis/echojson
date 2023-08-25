package controller

import (
	"echojson/db"
	"echojson/models"
	"echojson/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

/*
==================================

# LOGIN

==================================
*/
var container models.User

func GETLogin(username string) (containerArray []models.User, err error) { //declare as array of User struct
	db := db.Conn()
	//1. grabs username and password WHERE typed-username
	//2. dumps username, password, firstname from database into container
	//3. compare password in container to typed-password
	//4. if no errors, generate token
	statement, err := db.Prepare("SELECT username, password, first_name FROM users WHERE username=$1") //only select by Username
	if err != nil {
		fmt.Println("Prepare err in controller:", err)
	}
	var container models.User
	err = statement.QueryRow(username).Scan(&container.Username, &container.Password, &container.First_name) //container = [user{Username,Password,Firstname}]; scan scans into each of them
	if err != nil {
		fmt.Println("Query err in controller:", err)
	}
	containerArray = append(containerArray, container)
	return
}

/*
=================================

GET PRODUCTS (SELECT)

==================================
*/
func GetProducts(id int) (products []models.Item, err error) {
	db := db.Conn()

	var data []interface{}
	query := "SELECT name, description, price, product_id FROM products"

	if id != 0 { //if id is not presented
		query += " WHERE product_id=$1" //append "where" to query
		data = append(data, id)         //then append the id arg to the data interface. in the case of there being a lot of arguments for a lot of WHERE conditions
	}
	fmt.Println(query)
	rows, err := db.Query(query, data...) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() { //for every row result, run Scan then Append the result into the products struct
		var product models.Item
		err := rows.Scan(&product.Name, &product.Description, &product.Price, &product.Product_id)
		if err != nil {
			fmt.Println(err)
		}
		products = append(products, product)
	} //products already declared as return value, instead of old method returning c.JSON(products)
	return
}

/*========================================

ADD PRODUCTS

========================================*/

func AddProducts(items []models.Item) (err error) {
	db := db.Conn()

	//1. declare array of Item struct (items)
	//2. bind items to json input
	//3. declare vals as an array
	//4. loop the items array, append into vals each Name, Description, price from items
	//5. return null if OK, error if error

	vals := []interface{}{}
	sqlStr := `INSERT INTO products (name, description, price) VALUES `

	for _, row := range items { //index,name_of_
		sqlStr += "(?, ?, ?),"
		vals = append(vals, row.Name, row.Description, row.Price)
	}
	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	// replacing ? with $n for postgres
	sqlStr = utils.ReplaceSQL(sqlStr, "?")
	//prepare the statement
	statement, _ := db.Prepare(sqlStr)

	//format all vals at once
	_, err = statement.Exec(vals...)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}
	return
}

/*
================================

# UPDATE PRODUCTS

====================================
*/
func UpdateProducts(itemContainer models.Item) (response models.Response) { //returns response
	db := db.Conn() // declare "user" as new User struct

	statement, err := db.Prepare(`UPDATE public.products SET name=$1, description=$2, price=$3 WHERE product_id=$4 RETURNING product_id;`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		response = models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  err.Error(),
		}
		return response
	}
	var items models.Item
	err = statement.QueryRow(&itemContainer.Name, &itemContainer.Description, &itemContainer.Price, &itemContainer.Product_id).Scan(&items.Product_id)
	if err != nil {
		fmt.Println("Exec Error:", err)
		response = models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  err.Error(),
		}
		return response
	}

	response = models.Response{
		Message: "SUCCESS updated by NIL",
		Status:  "SUCCESS",
		Result:  items.Product_id,
	}

	return response
}

/*
================================

# DELETE PRODUCTS

====================================
*/
func DeleteProducts(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	name := claims.Name
	db := db.Conn()
	var itemContainer models.Item // declare "user" as new User struct
	if err := c.Bind(&itemContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user. if err is not nil, print out the response struct
		response := models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	statement, err := db.Prepare(`DELETE FROM products WHERE product_id=$1 RETURNING product_id;`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return err
	}
	var items models.Item
	err = statement.QueryRow(&itemContainer.Product_id).Scan(&items.Product_id)
	if err != nil {
		fmt.Println("Exec Error:", err)
		response := models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := models.Response{
		Message: "SUCCESS deleted by" + name,
		Status:  "SUCCESS",
		Result:  items.Product_id,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
=================================

GET TRANSACTION (SELECT)

==================================
*/
func GetTransaction(c echo.Context) error {
	db := db.Conn()
	statement, err := db.Prepare(`SELECT id, transaction_id, product_id, transaction_date, user_id, payment_method, quantity FROM public.transaction_history;`)
	if err != nil {
		fmt.Println(err)
	}
	rows, err := statement.Query()
	if err != nil {
		fmt.Println(err)
	}
	var trsContainer []models.Transaction

	for rows.Next() {
		var trans models.Transaction
		err := rows.Scan(&trans.ID, &trans.Transaction_id, &trans.Product_id, &trans.Transaction_date, &trans.User_id, &trans.Payment_method, &trans.Quantity)
		if err != nil {
			fmt.Println(err)
		}
		trsContainer = append(trsContainer, trans)
	}
	res := models.Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  trsContainer,
	}
	return c.JSON(http.StatusOK, res)
}

/*========================================

ADD TRANSACTION

========================================*/

func AddTransaction(c echo.Context) error {
	db := db.Conn()
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	name := claims.Name                   //returns dummyuser instead of Dummy A?
	var transactions []models.Transaction // declare "user" as new User struct
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
	sqlStr = utils.ReplaceSQL(sqlStr, "?")
	//prepare the statement
	statement, _ := db.Prepare(sqlStr)

	//format all vals at once
	_, err := statement.Exec(vals...)
	if err != nil {
		fmt.Println("Bind Error:", err)
		return err
	}

	response := models.Response{
		Message: "SUCCESS added by " + name,
		Status:  "SUCCESS",
		Result:  vals,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
================================

# UPDATE TRANSACTION

====================================
*/
func UpdateTransaction(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	name := claims.Name
	db := db.Conn()
	var transaction models.Transaction // declare "transaction" as new Transaction struct
	if err := c.Bind(&transaction); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	statement, err := db.Prepare(`UPDATE 
	public.transaction_history 
	SET 
	transaction_id=$1, product_id=$2, transaction_date=$3, user_id=$4, payment_method=$5, 
	quantity=$6 
	WHERE id=$7 RETURNING id;`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return err
	}
	var transactions models.Transaction
	err = statement.QueryRow(transaction.Transaction_id, transaction.Product_id, transaction.Transaction_date, transaction.User_id, transaction.Payment_method, transaction.Quantity, transaction.ID).Scan(&transactions.ID)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return err
	}

	response := models.Response{
		Message: "SUCCESS updated by" + name,
		Status:  "SUCCESS",
		Result:  statement,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
================================

# DELETE TRANSACTION

====================================
*/
func DeleteTransaction(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	name := claims.Name
	db := db.Conn()
	var transContainer models.Transaction // declare "user" as new User struct
	if err := c.Bind(&transContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	statement, err := db.Prepare(`DELETE FROM public.transaction_history WHERE transaction_id=$1 RETURNING transaction_id`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return err
	}
	var transactions models.Transaction
	err = statement.QueryRow(&transContainer.Transaction_id).Scan(&transactions.Transaction_id)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return err
	}

	response := models.Response{
		Message: "SUCCESS updated by" + name,
		Status:  "SUCCESS",
		Result:  statement,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
=================================

GET USERS (SELECT)

==================================
*/
func GetUsers(c echo.Context) error {
	db := db.Conn()

	statement, err := db.Prepare(`SELECT id, age, first_name, last_name, email, username, "group" FROM public.users;`)
	if err != nil {
		fmt.Println(err)
	}
	rows, err := statement.Query()
	if err != nil {
		fmt.Println(err)
	}
	var userContainer []models.User
	for rows.Next() {
		var users models.User
		err := rows.Scan(&users.ID, &users.Age, &users.First_name, &users.Last_name, &users.Email, &users.Username, &users.Group)
		if err != nil {
			fmt.Println(err)
		}
		userContainer = append(userContainer, users)
	}
	res := models.Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  userContainer,
	}
	return c.JSON(http.StatusOK, res)
}

/*========================================

ADD users

========================================*/

func AddUsers(c echo.Context) error {
	db := db.Conn()
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	name := claims.Name     //returns dummyuser instead of Dummy A?
	var users []models.User // declare "user" as new User struct
	if err := c.Bind(&users); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		response := models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	vals := []interface{}{}
	sqlStr := `INSERT INTO public.users(age, first_name, last_name, email, username, password, "group") VALUES `

	for _, row := range users { //index,name_of_
		sqlStr += " (?, ?, ?, ?, ?, ?, ?),"
		vals = append(vals, row.Age, row.First_name, row.Last_name, row.Email, row.Username, row.Password, row.Group)
	}
	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	// replacing ? with $n for postgres
	sqlStr = utils.ReplaceSQL(sqlStr, "?")
	//prepare the statement
	statement, _ := db.Prepare(sqlStr)

	//format all vals at once
	_, err := statement.Exec(vals...)
	if err != nil {
		fmt.Println("exec Error:", err)
		response := models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := models.Response{
		Message: "SUCCESS added by " + name,
		Status:  "SUCCESS",
		Result:  users,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
================================

# UPDATE USERS

====================================
*/
func UpdateUsers(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	name := claims.Name
	db := db.Conn()
	var userContainer models.User // declare "users" as new User struct for binding
	if err := c.Bind(&userContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	statement, err := db.Prepare(`UPDATE 
									public.users 
									SET 
									age=$1, first_name=$2, last_name=$3, email=$4, username=$5, 
									password=$6, "group"=$7 
									WHERE 
									id=$8 
									RETURNING id;`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return err
	}
	var users models.User //for Scan() container
	err = statement.QueryRow(&userContainer.Age, &userContainer.First_name, &userContainer.Last_name, &userContainer.Email, &userContainer.Username, &userContainer.Password, &userContainer.Group, &userContainer.ID).Scan(&users.ID)
	if err != nil {
		fmt.Println("Exec Error:", err)
		response := models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := models.Response{
		Message: "SUCCESS updated by" + name,
		Status:  "SUCCESS",
		Result:  users.ID,
	}

	return c.JSON(http.StatusCreated, response)
}

/*
================================

# DELETE USERS

====================================
*/
func DeleteUsers(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.JwtCustomClaims)
	name := claims.Name
	db := db.Conn()
	var userContainer models.User // declare "user" as new User struct
	if err := c.Bind(&userContainer); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	statement, err := db.Prepare(`DELETE FROM public.users WHERE id=$1 RETURNING id`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return err
	}
	var users models.User
	err = statement.QueryRow(&userContainer.ID).Scan(&users.ID)
	if err != nil {
		fmt.Println("Exec Error:", err)
		response := models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := models.Response{
		Message: "SUCCESS updated by" + name,
		Status:  "SUCCESS",
		Result:  users.ID,
	}

	return c.JSON(http.StatusCreated, response)
}
