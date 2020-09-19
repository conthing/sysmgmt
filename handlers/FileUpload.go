package handlers

import (
	"github.com/conthing/sysmgmt/dto"
	"github.com/conthing/sysmgmt/services"

	"github.com/conthing/utils/common"
	"github.com/gin-gonic/gin"
)

// todo 返回故障码混乱
//FileUpload 文件上传
func FileUpload(c *gin.Context) {
	file, err := c.FormFile("file.zip")
	if err != nil {
		common.Log.Errorf("Form file failed %v", err)
		c.JSON(400, err)
		return
	}

	err = c.SaveUploadedFile(file, "/tmp/file.zip")
	if err != nil {
		common.Log.Errorf("Save file failed %v", err)
		c.JSON(400, err)
		return
	}

	err = services.UpdateService()
	if err != nil {
		common.Log.Errorf("Update failed %v", err)
		c.JSON(400, err)
		return
	}

	var resp dto.FileInfo
	resp.Downloading = true
	c.JSON(200, resp)

}
