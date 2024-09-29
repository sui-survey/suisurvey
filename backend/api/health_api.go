package api

import (
	"github.com/gin-gonic/gin"
	"kevinsheeran/walrus-backend/model"
	"kevinsheeran/walrus-backend/result"
)

// HealthCheck
// @Summary HealthCheck API
// @Tags HealthCheck
// @Produce json
// @Success 200 {object} result.Result
// @router /health [get]
func HealthCheck(c *gin.Context) {

	var response model.WalrusResponse
	result.Success(c, response)

}
