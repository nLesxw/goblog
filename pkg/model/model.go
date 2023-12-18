package model

import (
	"GoBlog/pkg/logger"
	"gorm.io/gorm"

	//GORM 的 MySQL 数据库驱动导入
	"gorm.io/driver/mysql"
	
)

//DB gorm.DB对象
var DB *gorm.DB

//ConnectDB 初始化模型
func ConnectDB() *gorm.DB {

	var err error
	//mysql.Config{DSN: "root:123456@tcp(127.0.0.1:3306)/goblog?charset=utf8&parseTime=True&loc=Local",}
	config := mysql.New(mysql.Config{
		DSN: "root:123456@tcp(127.0.0.1:3306)/goblognew?charset=utf8&parseTime=True&loc=Local",
	})
	//准备数据库连接池
	DB, err = gorm.Open(config, &gorm.Config{})

	logger.LogError(err)

	return DB
}