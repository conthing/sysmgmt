package redis

import (
	"github.com/conthing/utils/common"
	"github.com/mediocregopher/radix/v3"
)

// SaveAlias 保存别名
func SaveAlias(alias string) error {
	key := "alias"
	err := Client.Do(radix.Cmd(nil, "SET", key, alias))
	if err != nil {
		common.Log.Error("设置 PubHTTPInfo 错误", err)
		return err
	}
	return nil
}

// GetAlias 获取别名
func GetAlias() (string, error) {
	key := "alias"
	var alias string
	err := Client.Do(radix.Cmd(&alias, "GET", key))
	if err != nil {
		common.Log.Error("获取 PubHTTPInfo 错误", err)
		return alias, err
	}
	return alias, nil
}
