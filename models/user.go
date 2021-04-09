package models

import (
	"github.com/jinzhu/gorm"
)

// User 数据库模型
type User struct {
	gorm.Model
	Username    string `gorm:"uniqueindex;<-:create;not null"`
	Hash        string //hash后的密码
	TokenRandom int    //jwt中所带的信息，用于logout或修改密码后让token失效
}
