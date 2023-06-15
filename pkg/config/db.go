package config

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func NewDB(config *Config) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,
			// lennon todo db log
			//LogLevel:      logger.Info,
			LogLevel: logger.Error,
			//LogLevel: logger.Info,
			Colorful: true,
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       config.DB.DataSourceName, // DSN data source name
		DefaultStringSize:         256,                      // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                    // 根据版本自动配置
	}), &gorm.Config{

		Logger: newLogger,
		//Logger: logger.Default.LogMode(logger.Error),

		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if db != nil {
		gorm.ErrRecordNotFound = errors.New("未找到数据")
	}
	return db, err
}
