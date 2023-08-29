package configs

import (
	"fmt"
	"os"
	"publishService/models"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v3"
)

type OssConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyID"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	BucketName      string `yaml:"bucketName"`
}
type MySQLConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Url      string `yaml:"url"`
	Port     string `yaml:"port"`
}

type GlobalConfig struct {
	MySQLConf MySQLConfig `yaml:"mysql"`
	OssConf   OssConfig   `yaml:"oss"`
}

// 初始化ossClient
var OssClient *oss.Client

// 初始化redis对象
var RedisSession *redis.Client

// 初始化mysql对象
var SqlSession *gorm.DB

var Conf GlobalConfig

func (c *GlobalConfig) getConf() *GlobalConfig {
	yamlFile, err := os.ReadFile("application.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic(err)
	}
	return c
}

func Init() {
	Conf.getConf()
	InitMysql()
	InitOssClient()
	InitRedis()
}

func InitRedis() {
	RedisSession = redis.NewClient(&redis.Options{
		Addr:     "172.19.0.11:6379",
		Password: "",
		DB:       0,
	})
}
func InitMysql() {
	mysqlConf := Conf.MySQLConf
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConf.Username,
		mysqlConf.Password,
		mysqlConf.Url,
		mysqlConf.Port,
		mysqlConf.DBName,
	)
	db, err := gorm.Open("mysql", dsn)
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
	models.InitMySQL(SqlSession)
	err = SqlSession.DB().Ping()
	if err != nil {
		panic(err)
	}
}
func InitOssClient() {
	ossConf := Conf.OssConf
	client, err := oss.New(ossConf.Endpoint, ossConf.AccessKeyID, ossConf.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	OssClient = client
}
