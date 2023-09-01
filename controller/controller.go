package controller

import (
	"echojson/db"
	"echojson/models"
	"echojson/utils"
	"fmt"
	"strings"
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

/*
=================================

GET TRANSACTION (SELECT)

==================================
*/
func GetTransaction(id string) (transactions []models.Transaction, err error) {
	db := db.Conn()

	var data []interface{}
	query := "SELECT id, transaction_id, product_id, transaction_date, payment_method, quantity, user_id FROM public.transaction_history"

	if id != "" { //if id is not presented
		query += "WHERE transaction_id=$1" //append "where" to query
		data = append(data, id)            //then append the id arg to the data interface. in the case of there being a lot of arguments for a lot of WHERE conditions
	}

	fmt.Println(query)
	rows, err := db.Query(query, data...) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	if err != nil {
		fmt.Println(err)
	}
	var transaction models.Transaction
	for rows.Next() { //for every row result, run Scan then Append the result into the products struct
		err := rows.Scan(&transaction.ID, &transaction.Transaction_id, &transaction.Payment_method, &transaction.Product_id, &transaction.Transaction_date, &transaction.Quantity, &transaction.User_id)
		if err != nil {
			fmt.Println(err)
		}
		transactions = append(transactions, transaction)
	} //products already declared as return value, instead of old method returning c.JSON(products)
	return
}

/*
=================================

GET USERS (SELECT)
todo: add "if id=0"
==================================
*/
func GetUsers(id int) (userContainer []models.User) { //return userContainer, yang map ke response di usecase
	db := db.Conn()
	query := `SELECT id, age, first_name, last_name, email, username, "group" FROM public.users`
	var data []interface{}

	if id != 0 { //if id is not presented
		query += "WHERE id=$1"  //append "where" to query
		data = append(data, id) //then append the id arg to the data interface. in the case of there being a lot of arguments for a lot of WHERE conditions
	}

	fmt.Println("query: " + query)
	rows, err := db.Query(query, data...) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var users models.User
		err := rows.Scan(&users.ID, &users.Age, &users.First_name, &users.Last_name, &users.Email, &users.Username, &users.Group)
		if err != nil {
			fmt.Println(err)
		}
		userContainer = append(userContainer, users)
	}
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

/*========================================

ADD TRANSACTION

========================================*/

func AddTransaction(transactions []models.Transaction, user string) (err error) {
	db := db.Conn()
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
	_, err = statement.Exec(vals...)
	if err != nil {
		fmt.Println("QUERY ERROR:", err)
		return
	}
	return
}

/*========================================

ADD users

========================================*/

func AddUsers(users []models.User, user string) (vals []interface{}, err error) {
	db := db.Conn()
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
	_, err = statement.Exec(vals...)
	if err != nil {
		fmt.Println("exec Error:", err)
		return
	}
	return
}

/*========================================

ADD cart

========================================*/

func AddCart(cartScan []models.Cart, user string) (err error) {
	db := db.Conn()
	vals := []interface{}{}
	sqlStr := `INSERT INTO public.cart(id, user_id, product_id, quantity) VALUES `

	for _, row := range cartScan { //index,name_of_
		sqlStr += " (?, ?, ?, ?),"
		vals = append(vals, row.ID, row.User_id, row.Product_id, row.Quantity)
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
		fmt.Println("exec Error:", err)
		return err
	}
	return
}

/*
================================

# UPDATE PRODUCTS

====================================
*/
func UpdateProducts(itemContainer models.Item, user string) (updated_id int, err error) { //returns response
	db := db.Conn()

	statement, err := db.Prepare(`UPDATE public.products SET name=$1, description=$2, price=$3 WHERE product_id=$4 RETURNING product_id;`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	var items models.Item
	err = statement.QueryRow(&itemContainer.Name, &itemContainer.Description, &itemContainer.Price, &itemContainer.Product_id).Scan(&items.Product_id)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}
	updated_id = items.Product_id
	return
}

/*
================================

# UPDATE TRANSACTION

====================================
*/
func UpdateTransaction(transaction models.Transaction, user string) (updated_id int, err error) {
	db := db.Conn()

	statement, err := db.Prepare(`UPDATE 
	public.transaction_history 
	SET 
	transaction_id=$1, product_id=$2, transaction_date=$3,
	user_id=$4, payment_method=$5, quantity=$6 
	WHERE 
	id=$7 
	RETURNING id;`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	var transactions models.Transaction
	err = statement.QueryRow(transaction.Transaction_id,
		transaction.Product_id,
		transaction.Transaction_date,
		transaction.User_id,
		transaction.Payment_method,
		transaction.Quantity,
		transaction.ID).Scan(&transactions.ID)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}
	updated_id = transactions.ID
	return
}

/*
================================

# UPDATE USERS
this part updates users
====================================
*/
func UpdateUsers(userContainer models.User, user string) (updated_id int, err error) {

	db := db.Conn()

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
		return
	}
	var users models.User //for Scan() container
	err = statement.QueryRow(&userContainer.Age,
		&userContainer.First_name,
		&userContainer.Last_name,
		&userContainer.Email,
		&userContainer.Username,
		&userContainer.Password,
		&userContainer.Group,
		&userContainer.ID).Scan(&users.ID)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}
	updated_id = users.ID
	return
}

/*
================================

# UPDATE Cart
this part updates cart
====================================
*/
func UpdateCart(cartContainer models.Cart, user string) (cartScan models.Cart, err error) {

	db := db.Conn()

	statement, err := db.Prepare(`UPDATE public.cart
								SET id=?, user_id=?, product_id=?, quantity=?
								WHERE <condition>;`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	err = statement.QueryRow(&cartContainer.ID, &cartContainer.User_id, &cartContainer.Product_id, &cartContainer.Quantity).Scan(&cartScan.ID)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}
	return
}

/*
================================

# DELETE TRANSACTION

====================================
*/
func DeleteTransaction(transContainer models.Transaction, user string) (transaction_id string, err error) {
	db := db.Conn()

	statement, err := db.Prepare(`DELETE FROM public.transaction_history WHERE transaction_id=$1 RETURNING transaction_id`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	var transactions models.Transaction
	err = statement.QueryRow(&transContainer.Transaction_id).Scan(&transactions.Transaction_id)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}
	transaction_id = transactions.Transaction_id
	return
}

/*
================================

# DELETE PRODUCTS

====================================
*/
func DeleteProducts(itemContainer models.Item, user string) (product_id string, err error) {

	db := db.Conn()

	statement, err := db.Prepare(`DELETE FROM products WHERE product_id=$1 RETURNING product_id;`)
	if err != nil {
		fmt.Println("Prep Error in controller:", err)
		return
	}
	var items models.Item
	err = statement.QueryRow(&itemContainer.Product_id).Scan(&items.Product_id)
	if err != nil {
		fmt.Println("Exec Error in controller:", err)
		return
	}

	return
}

/*
================================

# DELETE USERS

====================================
*/
func DeleteUsers(userContainer models.User, user string) (users models.User, err error) {

	db := db.Conn()

	statement, err := db.Prepare(`DELETE FROM public.users WHERE id=$1 RETURNING id`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	err = statement.QueryRow(&userContainer.ID).Scan(&users.ID)
	if err != nil {
		fmt.Println("Exec Error:", err) //response pindah ke usecase, jadi di controller cuma get value dari query;
		return
	}
	return
}

/*
================================

# DELETE CART

====================================
*/
func DeleteCart(userContainer models.User, user string) (users models.User, err error) {

	db := db.Conn()

	statement, err := db.Prepare(`DELETE FROM public.users WHERE id=$1 RETURNING id`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	err = statement.QueryRow(&userContainer.ID).Scan(&users.ID)
	if err != nil {
		fmt.Println("Exec Error:", err)

		return
	}
	return
}
