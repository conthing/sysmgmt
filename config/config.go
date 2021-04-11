package config

import "github.com/conthing/sysmgmt/db"

// Config 配置文件结构
type Config struct {
	ControlLed       StLedControl
	HTTP             HTTPConfig
	ShellPath        string
	MainInterface    string
	MDNS             MDNSConfig
	MicroServiceList []MicroService
	Recovery         StRecovery
	DB               db.DBConfig
}

type HTTPConfig struct {
	Port int
}

// 发现服务配置
type MDNSConfig struct {
	Enable bool
	Name   string
	Port   int
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

var Conf = &Config{}
