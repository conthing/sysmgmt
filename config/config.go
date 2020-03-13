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

// todo review 这段注释和上面的意思重复了，需要你添加的是指示灯的表现，亮灭闪这些，这些我不记得
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

// todo again 这里应该是没有配置时的默认值，应该是空，代码中不应该出现具体服务的名字和端口
// Conf 全局配置
//var Conf = Config{
//	ControlLed: StLedControl{
//		URLForWWWLed:  []string{},
//		URLForLinkLed: []string{},
//	},
//	Port:      0,
//	ShellPath: "",
//	MDNS: MDNS{
//		Name: "",
//		Port: 0,
//	},
//	MicroServiceList: []MicroService{
//		{
//			Name:         "",
//			Port:         0,
//			EnableHealth: false,
//		},
//	},
//}
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
