package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/conthing/sysmgmt/config"
	"github.com/conthing/sysmgmt/router"
	"github.com/conthing/sysmgmt/services"
	"time"

	"github.com/conthing/utils/common"
)

func boot(_ interface{}) (needRetry bool, err error) {
	var cfgfile string

	//解析命令行参数 -c <cfgfile> 默认config.yaml
	flag.StringVar(&cfgfile, "config", "config.yaml", "Specify a config file other than default.")
	flag.StringVar(&cfgfile, "c", "config.yaml", "Specify a config file other than default.")
	flag.Parse()

	common.InitLogger(&common.LoggerConfig{Level: "DEBUG", SkipCaller: true})

	err = common.LoadYaml(cfgfile, &config.Conf)
	if err != nil {
		return false, fmt.Errorf("Failed to load config %w", err)
	}
	common.Log.Infof("Load config success %+v", config.Conf)

	if config.Conf.MainInterface != "" {
		err = common.SetMajorInterface(config.Conf.MainInterface)
		if err != nil {
			return false, fmt.Errorf("Failed to set main interface: %w", err)
		}
		common.Log.Infof("Set main interface %s success, IP: %s", config.Conf.MainInterface, common.GetMajorInterfaceIP())
	}

	err = services.StartMDNS(&config.Conf.MDNS)
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

	// 1.SIGINT 2.httpserver 3.serial recv
	errs := make(chan error, 8)

	// WatchDog()

	services.ScheduledHealthCheck()
	go router.Service(&config.Conf.HTTP)

	listenForInterrupt(errs)

	// recv error channel
	c := <-errs
	common.Log.Errorf("terminating: %v", c)

	//优雅退出的操作
	services.StopMDNS()

	os.Exit(0)
}

func listenForInterrupt(errChan chan error) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
}

// WatchDog 看门狗
func WatchDog() {
	go func() {
		wdt, err := services.GetWatchDog(10) //10s超时
		if err == nil {
			for {
				select {
				case <-time.After(time.Second * 4):
					err = services.KeepAlive(wdt) //10s超时
					if err != nil {
						common.Log.Errorf("feed dog failed: %v", err)
					} else {
						common.Log.Debug("feed dog ok")
					}
				}

			}

		} else {
			common.Log.Errorf("watchdog init failed: ", err)
		}
	}()
}
