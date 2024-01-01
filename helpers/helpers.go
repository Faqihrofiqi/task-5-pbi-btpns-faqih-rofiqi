package helpers

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/database"
	makan "github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/helpers/jwt"
	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/models"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func IsValidJWT(token string) bool {

	return govalidator.IsUUIDv4(token)
}

func GetUserFromToken(c *gin.Context) (*models.User, error) {

	tokenInterface, exists := c.Get("token")
	if !exists {
		fmt.Println(tokenInterface)
		return nil, errors.New("token not found in context")
	}

	tokenString, ok := tokenInterface.(string)
	if !ok {
		return nil, errors.New("failed to convert token to string")
	}

	token, err := jwt.ParseWithClaims(tokenString, &makan.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return makan.JWT_KEY, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*makan.JWTClaims)
	if !ok {
		return nil, errors.New("failed to get jwt claims")
	}

	var user models.User
	if err := database.PG.Where("username = ?", claims.Username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
