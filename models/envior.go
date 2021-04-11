package models

import (
	"github.com/jinzhu/gorm"
)

// Envior 环境变量模型
type Envior struct {
	gorm.Model
	Name  string `gorm:"uniqueindex;<-:create;not null"`
	Value string
}
