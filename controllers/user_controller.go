package controllers

import (
	"encoding/json"
	"strings"

	"time"

	"net/http"

	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/database"
	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/helpers"
	makan "github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/helpers/jwt"
	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var userInput models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	var user models.User
	var queryCondition string

	if strings.Contains(userInput.Username, "@") {

		queryCondition = "email = ?"
	} else {

		queryCondition = "username = ?"
	}

	if err := database.PG.Where(queryCondition, userInput.Username).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": "Username or email not found"}
			helpers.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helpers.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{"message": "Wrong password"}
		helpers.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	expTime := time.Now().Add(5 * time.Minute)
	claims := &makan.JWTClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Faqihrofiqi",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenAlgo.SignedString(makan.JWT_KEY)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	response := map[string]string{"message": "Login successful"}
	helpers.ResponseJSON(w, http.StatusOK, response)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "Logout Sukses"}
	helpers.ResponseJSON(w, http.StatusOK, response)

}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashPassword)

	if err := database.PG.Create(&user).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "success"}
	helpers.ResponseJSON(w, http.StatusOK, response)
	database.PG.Save(&user)

}

func GetUsers(c *gin.Context) {
	var user []models.User
	database.PG.Find(&user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}
func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("userId")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if database.PG.Model(&user).Where("id = ?", id).Updates(&user).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Update failed"})
		return
	}
	database.PG.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Update success"})
}

func DeleteUserHandler(c *gin.Context) {
	var user models.User

	userID := c.Param("userId")

	result := database.PG.Where("id = ?", userID).Delete(&user)
	if result.Error != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}

	if result.RowsAffected == 0 {

		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
