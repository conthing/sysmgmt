package config

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const file = "config.yaml"

// Config 配置文件结构
type Config struct {
	ControlLed       StLedControl
	Port             int
	ShellPath        string
	MDNS             MDNS
	MicroServiceList []MicroService
}

// 健康检查有失败的，status灯就闪烁，全部健康则常亮
// WWW和Link灯，每个灯对应一个URL列表。对每个URL的GET返回均正常，指示灯常亮，任何一个URL返回不正常，指示灯常灭
// URL的GET返回正常是指：HTTP返回码等于200，且body里不包含以下字符串的任意一个"err, fail, disconnect, timeout"

type StLedControl struct {
	URLForWWWLed  []string
	URLForLinkLed []string
}

// MicroService 微服务配置
type MicroService struct {
	Name         string
	Port         int
	EnableHealth bool
}

// MDNS 发现服务
type MDNS struct {
	Name string
	Port int
}

var Conf = &Config{}

//Service 配置初始化服务
func Service() {

	if !exists(file) {
		createConfigFile()
	}
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("读取配置文件失败", err)
	}

	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Fatal("配置文件序列化失败", err)
	}

}

// exists 判断所给路径文件/文件夹是否存在
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

func createConfigFile() {
	buf := new(bytes.Buffer)
	err := yaml.NewEncoder(buf).Encode(Conf)
	if err != nil {
		log.Fatal("配置文件编码失败", err)
	}

	f, err := os.Create(file)
	if err != nil {
		log.Fatal("配置文件创建失败", err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal("配置文件关闭失败", err)
		}

	}()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		log.Fatal("配置文件写入失败", err)
	}
}
