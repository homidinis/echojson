package repository

import (
	"database/sql"
	"echojson/db"
	"echojson/models"
	"echojson/utils"
	"errors"
	"fmt"
	"strings"
)

func GetCart(id int, userid int) (cartArray []models.Cart, err error) {
	db := db.Conn()

	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	var data []interface{}
	query := "SELECT id, quantity, product_id, price, user_id FROM public.cart WHERE user_id=$1"
	data = append(data, userid)
	if id != 0 { //if id is presented
		query += "AND id=$2"    //append "where" to query
		data = append(data, id) //then append the id arg to the data interface. in the case of there being a lot of arguments for a lot of WHERE conditions
	}
	fmt.Println("getCart data (userid):")
	fmt.Println(data)
	rows, err := db.Query(query, data...) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	if err != nil {
		fmt.Println("err in getstock")
		fmt.Println(err)
		if err == sql.ErrNoRows {
			fmt.Println("rows kosong")
		}
	}

	for rows.Next() { //for every row result, run Scan then Append the result into the products struct
		var cartScan models.Cart
		rows.Scan(&cartScan.ID, &cartScan.Quantity, &cartScan.Product_id, &cartScan.Price, &cartScan.User_id)
		cartArray = append(cartArray, cartScan)
	}
	return
}

/*
========================================

# ADD PRODUCTS

========================================
*/
func GetStock(id int) (amount int) {
	db := db.Conn()
	query := "SELECT quantity FROM products WHERE product_id=$1"
	err := db.QueryRow(query, id).Scan(&amount) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	if err != nil {
		fmt.Println("error in getstock, repository")
		fmt.Println("where product_id:")
		fmt.Println(id)
		fmt.Println(err)
	}
	return
}

// 1. declare array of Item struct (items)
// 2. bind items to json input
// 3. declare vals as an array
// 4. loop the items array, append into vals each Name, Description, price from items
// 5. return null if OK, error if error
func AddCart(cartRequest models.RequestCart, userid int) (vals []interface{}, err error) {
	db := db.Conn()

	var qty int
	sqlStr := `INSERT INTO cart (user_id, product_id, quantity, price) VALUES `
	//if exists, UPDATE instead of INSERT
	for _, row := range cartRequest.Request { //index,name_of_
		sqlStr += "(?, ?, ?,?),"
		vals = append(vals, userid, row.Product_id, row.Quantity, row.Price)
	}
	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	// replacing ? with $n for postgres
	sqlStr = utils.ReplaceSQL(sqlStr, "?")
	//prepare the statement
	statement, _ := db.Prepare(sqlStr)

	//format all vals at once
	_, err = statement.Exec(vals...)
	fmt.Println(cartRequest.Request)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}

	for _, cart := range cartRequest.Request { //loop over each request
		query := "UPDATE products SET quantity=$1 WHERE product_id=$2;" //once added to cart, update quantity
		_, err := db.Exec(query, qty-cart.Quantity, cart.Product_id)    //set quantity as quantity (stock we acquired from db)
		fmt.Println("current stock, supposedly:")
		fmt.Println(qty - cart.Quantity)
		fmt.Println("Updated cart")
		if err != nil {
			fmt.Println(err)
			// Handle the error appropriately
		}
	}
	return
}
func TransactionDetailInsert(cart models.Cart) (err error) {
	db := db.Conn()

	query2 := "INSERT INTO public.transaction_detail(transaction_id,product_id, quantity, price) VALUES ($1, $2, $3,$4);"

	_, err = db.Exec(query2, utils.IncrementTrxID(), cart.Product_id, cart.Quantity, cart.Price) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	if err != nil {
		fmt.Println("error in trxdi queryrow:")
		fmt.Println(err)
	}
	fmt.Println("TransactionDetailInsert ran")
	fmt.Println("Getcart:")
	fmt.Println(cart)
	return
}
func TransactionHistoryInsert(cart models.Cart) (err error) {
	db := db.Conn()
	query3 := "INSERT INTO public.transaction_history(transaction_id,user_id) VALUES ($1,$2);"
	_, err = db.Exec(query3, utils.IncrementTrxID(), cart.User_id) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	if err != nil {
		fmt.Println("error in trxhi queryrow:")
		fmt.Println(err)
	}
	fmt.Println("TransactionHistoryInsert ran")
	return
}

/*
================================

# UPDATE PRODUCTS

====================================
*/
func UpdateCart(cartContainer models.Cart, user int) (updated_id int, err error) { //returns response
	db := db.Conn()

	statement, err := db.Prepare(`UPDATE public.cart SET product_id=$1, quantity=$2 WHERE id=$3 RETURNING product_id`)
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	var cartScan models.Cart
	err = statement.QueryRow(&cartContainer.Product_id, &cartContainer.Quantity, &cartContainer.ID).Scan(&cartScan.Product_id)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}
	updated_id = cartContainer.ID
	return
}

/*
================================

# DELETE PRODUCTS

====================================
*/
func DeleteCart(cartContainer models.Cart, user int) (cart_id string, err error) {

	db := db.Conn()

	statement, err := db.Prepare(`DELETE FROM cart WHERE product_id=$1 RETURNING id;`)
	if err != nil {
		fmt.Println("Prep Error in controller:", err)
		return
	}
	var items models.Item
	err = statement.QueryRow(&cartContainer.Product_id).Scan(&items.Product_id)
	if err != nil {
		fmt.Println("Exec Error in controller:", err)
		return
	}

	return
}
