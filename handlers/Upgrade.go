package handlers

import (
	"net/http"

	"github.com/conthing/sysmgmt/dto"
	"github.com/conthing/sysmgmt/services"

	"github.com/conthing/utils/common"
	"github.com/gin-gonic/gin"
)

// todo 导入导出功能没做
//Upgrade 升级程序
func Upgrade(c *gin.Context) {
	common.Log.Debugf("Upgrade start...")

	file, err := c.FormFile("file")
	if err != nil {
		common.Log.Errorf("Form file failed %v", err)
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	services.Clean()
	common.Log.Debugf("formed file.zip")

	err = c.SaveUploadedFile(file, "/tmp/file.zip")
	if err != nil {
		common.Log.Errorf("Save file failed %v", err)
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	common.Log.Debugf("saved to file.zip")

	err = services.UpdateService()
	if err != nil {
		common.Log.Errorf("Update failed %v", err)
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	var resp dto.FileInfo
	resp.Downloading = true

	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Data: resp,
	})

}

//Import 导入设置
func Import(c *gin.Context) {
	common.Log.Debugf("Import start...")

	file, err := c.FormFile("file")
	if err != nil {
		common.Log.Errorf("Form file failed %v", err)
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	services.Clean()
	common.Log.Debugf("formed file.zip")

	err = c.SaveUploadedFile(file, "/tmp/file.zip")
	if err != nil {
		common.Log.Errorf("Save file failed %v", err)
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	common.Log.Debugf("saved to file.zip")

	err = services.ImportService()
	if err != nil {
		common.Log.Errorf("Import failed %v", err)
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	var resp dto.FileInfo
	resp.Downloading = true

	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Data: resp,
	})

}

// FileInfo 文件信息
type ExportInfo struct {
	URL string `json:"url"`
}

//Export 导出设置
func Export(c *gin.Context) {

	if err := services.ZipData(); err != nil {
		common.Log.Errorf("Zip data failed %v", err)
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Data: ExportInfo{URL: "files/data.zip"},
	})

}
