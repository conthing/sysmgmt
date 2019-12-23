package services

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

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
	AddItemList []AddItem `yaml:"add-item-list"`
	PutItemList []PutItem `yaml:"put-item-list"`
	DelItemList []DelItem `yaml:"del-item-list"`
}

// UpdateService 升级服务
func UpdateService() error {
	var c = ItemCollection{}

	err := UnZip()
	if err != nil {
		log.Println("Unzip", err)

		return err
	}

	err = ReadYAML(&c)
	if err != nil {
		log.Println("Read", err)

		return err
	}

	err = Install(&c)
	if err != nil {
		log.Println("Install", err)

		return err
	}

	err = Clean()
	if err != nil {
		log.Println("Clean 2", err)

		return err
	}

	return nil
}

// Clean 清理
func Clean() error {
	if exists("/tmp/file") {
		err := os.RemoveAll("/tmp/file/")
		if err != nil {
			return err
		}
	}

	if exists("/tmp/__MACOSX") {
		err := os.RemoveAll("/tmp/__MACOSX")
		if err != nil {
			return err
		}
	}

	if exists("/tmp/file.zip") {
		err := os.RemoveAll("/tmp/file.zip")
		if err != nil {
			return err
		}
	}

	return nil
}

// Install 安装
func Install(c *ItemCollection) error {
	for _, item := range c.AddItemList {
		err := add(item)
		if err != nil {
			return err
		}
	}

	for _, item := range c.PutItemList {
		err := put(item)
		if err != nil {
			return err
		}
	}

	for _, item := range c.DelItemList {
		err := del(item)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnZip 目前在 /tmp/ 文件夹下执行
// arm 板需要安装 zip (apt install zip)
func UnZip() error {
	out, err := exec.Command("bash", "-c", "cd /tmp/ && unzip /tmp/file.zip").Output()
	log.Println(string(out))
	if err != nil {
		return err
	}
	return nil
}

// ReadYAML 读取YAML
// update.yaml 为更新清单
// 更新清单解压后放在"/tmp/file/update.yaml"
func ReadYAML(c *ItemCollection) error {
	yamlFile, err := ioutil.ReadFile("/tmp/file/update.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	log.Println(c)
	if err != nil {
		return err
	}
	return nil

}

// add 增加
func add(item AddItem) error {
	err := os.Rename(item.Src, item.Dst)
	return err
}

// put 覆盖
func put(item PutItem) error {
	err := os.RemoveAll(item.Dst)
	if err != nil {
		log.Println("removeall 失败", err)
		return err
	}
	err = os.Rename(item.Src, item.Dst)
	if err != nil {
		return err
	}
	return nil
}

// del 删除
func del(item DelItem) error {
	err := os.RemoveAll(item.Dst)
	if err != nil {
		return err
	}
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			log.Println(err)
			return true
		}
		return false
	}
	return true
}
