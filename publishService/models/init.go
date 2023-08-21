package models

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const DRIVER = "mysql"

var SqlSession *gorm.DB

var Ctx = context.Background()
var RedisSession *redis.Client

func InitMySQL(dsn string) (err error) {
	db, err := gorm.Open(DRIVER, dsn)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	if gin.Mode() == "release" {
		db.LogMode(false)
	}
	//默认不加复数
	db.SingularTable(true)
	//设置连接池
	//空闲
	db.DB().SetMaxIdleConns(20)
	//打开
	db.DB().SetMaxOpenConns(100)
	//超时
	db.DB().SetConnMaxLifetime(time.Second * 30)
	SqlSession = db
	SqlSession.AutoMigrate(&Video{})
	return SqlSession.DB().Ping()
}

func CloseMySQL() {
	err := SqlSession.Close()
	if err != nil {
		panic(err)
	}
}

func InitRedis() {
	RedisSession = redis.NewClient(&redis.Options{
		Addr:     "172.19.0.11:6379",
		Password: "",
		DB:       0,
	})
}
