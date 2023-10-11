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

/*
=================================

GET TRANSACTION (SELECT)

==================================
*/
func GetTransaction(id string, userid int) (transactions []models.Transaction, err error) {
	db := config.Conn()

	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	var data []interface{}
	query := "SELECT transaction_id, transaction_date, id, user_id, payment_method FROM public.transaction_history WHERE user_id=$1 "
	data = append(data, userid)
	if id != "" { //if id is presented
		query += "AND transaction_id=$2" //append "where" to query
		data = append(data, id)          //then append the id arg to the data interface. in the case of there being a lot of arguments for a lot of WHERE conditions
	}

	fmt.Println(query)
	fmt.Println(data)
	rows, err := db.Query(query, data...) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			fmt.Println("rows kosong")
		}
	}

	for rows.Next() { //for every row result, run Scan then Append the result into the products struct
		var trx models.Transaction
		rows.Scan(&trx.Transaction_id, &trx.Transaction_date, &trx.ID, &trx.User_id, &trx.Payment_method)
		transactions = append(transactions, trx)
		fmt.Println(transactions)
	}
	fmt.Println(transactions) //products already declared as return value, instead of old method returning c.JSON(products)
	return
}

/*========================================

ADD TRANSACTION

========================================*/

func AddTransaction(transactions models.RequestTransaction, user int, tx *sql.Tx) (vals []interface{}, err error) {
	sqlStr := `INSERT INTO public.transaction_history(transaction_id, transaction_date, id, user_id, payment_method)VALUES`
	for _, row := range transactions.Request { //index,name_of_ ; for every data inputted in, run loop use Request for bulk inserts
		sqlStr += " (?, ?, ?, ?, ?),"
		vals = append(vals, row.Transaction_id, row.Transaction_date, row.User_id, row.Payment_method)
	}
	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	// replacing ? with $n for postgres
	sqlStr = utils.ReplaceSQL(sqlStr, "?")
	//format all vals at once
	_, err = tx.Exec(sqlStr, vals...)

	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}

	// Check if the number of placeholders matches the number of values
	if len(vals) != strings.Count(sqlStr, "?") {
		fmt.Println("Number of placeholders does not match the number of values.")
		return
	}
	return
}

/*
================================

# UPDATE TRANSACTION

====================================
*/
func UpdateTransaction(trx models.Transaction, user int, tx *sql.Tx) (updated_id int, err error) {
	query := `UPDATE 
	public.transaction_history 
	SET 
	transaction_id=$1, product_id=$2, transaction_date=$3,
	user_id=$4, payment_method=$5, quantity=$6 
	WHERE 
	id=$7 
	RETURNING id;`
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	var transactions models.Transaction
	err = tx.QueryRow(query, trx.Transaction_id,
		trx.Transaction_date,
		trx.User_id,
		trx.Payment_method,
		trx.ID).Scan(&transactions.ID)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}
	updated_id = transactions.ID
	return
}

/*
================================

# DELETE TRANSACTION

====================================
*/
func DeleteTransaction(transContainer models.Transaction, user int, tx *sql.Tx) (transaction_id string, err error) {
	query := `DELETE FROM public.transaction_history WHERE transaction_id=$1 RETURNING transaction_id`
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	var transactions models.Transaction //output
	err = tx.QueryRow(query, &transContainer.Transaction_id).Scan(&transactions.Transaction_id)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}
	transaction_id = transactions.Transaction_id
	return
}
