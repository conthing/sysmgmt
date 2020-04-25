package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	Recovery         StRecovery
}

type StRecovery struct {
	Contains    string
	Command     string
	Parameter   []string
	Environment []string
	OutputFile  string
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
func Service() error {
	var err error
	if !exists(file) {
		err = createConfigFile()
		if err != nil {
			return fmt.Errorf("create file failed %v", err)
		}
	}
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("ReadFile failed %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		return fmt.Errorf("Unmarshal failed %v", err)
	}
	return nil
}

// exists 判断所给路径文件/文件夹是否存在
func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func createConfigFile() error {
	buf := new(bytes.Buffer)
	err := yaml.NewEncoder(buf).Encode(Conf)
	if err != nil {
		return fmt.Errorf("Encode config file failed %v", err)
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("Create config file failed %v", err)
	}
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("Write config file failed %v", err)
	}
	return nil
}
