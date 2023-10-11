package repository

import (
	"database/sql"
	"echojson/config"
	"echojson/models"
	"echojson/utils"
	"errors"
	"fmt"
	"strings"
)

func GetCart(id int, userid int) (cartArray []models.Cart, err error) {
	db := config.Conn()

	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	var data []interface{}
	query := "SELECT id, quantity, product_id, price, user_id FROM public.cart WHERE user_id=$1 "
	data = append(data, userid)
	if id != 0 { //if id is presented
		query += "AND product_id=$2" //append "where" to query
		data = append(data, id)      //then append the id arg to the data interface. in the case of there being a lot of arguments for a lot of WHERE conditions
	}
	fmt.Println("getCart data (userid):")
	fmt.Println(data)
	fmt.Println(query)
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
		fmt.Println(cartArray)
	}
	return
}

/*
========================================

# ADD PRODUCTS

========================================
*/
func GetStock(id int) (amount int) {
	db := config.Conn()
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
func AddCart(cartRequest models.RequestCart, userid int, tx *sql.Tx) (vals []interface{}, err error) {
	getCart, err := GetCart(cartRequest.Request[0].Product_id, userid)
	if err != nil {
		return nil, err
	}
	fmt.Println("cartrequest.request[0].product_id: ")
	fmt.Println(cartRequest.Request[0].Product_id)
	fmt.Println("len of getcart: ")
	fmt.Println(len(getCart))

	if len(getCart) > 0 {

		for _, cart := range cartRequest.Request { //loop through every request given
			// Check if the product exists in the cart
			var existingQuantity int

			err := tx.QueryRow("SELECT quantity FROM public.cart WHERE product_id = $1", cart.Product_id).Scan(&existingQuantity)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
			sqlUpdate := "UPDATE public.cart SET quantity = $1 WHERE product_id = $2;"
			_, err = tx.Exec(sqlUpdate, existingQuantity+cart.Quantity, cart.Product_id)
			if err != nil {
				return nil, err
			}
		}
	} else {
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

		//format all vals at once
		_, err = tx.Exec(sqlStr, vals...)
		fmt.Println(cartRequest.Request)
		if err != nil {
			fmt.Println("Exec Error:", err)
			return
		}
	}
	return
}
func TransactionDetailInsert(cart models.Cart, trx_id string, tx *sql.Tx) (err error) {

	query2 := "INSERT INTO public.transaction_detail(product_id,transaction_id,quantity, price) VALUES ($1, $2, $3,$4);"
	// trxID, err := utils.IncrementTrxID()
	if err != nil {
		fmt.Println("Transaction ID error")
		return err
	}
	_, err = tx.Exec(query2, cart.Product_id, trx_id, cart.Quantity, cart.Price) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	if err != nil {
		fmt.Println("error in trxdi queryrow:")
		fmt.Println(err)
	}
	fmt.Println("TransactionDetailInsert ran")
	return
}
func TransactionHistoryInsert(cart models.PaymentMethodCart, trx_id string, tx *sql.Tx) (err error) {
	// trxID, err := utils.IncrementTrxID()
	if err != nil {
		fmt.Println("Transaction ID error")
		return err
	}
	query3 := "INSERT INTO public.transaction_history(transaction_id,user_id,payment_method) VALUES ($1,$2,$3);"
	_, err = tx.Exec(query3, trx_id, cart.Cart.User_id, cart.Payment_method) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
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
func UpdateCart(cartContainer models.Cart, user int, tx *sql.Tx) (updated_id int, err error) { //returns response

	query := `UPDATE public.cart SET product_id=$1, quantity=$2 WHERE id=$3 RETURNING product_id`
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	var cartScan models.Cart
	err = tx.QueryRow(query, &cartContainer.Product_id, &cartContainer.Quantity, &cartContainer.ID).Scan(&cartScan.Product_id)
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

	db := config.Conn()

	statement, err := db.Prepare(`DELETE FROM cart WHERE product_id=$1 AND user_id = $2 RETURNING id;`)
	if err != nil {
		fmt.Println("Prep Error in controller:", err)
		return
	}
	var items models.Item
	err = statement.QueryRow(&cartContainer.Product_id, &cartContainer.User_id).Scan(&items.Product_id)
	if err != nil {
		fmt.Println("Delete Error in controller:", err)
		return
	}

	return
}

// 1. loop:
// 2. insert into transaction detail (LOOP THROUGH CARTS AGAIN)
//stop loop
// 3. insert into transaction history (separate from processcart?)
// 4. delete cart content
