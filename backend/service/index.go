package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetIndex(r *gin.Context) {
	r.JSON(http.StatusOK, gin.H{
		"message": "welcome",
	})
}
