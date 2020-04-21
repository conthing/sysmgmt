package redis

import (
	"github.com/conthing/utils/common"
	"github.com/mediocregopher/radix/v3"
)

// SaveRegion 保存地区
func SaveRegion(alias string) error {
	key := "region"
	err := Client.Do(radix.Cmd(nil, "SET", key, alias))
	if err != nil {
		common.Log.Error("SaveRegion ERR: ", err)
		return err
	}
	return nil
}

// GetRegion 获取地区
func GetRegion() (string, error) {
	key := "region"
	var alias string
	err := Client.Do(radix.Cmd(&alias, "GET", key))
	if err != nil {
		common.Log.Error("SaveRegion ERR: ", err)
		return alias, err
	}
	return alias, nil
}
