package services

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/conthing/utils/common"
	"gopkg.in/yaml.v2"
)

// AddItem 新增
type AddItem struct {
	Src string
	Dst string
}

// PutItem 覆盖
type PutItem struct {
	Src string
	Dst string
}

// DelItem 删除
type DelItem struct {
	Dst string
}

// ItemCollection 更新说明合集
type ItemCollection struct {
	PreScript   string    `yaml:"pre-script"`
	PostScript  string    `yaml:"post-script"`
	AddItemList []AddItem `yaml:"add-item-list"`
	PutItemList []PutItem `yaml:"modify-item-list"`
	DelItemList []DelItem `yaml:"del-item-list"`
	Reboot      bool
}

// UpdateService 升级服务
func UpdateService() error {
	var c = ItemCollection{}

	//defer Clean()
	err := UnZip()
	if err != nil {
		common.Log.Error("Unzip failed ", err)
		return err
	}
	common.Log.Debugf("unzipped to /tmp")

	err = ReadYAML(&c)
	if err != nil {
		common.Log.Error("Read failed ", err)

		return err
	}
	common.Log.Debugf("ReadYAML success")

	err = Install(&c)
	if err != nil {
		common.Log.Error("Install failed ", err)

		return err
	}
	common.Log.Debugf("Installed")

	return nil
}

// ImportService 导入服务
func ImportService() error {
	//defer Clean()
	err := UnZip()
	if err != nil {
		common.Log.Error("Unzip failed ", err)
		return err
	}
	common.Log.Debugf("unzipped to /tmp")

	if exists("/app/scripts/import.sh") {
		out, err := exec.Command("bash", "/app/scripts/import.sh", "/tmp/app/data").Output() // 用bash -c 的话$1参数会错位
		if err != nil {
			return fmt.Errorf("import.sh failed: %w", err)
		}
		common.Log.Debugf("import.sh get %s", string(out))
	} else {
		err := os.RemoveAll("/app/data/")
		if err != nil {
			return err
		}
		err = os.Rename("/tmp/app/data/", "/app/data/")
		if err != nil {
			return err
		}
	}
	common.Log.Debugf("Imported")

	RebootLater()

	return nil
}

// Clean 清理
func Clean() {
	if exists("/tmp/file") {
		err := os.RemoveAll("/tmp/file/")
		if err != nil {
			common.Log.Errorf("Remove /tmp/file/ failed: %v", err)
		} else {
			common.Log.Debugf("/tmp/file/ Removed")
		}
	}

	if exists("/tmp/app") {
		err := os.RemoveAll("/tmp/app/")
		if err != nil {
			common.Log.Errorf("Remove /tmp/app/ failed: %v", err)
		} else {
			common.Log.Debugf("/tmp/app/ Removed")
		}
	}

	if exists("/tmp/__MACOSX") {
		err := os.RemoveAll("/tmp/__MACOSX")
		if err != nil {
			common.Log.Errorf("Remove /tmp/__MACOSX failed: %v", err)
		} else {
			common.Log.Debugf("/tmp/__MACOSX Removed")
		}
	}

	if exists("/tmp/file.zip") {
		err := os.RemoveAll("/tmp/file.zip")
		if err != nil {
			common.Log.Errorf("Remove /tmp/file.zip failed: %v", err)
		} else {
			common.Log.Debugf("/tmp/file.zip Removed")
		}
	}
}

// Install 安装
func Install(c *ItemCollection) error {
	if c.PreScript != "" {
		out, err := exec.Command("bash", c.PreScript).Output()
		if err != nil {
			return fmt.Errorf("PreScript %q failed: %w", c.PreScript, err)
		}
		common.Log.Debugf("PreScript: %q get %s", c.PreScript, string(out))
	}

	for _, item := range c.AddItemList {
		err := os.Rename(item.Src, item.Dst)
		if err != nil {
			common.Log.Errorf("add item %q failed: %v", item.Dst, err)
			return fmt.Errorf("add item %q failed: %w", item.Dst, err)
		} else {
			common.Log.Debugf("added item: %q", item.Dst)
		}
	}

	for _, item := range c.PutItemList {
		//err := os.RemoveAll(item.Dst)
		//if err != nil {
		//	common.Log.Error("delete item %q failed: %v", item.Dst, err)
		//	return fmt.Errorf("delete item %q failed: %w", item.Dst, err)
		//} else {
		//	common.Log.Debugf("deleted item: %q", item.Dst)
		//}
		err := os.Rename(item.Src, item.Dst)
		if err != nil {
			common.Log.Errorf("modify item %q failed: %v", item.Dst, err)
			return fmt.Errorf("modify item %q failed: %w", item.Dst, err)
		} else {
			common.Log.Debugf("modified item: %q", item.Dst)
		}
	}

	for _, item := range c.DelItemList {
		err := os.RemoveAll(item.Dst)
		if err != nil {
			common.Log.Errorf("delete item %q failed: %v", item.Dst, err)
			return fmt.Errorf("delete item %q failed: %w", item.Dst, err)
		} else {
			common.Log.Debugf("deleted item: %q", item.Dst)
		}
	}

	if c.PostScript != "" {
		out, err := exec.Command("bash", c.PostScript).Output()
		if err != nil {
			return fmt.Errorf("PostScript %q failed: %w", c.PostScript, err)
		}
		common.Log.Debugf("PostScript: %q get %s", c.PostScript, string(out))
	}

	if c.Reboot {
		RebootLater()
	}
	return nil
}

// UnZip 目前在 /tmp/ 文件夹下执行
// arm 板需要安装 zip (apt install zip)
func UnZip() error {
	out, err := exec.Command("bash", "-c", "cd /tmp/ && unzip /tmp/file.zip").Output()
	common.Log.Debug(string(out))
	if err != nil {
		return err
	}
	return nil
}

func addCurrentDirIfNeeded(s string) string {
	if s != "" && s[0] != '/' {
		return "/tmp/file/" + s
	} else {
		return s
	}
}

// ReadYAML 读取YAML
// update.yaml 为更新清单
// 更新清单解压后放在"/tmp/file/update.yaml"
func ReadYAML(c *ItemCollection) error {
	yamlFile, err := ioutil.ReadFile("/tmp/file/update.yaml")
	if err != nil {
		common.Log.Error("ReadFile failed ", err)
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		common.Log.Error("Unmarshal failed ", err)
		return err
	}

	// 加上默认当前路径
	c.PreScript = addCurrentDirIfNeeded(c.PreScript)
	c.PostScript = addCurrentDirIfNeeded(c.PostScript)
	for i, e := range c.AddItemList {
		c.AddItemList[i].Src = addCurrentDirIfNeeded(e.Src)
	}
	for i, e := range c.PutItemList {
		c.PutItemList[i].Src = addCurrentDirIfNeeded(e.Src)
	}

	common.Log.Debugf("update.yaml: %+v", c)
	return nil

}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 压缩导出数据
func ZipData() error {
	out, err := exec.Command("zip", "-r", "/app/log/files/data.zip", "/app/data/").Output()
	common.Log.Debug(string(out))
	if err != nil {
		return err
	}
	return nil
}
