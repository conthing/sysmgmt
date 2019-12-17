package handlers

import (
	"log"
	"os/exec"
	"sysmgmt-next/config"
	"sysmgmt-next/dto"

	"github.com/gin-gonic/gin"
)

//FileUpload 文件上传
func FileUpload(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Println(file.Filename)
	c.SaveUploadedFile(file, "/tmp/")

	shellPath := config.Conf.ShellPath + "unzipfile.sh"
	command := exec.Command(shellPath, "/tmp/"+file.Filename, "/tmp")
	out, _ := command.Output()
	log.Println("unzip " + file.Filename + " -d /ota/ output:" + string(out))

	var resp dto.FileInfo
	resp.Downloading = true
	c.JSON(200, resp)
}
