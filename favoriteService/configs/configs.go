package configs

import (
	"favoriteService/models"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type MySQLConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Url      string `yaml:"url"`
	Port     string `yaml:"port"`
}

type GlobalConfig struct {
	MySQLConf MySQLConfig `yaml:"mysql"`
}

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
	var config GlobalConfig
	conf := config.getConf()
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
	models.InitRedis()
}

func Close() {
	models.CloseMySQL()
}
