package config

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const file = "config.yaml"

// todo 此结构体去除 ServiceNamelist 和 ServicePortlist 的定义，增加LED控制结构体，改完此配置结构体后再改其他地方
// Config 配置文件结构
type Config struct {
	ServiceNamelist  []string
	ServicePortlist  []string
	Port             int
	ShellPath        string
	MDNS             MDNS
	MicroServiceList []MicroService
}

// todo 指示灯正常和异常分别对应什么表现，写到此处的注释里
// 健康建查有失败的，status灯就正常，全部健康则异常
// WWW和Link灯，每个灯对应一个URL列表。对每个URL的GET返回均正常，指示灯正常，任何一个URL返回不正常，指示灯异常
// URL的GET返回正常是指：HTTP返回码等于200，且body里不包含以下字符串的任意一个"err, fail, disconnect, timeout"
type StLedControl struct {
	URLForWWWLed    []string
	URLForStatusLed []string
}

// todo 这个结构体里增加端口号定义，增加是否健康检查的bool字段，去除LED定义
// MicroService 微服务配置
type MicroService struct {
	Name string
	URL  string
	LED  string
	Type string
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
		Name: "conthing",
		Port: 42424,
	},
	MicroServiceList: []MicroService{
		{
			Name: "lpr",
			Type: "ping",
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
