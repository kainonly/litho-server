package bootstrap

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"os"
	"taste-api/config"
	"time"
)

var (
	LoadConfigurationNotExists = errors.New("the configuration file does not exist")
)

// Load application configuration
// reference config.example.yml
func LoadConfiguration() (config *config.Config, err error) {
	if _, err = os.Stat("./config.yml"); os.IsNotExist(err) {
		err = LoadConfigurationNotExists
		return
	}
	var buf []byte
	buf, err = ioutil.ReadFile("./config.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return
	}
	return
}

// Initialize database configuration
// If it is another database, replace the driver here
// gorm.Open(mysql.Open(option.Dsn),...)
// reference https://gorm.io/docs/connecting_to_the_database.html
func InitializeDatabase(config *config.Config) (db *gorm.DB, err error) {
	option := config.Database
	db, err = gorm.Open(mysql.Open(option.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   option.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	if option.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(option.MaxIdleConns)
	}
	if option.MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(option.MaxOpenConns)
	}
	if option.ConnMaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(option.ConnMaxLifetime))
	}
	return
}

// Initialize the redis library configuration
// reference https://github.com/go-redis/redis
func InitializeRedis(config *config.Config) *redis.Client {
	option := config.Redis
	return redis.NewClient(&redis.Options{
		Addr:     option.Address,
		Password: option.Password,
		DB:       option.DB,
	})
}

func HttpServer(config *config.Config) {
	serve := gin.Default()
	serve.Run(config.Listen)
}
