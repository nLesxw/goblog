package model

import (
	"GoBlog/pkg/config"
	"GoBlog/pkg/logger"
	"GoBlog/pkg/types"
	"fmt"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	//GORM 的 MySQL 数据库驱动导入
	"gorm.io/driver/mysql"
)

//DB gorm.DB对象
var DB *gorm.DB

//ConnectDB 初始化模型
func ConnectDB() *gorm.DB {

	var err error
	
	// 初始化 MySQL 连接信息
	gormConfig := mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
			config.GetString("database.mysql.username"),
			config.GetString("database.mysql.password"),
			config.GetString("database.mysql.host"),
			config.GetString("database.mysql.port"),
			config.GetString("database.mysql.database"),
			config.GetString("database.mysql.charset")),
	})

	var level gormlogger.LogLevel
	if config.GetBool("app.debug") {
		// 读取不到数据也会显示
		level = gormlogger.Warn
	} else {
		// 只有错误才会显示
		level = gormlogger.Error
	}

	//准备数据库连接池
	DB, err = gorm.Open(gormConfig, &gorm.Config{
		Logger: gormlogger.Default.LogMode(level),
	})

	logger.LogError(err)

	return DB
}

// BaseModel 模型基类
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;not null"`

	CreatedAt time.Time `gorm:"column:created_at;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;index"`
}

// GetStringID 获取ID的字符串格式
func (b BaseModel) GetStringID() string{
	return types.Uint64ToString(b.ID)
}