package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	_db *gorm.DB
	red *redis.Client
)

// 数据库和redis缓存配置信息结构体
type Config struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Redis struct {
		Addr        string `yaml:"addr"`
		Password    string `yaml:"password"`
		DB          int    `yaml:"db"`
		PoolSize    int    `yaml:"poolSize"`
		MinIdleConn int    `yaml:"minIdleConn"`
	} `yaml:"redis"`
}

// 数据库连接和redis缓存的连接
func Init() {
	// 读取YAML配置文件
	configData, err := ioutil.ReadFile("config/conf.yml")
	if err != nil {
		panic(err)
	}
	Config := Config{}
	err = yaml.Unmarshal(configData, &Config)
	if err != nil {
		log.Fatal(err)
	}
	// 使用配置信息连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.Database.User,
		Config.Database.Password,
		Config.Database.Host,
		Config.Database.Name)
	//自定义日志模板，打印sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql阈值
			LogLevel:      logger.Info, //级别
			Colorful:      true,        //彩色
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		//数据库迁移时禁用表名复数形式
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	_db = db
	migration()
	//配置redis
	red = redis.NewClient(&redis.Options{
		Addr:         Config.Redis.Addr,
		Password:     Config.Redis.Password,
		DB:           Config.Redis.DB,
		PoolSize:     Config.Redis.PoolSize,
		MinIdleConns: Config.Redis.MinIdleConn,
	})
}

// 创建实体类数据库对象的方法
func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}

// 创建Redis对象的方法
func NewRedisClient(ctx context.Context) *redis.Client {
	_red := red
	return _red.WithContext(ctx)
}
