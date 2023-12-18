package bootstrap

import (
	"GoBlog/pkg/model"
	"time"
)

//SetupDB 初始化数据库和 ORM
func SetupDB() {

	//建立数据库连接池
	db := model.ConnectDB()

	//命令行打印数据库请求的信息
	sqlDB, _ := db.DB()

	//设置最大连接数
	sqlDB.SetMaxOpenConns(100)
	//设置最大空闲连接数
	sqlDB.SetConnMaxIdleTime(25)
	//设置每个连接的过期时间
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}