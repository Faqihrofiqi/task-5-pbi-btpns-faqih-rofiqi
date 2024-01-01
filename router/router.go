package router

import (
	"net/http"

	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/controllers"
	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(c *gin.Context) {
	controllers.Register(c.Writer, c.Request)
}
func LoginHandlers(c *gin.Context) {
	controllers.LoginUser(c.Writer, c.Request)
}
func LogoutHandlers(c *gin.Context) {
	controllers.Logout(c.Writer, c.Request)
}

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/users", controllers.GetUsers)
	r.POST("/users/register", RegisterHandlers)
	r.PUT("/users/:userId", controllers.UpdateUser)
	r.DELETE("/users/:userId", controllers.DeleteUserHandler)
	r.POST("/users/login", LoginHandlers)
	r.GET("/users/logout", LogoutHandlers)

	r.Use(middlewares.ExtractTokenFromCookieMiddleware())
	{

		r.POST("/photos", controllers.UploadProfilePhoto)
		r.GET("/photos", controllers.GetPhotos)
		r.PUT("/photos/:photoId", controllers.UpdatePhotoProfile)
		r.DELETE("/photos/:photoId", controllers.DeletePhoto)
	}
	return r

}
