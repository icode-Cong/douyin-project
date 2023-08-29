package models

import (
	"context"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var SqlSession *gorm.DB

var Ctx = context.Background()

func InitMySQL(sqlSession *gorm.DB) {
	SqlSession = sqlSession
	SqlSession.AutoMigrate(&Video{})
}

func CloseMySQL() {
	err := SqlSession.Close()
	if err != nil {
		panic(err)
	}
}
