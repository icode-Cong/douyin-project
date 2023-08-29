package configs

import (
	"fmt"
	"os"
	"publishService/models"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis/v8"
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

// 初始化一个ossClient
var OssClient *oss.Client

// 初始化redis对象
var RedisSession *redis.Client

var conf GlobalConfig

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
	conf.getConf()
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
	mysqlConf := conf.MySQLConf
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConf.Username,
		mysqlConf.Password,
		mysqlConf.Url,
		mysqlConf.Port,
		mysqlConf.DBName,
	)
	err := models.InitMySQL(dsn)
	if err != nil {
		panic(err)
	}
}
func InitOssClient() {
	ossConf := conf.OssConf
	client, err := oss.New(ossConf.Endpoint, ossConf.AccessKeyID, ossConf.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	OssClient = client
}
