package usecase

import (
	"database/sql"
	"echojson/db"
	"echojson/models"
	"echojson/repository"
	"echojson/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

/*
====================

# CHECKOUT
Get cart, increment trx id, then loop through every entry in the cart to insert them into TRX history and TRX detail
=======================
*/
func Checkout(c echo.Context) error {
	// tokenStr, err := utils.ExtractToken(c)
	// if err != nil {
	// 	response := models.Response{
	// 		Message: "ExtractToken MISSING TOKEN!",
	// 		Status:  "ERROR",
	// 		Result:  nil,
	// 		Errors:  err.Error(),
	// 	}
	// 	return c.JSON(http.StatusInternalServerError, response)
	// }

	// user, _, err := utils.ExtractAccessClaims(tokenStr)
	// if err != nil {
	// 	result := models.Response{
	// 		Message: "error in extractaccessclaims",
	// 		Status:  "ERROR",
	// 		Result:  nil,
	// 		Errors:  err.Error(),
	// 	}
	// 	return c.JSON(http.StatusInternalServerError, result)
	// }

	var cartScan models.UserID
	err := utils.BindValidateStruct(c, &cartScan)
	if err != nil {
		result := models.Response{
			Message: "ERROR IN BIND VALIDATE STRUCT",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, result)
	}

	carts, err := repository.GetCart(0, cartScan.User_id) //product id,user id
	if len(carts) == 0 {
		result := models.Response{
			Message: "ERROR CART EMPTY",
			Status:  "ERROR",
			Result:  nil,
			Errors:  "cart is empty",
		}
		return c.JSON(http.StatusInternalServerError, result)
	}

	var trxID string
	if trxID, err = utils.IncrementTrxID(); err != nil {
		fmt.Println("incrementing error in usecase line 57")
		return err
	}

	for _, cart := range carts {
		qty := repository.GetStock(cart.Product_id)

		if qty-cart.Quantity < 0 {
			fmt.Println("quantity bigger than stock")
			return err
		} else {
			err = utils.DBTransaction(db.Conn(), func(tx *sql.Tx) (err error) {
<<<<<<< HEAD
				cartReq := models.PaymentMethodCart{ //separate this into its own function (process cart)
					Cart:           cart,
					Payment_method: cartScan.Payment_method,
				}
				err = repository.TransactionHistoryInsert(cartReq, trxID, tx)
=======
				err = repository.TransactionDetailInsert(cart, trxID, tx)
>>>>>>> f8d3a2cd046eb17d2f021dc1cdabd26781c4c4b5
				if err != nil {
					fmt.Println("error in inserting transaction detail:")
					fmt.Println(err)
				}

				return err
			})

			result := models.Response{
				Message: "SUCCESS",
				Status:  "SUCCESS",
				Result:  err,
				Errors:  nil,
			}
			return c.JSON(http.StatusOK, result)
			for _, cart := range carts {
				err = utils.DBTransaction(db.Conn(), func(tx *sql.Tx) (err error) {
					fmt.Println("trx history inserted!")
					pmcart := models.PaymentMethodCart{
						Cart:           cart,
						Payment_method: cartScan.Payment_method,
					}
					err = repository.TransactionHistoryInsert(pmcart, trxID, tx)
					if err != nil {
						return err
					}
					return
				})
			}

		}

	}

	return err
}

/*
====================

# GETCART

=======================
*/
func GetCart(c echo.Context) error {

	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		response := models.Response{
			Message: "ExtractToken MISSING TOKEN!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	user, _, err := utils.ExtractAccessClaims(tokenStr)
	if err != nil {
		response := models.Response{
			Message: "ExtractAccessClaims error!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	var cart models.Cart
	if err := c.Bind(&cart); err != nil {
		fmt.Println("Bind Error:", err) //if err is nil, bind user
		return err
	}
	result, err := repository.GetCart(cart.ID, user)
	if err != nil {
		response := models.Response{
			Message: "ERROR IN GETCART, LINE 104",
			Status:  "ERROR",
			Result:  result,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	response := models.Response{
		UserID:  user,
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  result,
		Errors:  nil,
	}
	return c.JSON(http.StatusOK, response)
}

/*
====================

# GETCART

=======================
*/
func InsertCart(c echo.Context) error {
	var cart models.RequestCart
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		response := models.Response{
			Message: "ExtractToken MISSING TOKEN!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	user, _, err := utils.ExtractAccessClaims(tokenStr)
	if err != nil {
		return err
	}
	// var cart = new(models.Requestcart)
	err = utils.BindValidateStruct(c, &cart)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	for _, cartReq := range cart.Request {
		products, err := repository.GetProductsV2(cartReq.Product_id)
		if err != nil {
			fmt.Println("err in getproducts, insert product:")
			result := models.Response{
				Message: "ERROR IN GETPRODUCTS",
				Status:  "error",
				Result:  nil,
				Errors:  err.Error(),
			}
			fmt.Print(err)
			return c.JSON(http.StatusInternalServerError, result) //need to return json to stop
		} else if cartReq.Quantity > products.Quantity {
			fmt.Println("quantity empty! err: ")
			fmt.Print(err)
			result := models.Response{
				Message: "ERROR IN GETPRODUCTS",
				Status:  "error",
				Result:  nil,
				Errors:  "qty bigger than stock",
			}
			return c.JSON(http.StatusInternalServerError, result)
		}
	}
	err = utils.DBTransaction(db.Conn(), func(tx *sql.Tx) (err error) {
		result, err := repository.AddCart(cart, user, tx)
		if err != nil {
			response := models.Response{
				Message: "ERROR in AddCart calling",
				Status:  "ERROR",
				Result:  result,
				Errors:  err.Error(),
			}
			return c.JSON(http.StatusInternalServerError, response)
		}
		response := models.Response{
			Message: "SUCCESS",
			Status:  "SUCCESS",
			Result:  result,
			Errors:  nil,
		}
		return c.JSON(http.StatusOK, response)
	})
	return err
}

/*
====================

# UPDATE

=======================
*/
func UpdateCart(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		response := models.Response{
			Message: "ExtractToken MISSING TOKEN!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	user, _, err := utils.ExtractAccessClaims(tokenStr)
	if err != nil {
		response := models.Response{
			Message: "ExtractAccessClaims error in update cart",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	var cartContainer models.Cart // declare "users" as new User struct for binding
	err = utils.BindValidateStruct(c, &cartContainer)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	err = utils.DBTransaction(db.Conn(), func(tx *sql.Tx) (err error) {
		result, err := repository.UpdateCart(cartContainer, user, tx)
		if err != nil {
			response := models.Response{
				Message: "ERROR",
				Status:  "ERROR",
				Result:  nil,
				Errors:  err.Error(),
			}
			return c.JSON(http.StatusInternalServerError, response)
		}
		response := models.Response{
			Message: "SUCCESS",
			Status:  "SUCCESS",
			Result:  result,
			Errors:  nil,
		}
		return c.JSON(http.StatusOK, response)
	})
	return err
}

/*
====================

# DELETE

=======================
*/

func DeleteCart(c echo.Context) error {
	tokenStr, err := utils.ExtractToken(c)
	if err != nil {
		response := models.Response{
			Message: "ExtractToken MISSING TOKEN!",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	user, _, err := utils.ExtractAccessClaims(tokenStr)
	if err != nil {
		response := models.Response{
			Message: "ERROR in extractaccessclaims, delete cart",
			Status:  "ERROR in binding",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	var cart models.Cart
	err = utils.BindValidateStruct(c, &cart)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result, err := repository.DeleteCart(cart, user)
	if err != nil {
		response := models.Response{
			Message: "ERROR",
			Status:  "ERROR",
			Result:  result,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	response := models.Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  result,
		Errors:  nil,
	}

	return c.JSON(http.StatusOK, response)
}
