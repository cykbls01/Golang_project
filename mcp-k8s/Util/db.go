package Util

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mcp-k8s/Model"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := "cm@cmscreen#dkobc04:Docker@123@tcp(10.17.174.33:3306)/mcp?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = DB.AutoMigrate(&Model.Service{})
	return DB
}
