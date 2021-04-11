package db

import (
	"fmt"

	"github.com/conthing/sysmgmt/models"

	"gorm.io/gorm"
)

// 根据obj的可识别字段找到
func getEnviorEntry(obj *models.Envior) (tx *gorm.DB) {
	tx = dbClient
	//判断可识别字段是否有效
	if obj.ID != 0 {
		tx = tx.Where("id = ?", obj.ID)
	}
	if obj.Name != "" {
		tx = tx.Where("name = ?", obj.Name)
	}
	return
}

// GetEnvior
func GetEnvior(obj *models.Envior) error {
	err := getEnviorEntry(obj).First(obj).Error
	if err != nil {
		return fmt.Errorf("dbClient.First Envior failed: %w", err)
	}
	return nil
}

// AddEnvior
func AddEnvior(obj *models.Envior) error {
	err := dbClient.Create(obj).Error
	if err != nil {
		return fmt.Errorf("dbClient.Create Envior failed: %w", err)
	}
	return nil
}

// DeleteEnvior
func DeleteEnvior(obj *models.Envior) error {
	err := getEnviorEntry(obj).First(obj).Error
	if err != nil {
		return fmt.Errorf("Couldn't find deleting Envior in database: %w", err)
	}
	err = dbClient.Delete(obj).Error
	if err != nil {
		return fmt.Errorf("dbClient.Delete Envior failed:%w", err)
	}
	return nil
}

// SetEnvior
func SetEnvior(obj *models.Envior) error {
	value := obj.Value
	err := getEnviorEntry(obj).FirstOrCreate(obj).Updates(map[string]interface{}{"value": value}).Error

	if err != nil {
		return fmt.Errorf("dbClient.Update Envior failed: %w", err)
	}
	obj.Value = value
	return err
}
