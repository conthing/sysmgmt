package handlers

import (
	"log"
	"sysmgmt-next/dto"
	"sysmgmt-next/services"

	"github.com/gin-gonic/gin"
)

//FileUpload 文件上传
func FileUpload(c *gin.Context) {
	file, err := c.FormFile("file.zip")
	if err != nil {
		log.Println("文件上传错误", err)
		c.JSON(400, err)
		return
	}

	err = c.SaveUploadedFile(file, "/tmp/file.zip")
	if err != nil {
		log.Println("文件保存失败", err)
		c.JSON(400, err)
		return
	}

	err = services.UpdateService()
	if err != nil {
		log.Println("更新失败", err)
		c.JSON(400, err)
		return
	}

	var resp dto.FileInfo
	resp.Downloading = true
	c.JSON(200, resp)

}
