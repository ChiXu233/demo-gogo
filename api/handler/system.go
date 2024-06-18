package handler

import (
	"demo-gogo/config"
	"github.com/gin-gonic/gin"
)

func (handler *RestHandler) V1Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    200,
		"message": config.Conf.DB.Name,
	})
}

func (handler *RestHandler) V2Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    200,
		"message": "ok",
	})
}
