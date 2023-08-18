package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type User struct {
	First_name string `json:'firstname'`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

const (
	accessTokenCookieName = "access-token"
	jwtSecretKey          = "secret"
)

func GetJWTSecret() string {
	return jwtSecretKey
}

func GenerateAccessToken(user *User) (string, time.Time, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	return GenerateToken(user, expirationTime, []byte(GetJWTSecret()))
}
func GenerateToken(user *User, expirationTime time.Time, secret []byte) (string, time.Time, error) {

	claims := &jwtCustomClaims{
		user.First_name,
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", time.Now(), err
	}
	return tokenString, expirationTime, nil
}
func JWTErrorChecker(err error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
}
