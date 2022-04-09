package common

import (
	"fmt"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/passport"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

func SetValues(path string) (values *Values, err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("静态配置不存在，请检查路径 [%s]", path)
	}
	var config []byte
	if config, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if err = yaml.Unmarshal(config, &values); err != nil {
		return
	}
	return
}

type Values struct {
	TrustedProxies []string                 `yaml:"trusted_proxies"`
	Namespace      string                   `yaml:"namespace"`
	Key            string                   `yaml:"key"`
	Cors           Cors                     `yaml:"cors"`
	Database       Database                 `yaml:"database"`
	Redis          Redis                    `yaml:"redis"`
	Nats           Nats                     `yaml:"nats"`
	Passport       passport.Option          `yaml:"passport"`
	Engines        map[string]engine.Option `yaml:"engines"`
	QCloud         QCloud                   `yaml:"qcloud"`
}

func (x *Values) KeyName(v ...string) string {
	return fmt.Sprintf(`%s:%s`, x.Namespace, strings.Join(v, ":"))
}

func (x *Values) EventName(v string) string {
	return fmt.Sprintf(`%s.events.%s`, x.Namespace, v)
}

type Cors struct {
	AllowOrigins     []string `yaml:"allowOrigins"`
	AllowMethods     []string `yaml:"allowMethods"`
	AllowHeaders     []string `yaml:"allowHeaders"`
	ExposeHeaders    []string `yaml:"exposeHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	MaxAge           int      `yaml:"maxAge"`
}

type Database struct {
	Uri    string `yaml:"uri"`
	DbName string `yaml:"dbName"`
}

type Redis struct {
	Uri string `yaml:"uri"`
}

type Nats struct {
	Hosts []string `yaml:"hosts"`
	Nkey  string   `yaml:"nkey"`
}

type QCloud struct {
	SecretID  string    `yaml:"secret_id"`
	SecretKey string    `yaml:"secret_key"`
	Cos       QCloudCos `yaml:"cos"`
}

type QCloudCos struct {
	Bucket  string `yaml:"bucket"`
	Region  string `yaml:"region"`
	Expired int64  `yaml:"expired"`
}

type Subscriptions struct {
	*sync.Map
}
