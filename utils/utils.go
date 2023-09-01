package utils

import (
	"echojson/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	refreshTokenCookieName = "refresh-token"
	accessTokenCookieName  = "access-token"
	jwtSecretKey           = "secret"
	jwtRefreshSecretKey    = "secret"
)

func GetJWTSecret() string {
	return jwtSecretKey
}
func GetJWTRefresh() string {
	return jwtRefreshSecretKey
}

func GenerateAccessToken(user models.User) (string, error) {
	return GenerateToken(user, []byte(GetJWTSecret()))
}

/*
============================

# GENERATE TOKEN

==============================
*/
func GenerateToken(user models.User, secret []byte) (string, error) {

	claims := &models.JwtCustomClaims{ //need to put the struct in a common file exportable by main AND Products or it will complaim
		Name:  user.First_name,
		Admin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTErrorChecker(err error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("login"))
}

/*
========================================

# REPLACE SQL

========================================
*/
func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func ExtractToken(c echo.Context) (tokenStr string, err error) {
	authHeader := c.Request().Header.Get("Authorization")

	// Split the header into parts ("Bearer" and token )
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" { //if length is not 2 or the first part isn't "bearer"
		return "error: unauthorized (did you forget to set the Bearer token?)", c.String(http.StatusUnauthorized, "Unauthorized")
	}
	// Get the JWT token from the header
	tokenStr = parts[1]
	return
}

func ExtractAccessClaims(tokenStr string) (username string, err bool) {
	hmacSecretString := GetJWTSecret()
	hmacSecret := []byte(hmacSecretString)
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println("JWT Claims: ", claims)

		username := claims["name"]
		// password := claims["password"]
		return fmt.Sprintf("%v", username), true
	} else {
		log.Printf("Invalid JWT Token")
		return username, false
	}
}
