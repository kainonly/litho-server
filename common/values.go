package common

import (
	"fmt"
	"github.com/weplanx/utils/values"
	"strings"
)

type Values struct {
	Address               string `env:"ADDRESS" envDefault:":3001"`
	BaseUrl               string `env:"BASE_URL,required"`
	Namespace             string `env:"NAMESPACE,required"`
	Key                   string `env:"KEY,required"`
	Database              `envPrefix:"DATABASE_"`
	Influx                `envPrefix:"INFLUX_"`
	Nats                  `envPrefix:"NATS_"`
	Otlp                  `envPrefix:"OTLP_"`
	*values.DynamicValues `env:"-"`
}

type Database struct {
	Host  string `env:"HOST,required"`
	Name  string `env:"NAME,required"`
	Redis string `env:"REDIS,required"`
}

type Influx struct {
	Url    string `env:"URL"`
	Org    string `env:"ORG"`
	Token  string `env:"TOKEN"`
	Bucket string `env:"BUCKET"`
}

type Nats struct {
	Hosts []string `env:"HOSTS,required" envSeparator:","`
	Nkey  string   `env:"NKEY,required"`
}

type Otlp struct {
	Endpoint string `env:"ENDPOINT"`
	// TODO: Improve other configuration later
}

func (x Values) Name(v ...string) string {
	return fmt.Sprintf(`%s:%s`, x.Namespace, strings.Join(v, ":"))
}
