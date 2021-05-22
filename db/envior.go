package db

import (
	"errors"
	"fmt"

	"github.com/conthing/utils/common"
	"gorm.io/gorm"
)

func GetEnviorList() ([]Envior, error) {
	list := make([]Envior, 0)
	err := dbClient.Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("dbClient.Find Envior list failed: %w", err)
	}
	return list, nil
}

// 根据obj的可识别字段找到
func getEnviorEntry(obj *Envior) (tx *gorm.DB) {
	tx = dbClient
	//判断可识别字段是否有效
	if obj.Name != "" {
		tx = tx.Where("name = ?", obj.Name)
	}
	return
}

// GetEnvior
func GetEnvior(obj *Envior) error {
	err := getEnviorEntry(obj).First(obj).Error
	if err != nil {
		return fmt.Errorf("dbClient.First device failed: %w", err)
	}
	return nil
}

// AddEnvior
func AddEnvior(obj *Envior) error {
	err := dbClient.Create(obj).Error
	if err != nil {
		return fmt.Errorf("dbClient.Create Envior failed: %w", err)
	}
	return nil
}

// SetEnvior
func SetEnvior(obj *Envior) error {
	value := obj.Value

	err := getEnviorEntry(obj).FirstOrCreate(obj).Update("value", value).Scan(obj).Error
	if err != nil {
		return fmt.Errorf("dbClient.Update device failed: %w", err)
	}
	return err
}

// DeleteEnvior
func DeleteEnvior(obj *Envior) error {
	err := getEnviorEntry(obj).First(obj).Unscoped().Delete(obj).Error
	if err != nil {
		return fmt.Errorf("dbClient.Delete envior failed:%w", err)
	}
	return nil
}

// GetEnv 读取环境变量
func GetEnv(name string) string {
	if name == "" {
		return ""
	}
	obj := &Envior{Name: name}
	if err := GetEnvior(obj); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) { //是not found错误，不输出日志
			common.Log.Errorf("GetEnvior %q failed: %v", name, err)
		}
		return ""
	}
	return obj.Value
}

// SetEnv 设置环境变量，如果value为空，则删除此环境变量
func SetEnv(name string, value string) error {
	if name == "" {
		return fmt.Errorf("Set Envoir with empty name")
	}
	obj := &Envior{Name: name, Value: value}
	if value == "" {
		if err := DeleteEnvior(obj); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) { //是not found错误，不输出日志
				common.Log.Errorf("Delete Envior %q failed: %v", name, err)
				return err
			}
		}
	} else {
		if err := SetEnvior(obj); err != nil {
			common.Log.Errorf("Set Envior %q=%q failed: %v", name, value, err)
			return err
		}
	}

	return nil
}
