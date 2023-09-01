package models

import "github.com/golang-jwt/jwt/v4"

//needs to be declared somewhere that can be accessed by all packages AND main since packages cannot import main
//program will complain about interface conversion jwt.Claims is *main.jwtCustomClaims, not *products.jwtCustomClaims

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

type User struct {
	ID         int    `json:"id"`
	First_name string `json:"firstname"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Age        int    `json:"age"`
	Last_name  string `json:"lastname"`
	Group      string `json:"group"`
	Email      string `json:"email"`
}
type Response struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Result  interface{} `json:"result"`
	Errors  interface{} `json:"error"`
}
type Item struct {
	Product_id  int    `json:"product_id"`
	Name        string `json:"item"`
	Description string `json:"description"`
	Price       int    `json:"price"`
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
type Cart struct {
	ID         int    `json:"id"`
	User_id    string `json:"Userid"`
	Product_id int    `json:"ProductID"`
	Quantity   int    `json:"Quantity"`
}
