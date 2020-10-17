package bootstrap

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"os"
	"time"
	"van-api/types"
)

var (
	LoadConfigurationNotExists = errors.New("the configuration file does not exist")
)

func LoadConfiguration() (cfg types.Config, err error) {
	if _, err = os.Stat("./config/config.yml"); os.IsNotExist(err) {
		err = LoadConfigurationNotExists
		return
	}
	var buf []byte
	buf, err = ioutil.ReadFile("./config/config.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return
	}
	return
}

// Initialize database configuration
func InitializeDatabase(option *types.MysqlOption) (db *gorm.DB, err error) {
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

// Initialize the cache library configuration
func InitializeRedis(option *types.RedisOption) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     option.Address,
		Password: option.Password,
		DB:       option.DB,
	})
}
