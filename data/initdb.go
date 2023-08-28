package data

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func InitDb() {
 	server := os.Getenv("MYSQL_HOST")
	portStr := os.Getenv("MYSQL_PORT")
	dsn := "root:CysCOYio@tcp("+server+":"+portStr+")/douyindb?charset=utf8mb4&parseTime=True&loc=Local" //数据库信息
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 启用日志输出，输出包含SQL查询语句
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	// 迁移数据库
	Db.AutoMigrate()
}
