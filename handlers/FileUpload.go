package handlers

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, dto.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err = c.SaveUploadedFile(file, "/tmp/file.zip")
	if err != nil {
		common.Log.Errorf("Save file failed %v", err)
		c.JSON(http.StatusBadRequest, dto.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err = services.UpdateService()
	if err != nil {
		common.Log.Errorf("Update failed %v", err)
		c.JSON(http.StatusBadRequest, dto.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	var resp dto.FileInfo
	resp.Downloading = true

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: resp,
	})

}
