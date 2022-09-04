package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"

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

func hashSHA256File(filePath string) (string, error) {
	var hashValue string
	file, err := os.Open(filePath)
	if err != nil {
		return hashValue, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return hashValue, err
	}
	hashInBytes := hash.Sum(nil)
	hashValue = hex.EncodeToString(hashInBytes)
	return hashValue, nil
}

func UrlUpgrade(c *gin.Context) {
	var info dto.UrlUpgradeInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	common.Log.Debugf("GET from %s", info.URL)

	// Get the data
	resp, err := http.Get(info.URL)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	services.Clean()

	// 创建一个文件用于保存
	out, err := os.Create("/tmp/file.zip")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	common.Log.Debugf("saved to file.zip")

	sha256, err := hashSHA256File("/tmp/file.zip")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if sha256 != info.SHA256 {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: "SHA256 failed",
		})
		return
	}

	common.Log.Debugf("SHA256 success")

	err = services.UpdateService()
	if err != nil {
		common.Log.Errorf("Update failed %v", err)
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

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
