package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExtractTokenFromCookieMiddleware adalah middleware untuk mengambil token dari cookie
func ExtractTokenFromCookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token otentikasi dari cookie
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			fmt.Println(token)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			log.Fatal(err)
			return

		}

		// Setel token di context agar dapat diakses oleh handler selanjutnya
		c.Set("token", token)

		// Lanjutkan ke handler berikutnya
		c.Next()
	}
}
