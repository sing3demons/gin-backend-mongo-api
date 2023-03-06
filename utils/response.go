package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseJsonWithLogger(c *gin.Context, status int, data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.Writer.Header().Add("Response-Json", string(json))
	c.JSON(status, data)
}