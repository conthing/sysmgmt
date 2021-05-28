package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/conthing/sysmgmt/db"
	"github.com/conthing/sysmgmt/handlers"
	"github.com/conthing/sysmgmt/services"

	"github.com/conthing/utils/common"
)

// Config 配置文件结构
type Config struct {
	services.HealthCheckConfig
	HTTP          HTTPConfig
	MainInterface string
	MDNS          services.MDNSConfig
	DB            db.DBConfig
}

type HTTPConfig struct {
	Port int
}

var Conf = &Config{}

func boot(_ interface{}) (needRetry bool, err error) {
	var cfgfile string

	//解析命令行参数 -c <cfgfile> 默认config.yaml
	flag.StringVar(&cfgfile, "config", "config.yaml", "Specify a config file other than default.")
	flag.StringVar(&cfgfile, "c", "config.yaml", "Specify a config file other than default.")
	flag.Parse()

	common.InitLogger(&common.LoggerConfig{Level: "DEBUG", SkipCaller: true})

	err = common.LoadYaml(cfgfile, &Conf)
	if err != nil {
		return false, fmt.Errorf("Failed to load config %w", err)
	}
	common.Log.Infof("Load config success %+v", Conf)

	if Conf.MainInterface != "" {
		err = common.SetMajorInterface(Conf.MainInterface)
		if err != nil {
			return false, fmt.Errorf("Failed to set main interface: %w", err)
		}
		common.Log.Infof("Set main interface %s success, IP: %s", Conf.MainInterface, common.GetMajorInterfaceIP())
	}

	// 数据库初始化
	err = db.Init(&Conf.DB)
	if err != nil {
		return true, fmt.Errorf("Failed to init database: %w", err)
	}
	common.Log.Info("Init database success")

	err = services.StartMDNS(&Conf.MDNS)
	if err != nil {
		return true, err
	}

	return false, nil
}

func main() {
	//start := time.Now()
	if common.Bootstrap(boot, nil, 1000, 500) != nil {
		return
	}

	services.HealthCheckInit(&Conf.HealthCheckConfig)

	errs := make(chan error, 8)

	//startWatchDog(errs)
	startHealthCheck(errs)
	startHTTPServer(errs)
	listenForEvents(errs)
	listenForInterrupt(errs)

	// recv error channel
	c := <-errs
	common.Log.Errorf("terminating: %v", c)

	//优雅退出的操作
	services.StopMDNS()

	os.Exit(0)
}

func startHTTPServer(errChan chan error) {
	go func() {
		ret := handlers.Run(Conf.HTTP.Port)
		if ret != nil {
			errChan <- ret
		}
	}()
}

func listenForInterrupt(errChan chan error) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
}

func listenForEvents(errChan chan error) {
	go func() {
		ret := services.ButtonSevcie()
		if ret != nil {
			errChan <- ret
		}
	}()
}

func startHealthCheck(errChan chan error) {
	go func() {
		ret := services.ScheduledHealthCheck()
		if ret != nil {
			errChan <- ret
		}
	}()
}

func startWatchDog(errChan chan error) {
	go func() {
		ret := services.WatchDog()
		if ret != nil {
			errChan <- ret
		}
	}()
}
