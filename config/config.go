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
	ServiceNamelist  []string
	ServicePortlist  []string
	Port             int
	ShellPath        string
	MDNS             MDNS
	MicroServiceList []MicroService
}

// MicroService 微服务配置
type MicroService struct {
	Name string
	URL  string
	LED  string
}

// MDNS 发现服务
type MDNS struct {
	Name string
	Port int
}

// Conf 全局配置
var Conf = Config{
	ServiceNamelist: make([]string, 0),
	ServicePortlist: make([]string, 0),
	Port:            52035,
	ShellPath:       "/app/sysmgmt/res/",
	MDNS: MDNS{
		Name: "主机",
		Port: 42424,
	},
	MicroServiceList: []MicroService{
		MicroService{
			Name: "lpr",
			URL:  "http://localhost:52032/api/v1/status",
			LED:  "/dev/led-pwm3",
		},
	},
}

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
