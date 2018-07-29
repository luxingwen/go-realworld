package common

import (
	"github.com/gin-gonic/gin"
)

func HandleOk(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"code": 0, "data": data, "msg": "ok"})
}

func HandleErr(c *gin.Context, code int, msg string) {
	c.JSON(200, gin.H{"code": code, "msg": msg})
}
