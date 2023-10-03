package models

import (
	_ "github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

//needs to be declared somewhere that can be accessed by all packages AND main since packages cannot import main
//program will complain about interface conversion jwt.Claims is *main.jwtCustomClaims, not *products.jwtCustomClaims

type JwtCustomClaims struct {
	UserID int  `json:"userid"`
	Admin  bool `json:"admin"`
	jwt.RegisteredClaims
}

type User struct {
	ID         int    `json:"id"`
	First_name string `json:"firstName" validate:"required"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Age        int    `json:"age" validate:"required,gte=0,lte=80"`
	Last_name  string `json:"lastName" validate:"required"`
	Admin      *bool  `json:"isAdmin" `
	Email      string `json:"email" validate:"required,email"`
}
type RequestUser struct {
	Request []User `json:"request" validate:"required,dive"`
}
type RequestTransaction struct { //cuma untuk bulk insert/update
	Request []Transaction `json:"request" validate:"required,dive"`
}
type RequestItem struct {
	Request []Item `json:"request" validate:"required,dive"`
}
type RequestCart struct {
	Request []Cart `json:"request" validate:"required,dive"`
}
type Response struct {
	UserID  int         `json:"userID"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Result  interface{} `json:"result"`
	Errors  interface{} `json:"error"`
}
type Item struct {
	Product_id  int     `json:"product_id"`
	Name        string  `json:"item"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"min=0"`
	Quantity    int     `json:"quantity" validate:"min=0"`
}

type Transaction struct {
	Transaction_id   string  `json:"transactionID" `
	Transaction_date *string `json:"transactionDate" ` //pointer allows null values (one null value will turn everything blank!!)
	ID               int     `json:"id" `
	User_id          int     `json:"userID"  `
	Payment_method   string  `json:"paymentMethod" `
}
type Cart struct {
	ID         int `json:"id"`
	User_id    int `json:"userID" validate:"required"`
	Product_id int `json:"productID" validate:"required"`
	Quantity   int `json:"quantity" validate:"required,gte=0"`
	Price      int `json:"price" validate:"required"`
}
type PaymentMethodCart struct {
	Cart           Cart   `json:"cart"`
	Payment_method string `json:"payment_method"`
}

type UserID struct {
	User_id        int    `json:"userid"`
	Payment_method string `json:"payment_method"`
}

//cart:{userid etc etc} payment_method:etcetc
