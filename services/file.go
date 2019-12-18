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

// MyItemCollection 全局变量, 默认都为空数组
var MyItemCollection = ItemCollection{}

// UpdateService 升级服务
func UpdateService() error {
	err := UnZip()
	if err != nil {
		return err
	}

	err = ReadYAML()
	if err != nil {
		return err
	}

	err = Install()
	if err != nil {
		return err
	}

	err = Clean()
	if err != nil {
		return err
	}

	return nil
}

// Clean 清理
func Clean() error {
	err := os.RemoveAll("/tmp/file/")
	if err != nil {
		return err
	}
	err = os.RemoveAll("/tmp/__MACOSX")
	if err != nil {
		return err
	}
	err = os.RemoveAll("/tmp/file.zip")
	if err != nil {
		return err
	}
	return nil
}

// Install 安装
func Install() error {
	for _, item := range MyItemCollection.AddItemList {
		err := add(item)
		if err != nil {
			return err
		}
	}

	for _, item := range MyItemCollection.PutItemList {
		err := put(item)
		if err != nil {
			return err
		}
	}

	for _, item := range MyItemCollection.DelItemList {
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
func ReadYAML() error {
	yamlFile, err := ioutil.ReadFile("/tmp/file/update.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &MyItemCollection)
	if err != nil {
		return err
	}
	log.Println(MyItemCollection)
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
