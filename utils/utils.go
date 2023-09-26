package utils

import (
	"database/sql"
	"echojson/db"
	"echojson/models"
	"errors"
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

func GenerateAccessTokenAdmin(user models.User) (string, error) {
	return GenerateToken(user, []byte(GetJWTSecret()))
}
func GenerateAccessTokenUser(user models.User) (string, error) {
	return GenerateToken(user, []byte(GetJWTSecret()))
}

/*
============================

# GENERATE TOKEN

==============================
*/
func GenerateToken(user models.User, secret []byte) (string, error) {

	claims := &models.JwtCustomClaims{ //need to put the struct in a common file exportable by main AND Products or it will complaim
		UserID: user.ID,
		Admin:  user.Admin,
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
	if len(authHeader) == 0 {
		fmt.Println("authHeader is empty")
		return "", errors.New("token is empty")
	}

	// Split the header into parts ("Bearer" and token )
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" { //if length is not 2 or the first part isn't "bearer"
		return "error: unauthorized (did you forget to set the Bearer token?)", errors.New("Unauthorized")
	}
	// Get the JWT token from the header
	tokenStr = parts[1]
	return
}

func ExtractAccessClaims(tokenStr string) (username int, isAdmin bool, err error) {
	if len(tokenStr) == 0 {
		fmt.Println("tokenStr empty")
		return 0, false, errors.New("token is empty")
	}
	hmacSecretString := GetJWTSecret()
	hmacSecret := []byte(hmacSecretString)
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println("JWT Claims: ", claims)
		isAdmin = claims["admin"].(bool)
		if userID, exists := claims["userid"].(float64); exists {
			// Convert the float64 to int
			username = int(userID)
			fmt.Println(isAdmin)
			return username, isAdmin, err
		} else {
			log.Printf("UserID is not of type float64 in JWT Claims")
		}

		// password := claims["password"]
		return username, isAdmin, err
	} else {
		log.Printf("Invalid JWT Token")
		return username, isAdmin, err
	}
}

// func TranslateError(err error) (errs []error) {

//		english := en.New()
//		uni := ut.New(english, english)
//		trans, _ := uni.GetTranslator("en")
//		if err == nil {
//			return nil
//		}
//		validatorErrs := err.(validator.ValidationErrors)
//		for _, e := range validatorErrs {
//			translatedErr := fmt.Errorf(e.Translate(trans))
//			errs = append(errs, translatedErr)
//		}
//		return errs
//	}
func BindValidateStruct(ctx echo.Context, i interface{}) error {
	if err := ctx.Bind(i); err != nil {
		return err
	}

	if err := ctx.Validate(i); err != nil {
		return err
	}

	return nil
}

func IncrementTrxID() string {
	db := db.Conn()
	query := "SELECT transaction_id FROM transaction_history ORDER BY DESC LIMIT 1"
	row := db.QueryRow(query)

	var latestTransactionID string
	if latestTransactionID == "" {
		fmt.Println("Latest transaction_id is empty")
		return "T001"
	}

	err := row.Scan(&latestTransactionID)

	numericPart, err := strconv.Atoi(latestTransactionID[1:])
	if err != nil {
		fmt.Println("Error in parsing numeric part:", err)
		return ""
	}
	numericPart++

	return fmt.Sprintf("T%03d", numericPart)
}

func DBTransaction(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Rollback Panic
		} else if err != nil {
			tx.Rollback() // err is not nill
		} else {
			err = tx.Commit() // err is nil
		}
	}()
	err = txFunc(tx)
	return err
}
