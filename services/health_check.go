package services

import (
	"fmt"
	"net/http"
	"sysmgmt-next/config"
	"time"

	"github.com/conthing/utils/common"
)

// todo review 这里返回故障时要做到：如果两个微服务ping失败，error信息中String为"xxx yyy health check failed"，有点难，不要花太长时间，我会教你的
// HealthCheck 健康检查
func HealthCheck() error {
	servicename := []string{}
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		if microservice.EnableHealth == true {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/ping", microservice.Port))
			if err != nil || resp.StatusCode != 200 {
				// todo again 这里意思也对了，但如果有3个微服务出错呢？所以这种用1，2的方式不行，还有日志不要用中文
				servicename = append(servicename, microservice.Name)
				common.Log.Error("%v health check failed", servicename)
			}
			defer resp.Body.Close()
		} else {
			common.Log.Info("%s health check success", microservice.Name)
			return nil
		}
	}
	return nil
}

// todo again 1.这个函数名和结构体中冲突，最好按照命名规则区分开，2.这里遍历的是MicroServiceList吗？3.函数体内出现了mesh字眼，说明肯定写错了
func LedStatus() error {
	microservice := config.Conf.ControlLed
	for _, wwwurl := range microservice.URLForWWWLed {
		err := CheckURL(wwwurl)
		if err != nil {
			setLed(constLedWWW, constLedOff)
			common.Log.Error("WWW Led is off")
			return err
		} else {
			setLed(constLedWWW, constLedOn)
			common.Log.Info("WWW Led is on")
			return nil
		}
	}
	for _, linkurl := range microservice.URLForLinkLed {
		err := CheckURL(linkurl)
		if err != nil {
			setLed(constLedLink, constLedOff)
			common.Log.Error("Link Led is off")
			return err
		} else {
			setLed(constLedLink, constLedOn)
			common.Log.Info("Link Led is on")
			return nil
		}
	}
	return nil
}

// todo 改写此函数，将健康检查和LED控制放到同一个go程里
// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	go func() {
		LedStatus()
		//setLed() // todo review 因为setLed不对，所以这个写法不对；应该讲你现在实现的setLed里的一些内容搬过来
		HealthCheck()
		time.Sleep(30 * time.Second)
	}()
}
