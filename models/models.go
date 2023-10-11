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

type Response struct {
	UserID  int         `json:"userID"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Result  interface{} `json:"result"`
	Errors  interface{} `json:"error"`
}

//cart:{userid etc etc} payment_method:etcetc
