package jwtTest

import "github.com/golang-jwt/jwt/v4"

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

//needs to be declared somewhere that can be accessed by all packages AND main since packages cannot import main
//program will complain about interface conversion jwt.Claims is *main.jwtCustomClaims, not *products.jwtCustomClaims
