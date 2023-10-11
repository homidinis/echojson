package repository

import (
	"database/sql"
	"echojson/config"
	"echojson/models"
	"echojson/utils"
	"fmt"
	"strings"
)

func GetProducts(id int) (products []models.Item, err error) {
	db := config.Conn()

	var data []interface{}
	query := "SELECT name, description, price, product_id, quantity FROM products"

	if id != 0 { //if id is not presented
		query += " WHERE product_id=$1" //append "where" to query
		data = append(data, id)         //then append the id arg to the data interface. in the case of there being a lot of arguments for a lot of WHERE conditions
	}
	fmt.Println(query)
	rows, err := db.Query(query, data...) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	fmt.Println("err in getproducts?: ", err)

	if err != nil {
		fmt.Println("error in repo get products", err.Error())
		return products, err
	}

	for rows.Next() { //for every row result, run Scan then Append the result into the products struct
		var product models.Item
		err := rows.Scan(&product.Name, &product.Description, &product.Price, &product.Product_id, &product.Quantity)
		if err != nil {
			fmt.Println(err)
		}
		products = append(products, product)
	} //products already declared as return value, instead of old method returning c.JSON(products)
	return
}

func GetProductsV2(id int) (products models.Item, err error) {
	db := config.Conn()

	query := "SELECT name, description, price, product_id, quantity FROM products WHERE product_id=$1"
	fmt.Println(query)
	err = db.QueryRow(query, id).Scan(&products.Name, &products.Description, &products.Price, &products.Product_id, &products.Quantity) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	fmt.Println("err in getproductsv2?: ", err)

	if err != nil {
		fmt.Println("error in repo get products", err.Error())
		if err == sql.ErrNoRows {
			fmt.Println("product not exists")
		}
		return products, err
	}
	return
}

/*========================================

ADD PRODUCTS

========================================*/

func AddProducts(items models.RequestItem, tx *sql.Tx) (vals []interface{}, err error) {
	//1. declare array of Item struct (items)
	//2. bind items to json input
	//3. declare vals as an array
	//4. loop the items array, append into vals each Name, Description, price from items
	//5. return null if OK, error if error
	sqlStr := `INSERT INTO products (name, description, price) VALUES `

	for _, row := range items.Request { //index,name_of_
		sqlStr += "(?, ?, ?),"
		vals = append(vals, row.Name, row.Description, row.Price)
	}
	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	// replacing ? with $n for postgres
	sqlStr = utils.ReplaceSQL(sqlStr, "?")
	//prepare the statement

	//format all vals at once
	_, err = tx.Exec(sqlStr, vals...)
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
func UpdateProducts(itemContainer models.Item, user int, tx *sql.Tx) (updated_id int, err error) { //returns response
	query := `UPDATE public.products SET name=$1, description=$2, price=$3 WHERE product_id=$4 RETURNING product_id;`
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	var items models.Item
	err = tx.QueryRow(query, &itemContainer.Name, &itemContainer.Description, &itemContainer.Price, &itemContainer.Product_id).Scan(&items.Product_id)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}
	updated_id = items.Product_id
	return
}

func UpdateProductsQuantity(qty int, product_id int) (err error) {
	db := config.Conn()

	query := "UPDATE public.products SET quantity=$1 WHERE product_id=$2"
	_, err = db.Exec(query, qty, product_id)
	if err != nil {
		fmt.Println("update products quantity error: ")
		fmt.Println(err.Error())
		return
	}
	return
}

/*
================================

# DELETE PRODUCTS

====================================
*/
func DeleteProducts(itemContainer models.Item, user int, tx *sql.Tx) (product_id string, err error) {

	query := `DELETE FROM products WHERE product_id=$1 RETURNING product_id;`
	if err != nil {
		fmt.Println("Prep Error in controller:", err)
		return
	}
	var items models.Item
	err = tx.QueryRow(query, &itemContainer.Product_id).Scan(&items.Product_id)
	if err != nil {
		fmt.Println("Deleteproducts Error in controller:", err)
		return
	}

	return
}
