package jwt

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func GenerateJWT(id uint) (string, error) {
	tokenTTL, err := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}

func ValidateJWT(c *gin.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token")
}

func CurrentUserID(c *gin.Context) (uint, error) {
	err := ValidateJWT(c)
	if err != nil {
		return 0, err
	}

	token, err := getToken(c)
	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))
	return userID, nil
}

func getToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequests(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return privateKey, nil
	})
	return token, err
}

func getTokenFromRequests(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
