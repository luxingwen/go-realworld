package upload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/luxingwen/go-realworld/config"
)

func UploadAvatarRegister(router *gin.RouterGroup) {
	router.POST("/", UploadAvatarByFrom)
}

func UploadAvatarByFrom(c *gin.Context) {
	f, h, err := c.Request.FormFile("avatar")
	if err != nil {
		log.Println("111--->", err)
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	exp := filepath.Ext(h.Filename)
	fmt.Println("exp: ", exp)
	if exp != ".jpg" && exp != ".jpeg" && exp != ".png" {
		fmt.Println("非法文件,只能上传jpg jpeg png格式图片")
		c.JSON(400, gin.H{"err": "非法文件,只能上传jpg jpeg png格式图片"})
		return
	}
	upload := "apistatic/public/avatar"

	filename := time.Now().Format("20060102150405") + "_" + h.Filename
	fs, err := os.Create(upload + "/" + filename)
	if err != nil {
		fmt.Print("err2:", err)
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}
	io.Copy(fs, f)
	c.JSON(200, gin.H{"avatar": config.ServerConf.FilePath + "/file/public/avatar/" + filename})
}
