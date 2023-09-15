package usecase

import (
	"echojson/models"
	"echojson/repository"
	"echojson/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Checkout(c echo.Context) error {
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
		result := models.Response{
			Message: "error in extractaccessclaims",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, result)
	}

	var cartScan models.UserID
	err = utils.BindValidateStruct(c, &cartScan)
	if err != nil {
		result := models.Response{
			Message: "ERROR IN BIND VALIDATE STRUCT",
			Status:  "ERROR",
			Result:  nil,
			Errors:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, result)
	}

	var qty int
	carts, err := repository.GetCart(0, cartScan.User_id) //product id,user id
	fmt.Println("carts:")
	fmt.Println(carts)

	for _, cart := range carts {
		qty = repository.GetStock(cart.Product_id)
		if err != nil {
			fmt.Println("error in getstock")
			fmt.Println(cart.Product_id)
			return err
		}
		fmt.Println("Stock acquired:")
		fmt.Println(qty)
		if cart.Quantity > qty {
			result := models.Response{
				Message: "ERROR QTY EXCEEDS STOCK",
				Status:  "ERROR",
				Result:  nil,
				Errors:  errors.New("qty exceeds stock"),
			}
			fmt.Println(result)
			return c.JSON(http.StatusInternalServerError, result)
		} else if qty == 0 {
			fmt.Println("quantity 0")
			result := models.Response{
				Message: "ERROR QTY EMPTY",
				Status:  "ERROR",
				Result:  nil,
				Errors:  errors.New("qty empty"),
			}
			return c.JSON(http.StatusInternalServerError, result)
		} else {
			fmt.Println("if check passed")
			err = repository.TransactionDetailInsert(cart)
			if err != nil {
				fmt.Println(err)
			}
			err = repository.TransactionHistoryInsert(cart)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("id:")
			fmt.Println(cart.Product_id)
			_, err = repository.DeleteCart(cart, user)

			fmt.Println("cart delete:")
			fmt.Println(cart)
		}
	}

	result := models.Response{
		Message: "SUCCESS",
		Status:  "SUCCESS",
		Result:  err,
		Errors:  nil,
	}
	return c.JSON(http.StatusOK, result)
	//struct: {}
}

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
		Result:  nil,
		Errors:  nil,
	}
	return c.JSON(http.StatusOK, response)
}

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
		products, err := repository.GetProducts(cartReq.Product_id)
		if err != nil {
			fmt.Println("err in getproducts, insert product:")
			fmt.Print(err)
			return c.JSON(http.StatusInternalServerError, errors.New("product does not exist")) //need to return json to stop
		} else if cartReq.Quantity > products[0].Quantity {
			fmt.Println("quantity empty!")
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, errors.New("qty empty!"))
		}
	}
	result, err := repository.AddCart(cart, user)
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
}
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

	var cartContainer models.Cart // declare "users" as new User struct for binding
	err = utils.BindValidateStruct(c, &cartContainer)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	result, err := repository.UpdateCart(cartContainer, user)
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
}

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

	var cart models.Cart
	err = utils.BindValidateStruct(c, &cart)
	if err != nil {
		response := models.Response{
			Message: "ERROR in binding",
			Status:  "ERROR in binding",
			Result:  err.Error(),
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
