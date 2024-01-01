package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/database"
	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/helpers"
	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadProfilePhoto(c *gin.Context) {

	user, err := helpers.GetUserFromToken(c)
	if err != nil {
		helpers.ResponseJSON(c.Writer, http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	existingProfilePhoto := models.Photo{}
	if err := database.PG.Where("user_id = ? AND is_profile_photo = ?", user.ID, true).First(&existingProfilePhoto).Error; err == nil {

		helpers.ResponseJSON(c.Writer, http.StatusConflict, gin.H{"error": "User already has a profile photo"})
		return
	}

	var photoInput models.Photo
	photoInput.Title = c.PostForm("title")
	photoInput.Caption = c.PostForm("caption")
	photoInput.PhotoUrl = c.PostForm("photourl")

	if err := validatePhotoInput(photoInput); err != nil {
		helpers.ResponseJSON(c.Writer, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photoInput.UserID = user.ID
	photoInput.IsProfilePhoto = true

	if err := database.PG.Create(&photoInput).Error; err != nil {
		helpers.ResponseJSON(c.Writer, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = saveProfilePhoto(user.ID, photoInput.ID, c.Request)
	if err != nil {

		helpers.ResponseJSON(c.Writer, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	helpers.ResponseJSON(c.Writer, http.StatusCreated, gin.H{"message": "Profile photo added successfully"})
}

func validatePhotoInput(photo models.Photo) error {
	if photo.Title == "" {
		return errors.New("title cannot be empty")
	}

	return nil
}

func saveProfilePhoto(userID uint, photoID uint, r *http.Request) (string, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	storageDir := "./uploads"
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		err := os.Mkdir(storageDir, 0755)
		if err != nil {
			return "", err
		}
	}

	fileName := fmt.Sprintf("%d_%d_%s", userID, photoID, header.Filename)
	filePath := filepath.Join(storageDir, fileName)

	newFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return "", err
	}

	photoURL := getPhotoURL(filePath)

	if err := updatePhotoURL(photoID, photoURL); err != nil {
		return "", err
	}

	return photoURL, nil
}

func getPhotoURL(filePath string) string {
	baseURL := "public/photos/"
	return baseURL + "uploads/" + filepath.Base(filePath)
}

func updatePhotoURL(photoID uint, photoURL string) error {

	if err := database.PG.Model(&models.Photo{}).Where("id = ?", photoID).Update("photo_url", photoURL).Error; err != nil {
		return err
	}

	return nil
}

type PhotoResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Caption   string `json:"caption"`
	PhotoURL  string `json:"photo_url"`
	IsProfile bool   `json:"is_profile_photo"`
}

func GetPhotos(c *gin.Context) {
	user, err := helpers.GetUserFromToken(c)
	if err != nil {
		helpers.ResponseJSON(c.Writer, http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var photos []models.Photo
	if err := database.PG.Where("user_id = ?", user.ID).Find(&photos).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting user photos"})
		return
	}

	var photoResponses []PhotoResponse

	for _, photo := range photos {
		photoResponse := PhotoResponse{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoURL:  photo.PhotoUrl,
			IsProfile: photo.IsProfilePhoto,
		}
		photoResponses = append(photoResponses, photoResponse)
	}

	c.JSON(http.StatusOK, gin.H{"photos": photoResponses})
}

func DeletePhoto(c *gin.Context) {

	user, err := helpers.GetUserFromToken(c)
	if err != nil {
		helpers.ResponseJSON(c.Writer, http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	photoID := c.Param("photoId")

	var photo models.Photo
	result := database.PG.Model(&photo).Select("user_id").Where("id = ?", photoID).Scan(&photo.UserID)
	if result.Error != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user ID"})
		return
	}

	if user.ID != photo.UserID {

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You are not the owner of this photo"})
		return
	}

	result = database.PG.Where("id = ?", photoID).Delete(&photo)
	if result.Error != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting photo"})
		return
	}

	if result.RowsAffected == 0 {

		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	if err := database.PG.Model(&models.Photo{}).Where("id = ?", photoID).Update("IsProfilePhoto", false).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating isProfilePhoto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}

func UpdatePhotoProfile(c *gin.Context) {

	user, err := helpers.GetUserFromToken(c)
	if err != nil {
		helpers.ResponseJSON(c.Writer, http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	photoID := c.Param("photoId")

	var photo models.Photo
	if err := database.PG.First(&photo, photoID).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:

			helpers.ResponseJSON(c.Writer, http.StatusNotFound, gin.H{"error": "Photo not found"})
			return
		default:

			helpers.ResponseJSON(c.Writer, http.StatusInternalServerError, gin.H{"error": "Error querying photo"})
			return
		}
	}

	if photo.UserID != user.ID {
		helpers.ResponseJSON(c.Writer, http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var updatedPhoto models.Photo
	if err := c.ShouldBind(&updatedPhoto); err != nil {
		helpers.ResponseJSON(c.Writer, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validatePhotoInput(updatedPhoto); err != nil {
		helpers.ResponseJSON(c.Writer, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if file, header, err := c.Request.FormFile("file"); err == nil {
		defer file.Close()

		storageDir := "./uploads/"
		if _, err := helpers.CreateDirIfNotExist(storageDir); err != nil {
			helpers.ResponseJSON(c.Writer, http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fileName := fmt.Sprintf("%d_%d_%s", user.ID, photo.ID, header.Filename)
		filePath := filepath.Join(storageDir, fileName)

		newFile, err := helpers.CreateFile(filePath)
		if err != nil {
			helpers.ResponseJSON(c.Writer, http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer newFile.Close()

		if _, err := helpers.CopyFile(newFile, file); err != nil {
			helpers.ResponseJSON(c.Writer, http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		photo.PhotoUrl = getPhotoURL(filePath)
	} else if err != http.ErrMissingFile {

		helpers.ResponseJSON(c.Writer, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	photo.Title = updatedPhoto.Title
	photo.Caption = updatedPhoto.Caption

	if err := database.PG.Save(&photo).Error; err != nil {
		helpers.ResponseJSON(c.Writer, http.StatusInternalServerError, gin.H{"error": "Error updating photo"})
		return
	}

	helpers.ResponseJSON(c.Writer, http.StatusOK, gin.H{"message": "Photo updated successfully"})
}
