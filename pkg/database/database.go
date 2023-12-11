package database

import(
	"database/sql"
    "GoBlog/pkg/logger"
    "time"

    "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Initialize 初始化数据库
func Initialize() {
    initDB()
    createTables()
}

func initDB(){
    var err error
    config := mysql.Config{
        User: "root",
        Passwd: "123456",
        Addr: "127.0.0.1:3306",
        Net: "tcp",
        DBName: "goblogNew",
        AllowNativePasswords: true,
    }

    //准备数据库连接池
    DB, err = sql.Open("mysql", config.FormatDSN())
   logger.LogError(err)

    //设置最大的连接数
    DB.SetMaxOpenConns(25)
    //设置最大的空闲连接数
    DB.SetMaxIdleConns(25)
    //设置每个连接的过期时间
    DB.SetConnMaxLifetime(5 * time.Minute)

    //尝试连接，失败则会报错
    err = DB.Ping()
   logger.LogError(err)
}

func createTables(){
    createArticlesTale := `create table if not exists articles(
        id bigint(20) primary key auto_increment not null,
        title varchar(255) collate utf8mb4_unicode_ci not null,
        body longtext collate utf8mb4_unicode_ci
    );`

    _, err := DB.Exec(createArticlesTale)
   logger.LogError(err)
}