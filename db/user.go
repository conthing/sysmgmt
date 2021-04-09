package db

import (
	"fmt"

	"github.com/conthing/sysmgmt/models"

	"gorm.io/gorm"
)

// GetUserList 获取物理设备列表
func GetUserList() ([]models.User, error) {
	list := make([]models.User, 0)
	err := dbClient.Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("dbClient.Find User list failed: %w", err)
	}
	return list, nil
}

// 根据obj的可识别字段找到
func getUserEntry(obj *models.User) (tx *gorm.DB) {
	tx = dbClient
	//判断可识别字段是否有效
	if obj.ID != 0 {
		tx = tx.Where("id = ?", obj.ID)
	}
	if obj.Username != "" {
		tx = tx.Where("username = ?", obj.Username)
	}
	return
}

// GetUser
func GetUser(obj *models.User) error {
	err := getUserEntry(obj).First(obj).Error
	if err != nil {
		return fmt.Errorf("dbClient.First User failed: %w", err)
	}
	return nil
}

// AddUser
func AddUser(obj *models.User) error {
	err := dbClient.Create(obj).Error
	if err != nil {
		return fmt.Errorf("dbClient.Create User failed: %w", err)
	}
	return nil
}

// DeleteUser
func DeleteUser(obj *models.User) error {
	err := getUserEntry(obj).First(obj).Error
	if err != nil {
		return fmt.Errorf("Couldn't find deleting User in database: %w", err)
	}
	err = dbClient.Delete(obj).Error
	if err != nil {
		return fmt.Errorf("dbClient.Delete User failed:%w", err)
	}
	return nil
}

// ModifyUser
func ModifyUser(obj *models.User) error {
	hash := obj.Hash
	tokenRandom := obj.TokenRandom

	var err error
	if hash == "" {
		err = getUserEntry(obj).First(obj).Updates(map[string]interface{}{"token_random": tokenRandom}).Error
	} else {
		err = getUserEntry(obj).First(obj).Updates(map[string]interface{}{"token_random": tokenRandom, "hash": hash}).Error
	}
	if err != nil {
		return fmt.Errorf("dbClient.Update user failed: %w", err)
	}
	if hash != "" {
		obj.Hash = hash
	}
	obj.TokenRandom = tokenRandom
	return err
}
