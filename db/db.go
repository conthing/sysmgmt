package db

import (
	"fmt"
	"os"

	"github.com/conthing/sysmgmt/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DataBase 数据库配置
type DBConfig struct {
	User     string
	Password string
	Location string
	DBName   string
}

var dbClient *gorm.DB

// Init connect and migrate
func Init(cfg *DBConfig) error {
	var err error

	os.MkdirAll(cfg.Location, os.ModePerm)
	//dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Location, cfg.DBName)
	//db, err = gorm.Open(mysql.New(mysql.Config{
	//	DSN:                       dsn,   // DSN data source name
	//	DefaultStringSize:         256,   // string 类型字段的默认长度
	//	DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	//	DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	//	DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	//	SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	//}), &gorm.Config{})
	dbClient, err = gorm.Open(sqlite.Open(cfg.Location+cfg.DBName+".sqlite3"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("db connect failed: %v", err)
	}
	dbClient.AutoMigrate(
		&models.User{},
		&models.Envior{},
	)
	return nil
}
