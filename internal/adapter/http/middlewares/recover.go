package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch t := err.(type) {
				case error:
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": t,
					})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{})
				}
			}
		}()

		c.Next()
	}
}
