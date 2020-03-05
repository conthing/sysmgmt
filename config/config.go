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

// 健康检查有失败的，status灯就正常，全部健康则异常
// WWW和Link灯，每个灯对应一个URL列表。对每个URL的GET返回均正常，指示灯正常，任何一个URL返回不正常，指示灯异常
// URL的GET返回正常是指：HTTP返回码等于200，且body里不包含以下字符串的任意一个"err, fail, disconnect, timeout"

// todo review 这段注释和上面的意思重复了，需要你添加的是指示灯的表现，亮灭闪这些，这些我不记得
//Link灯(D9或led-pwm3)正常:当获取zigbee的mesh时返回是200且body里不包含以下字符串的任意一个"err, fail, disconnect, timeout"或者lpr的485运行正常
//Link灯(D9或led-pwm3)异常:当获取zigbee的mesh时返回不是200且body里包含以下字符串的任意一个"err, fail, disconnect, timeout"或者lpr的485运行不正常
//WWW灯(D8或led-pwm2)正常:当获取zap的status为connected或者(和)获取lpr的status为connected时
//WWW灯(D8或led-pwm2)异常:当获取zap的status为disconnected或者(和)获取lpr的status为disconnected时
//Status灯(D7或led-pwm1)正常:当对每个微服务健康检查(ping)时每个微服务返回都是pong
//Status灯(D7或led-pwm1)异常:当对每个微服务健康检查(ping)时只要有一个微服务返回的不是pong
type StLedControl struct {
	URLForWWWLed  []string
	URLForLinkLed []string
}

// todo review 不是ServicePortlist，应该是一个端口号；增加的bool是是否进行健康检查，不是是否健康
// MicroService 微服务配置
type MicroService struct {
	Name            string
	URL             string
	Type            string
	ServicePortlist []string
	IsHealth        bool
}

// MDNS 发现服务
type MDNS struct {
	Name string
	Port int
}

// Conf 全局配置
var Conf = Config{
	ControlLed: StLedControl{
		URLForWWWLed:  []string{"http://localhost:52032/api/v1/status", "http://localhost:52018/api/v1/status"},
		URLForLinkLed: []string{"http://localhost:52032/api/v1/mesh", "http://localhost:52032/api/v1/ping"},
	},
	Port:      52035,
	ShellPath: "/app/sysmgmt/res/",
	MDNS: MDNS{
		Name: "conthing",
		Port: 42424,
	},
	MicroServiceList: []MicroService{
		{
			Name:            "lpr",
			Type:            "ping",
			URL:             "http://localhost:52032/api/v1/status",
			ServicePortlist: []string{"52032", "52018"},
			IsHealth:        true,
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
