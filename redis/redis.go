package redis

import (
	"log"

	"github.com/conthing/utils/common"
	"github.com/mediocregopher/radix/v3"
)

// Client redis调用客户端
var Client *radix.Pool

// Connect 初始化连接池
func Connect() {
	var err error
	Client, err = radix.NewPool("tcp", "127.0.0.1:6379", 10)
	if err != nil {
		log.Fatal("数据库连接失败", err)
	}
	common.Log.Info("redis 启动成功")
}
