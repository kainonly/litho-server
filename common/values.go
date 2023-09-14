package common

import (
	"fmt"
	"github.com/weplanx/go/values"
	"strings"
)

type Values struct {
	Mode      string `env:"MODE" envDefault:"debug"`
	Address   string `env:"ADDRESS" envDefault:":3000"`
	Namespace string `env:"NAMESPACE,required"`
	Key       string `env:"KEY,required"`

	Database struct {
		Url   string `env:"URL,required"`
		Name  string `env:"NAME,required"`
		Redis string `env:"REDIS,required"`
	} `envPrefix:"DATABASE_"`

	Influx struct {
		Url    string `env:"URL"`
		Org    string `env:"ORG"`
		Token  string `env:"TOKEN"`
		Bucket string `env:"BUCKET"`
	} `envPrefix:"INFLUX_"`

	Nats struct {
		Hosts []string `env:"HOSTS,required" envSeparator:","`
		Nkey  string   `env:"NKEY,required"`
	} `envPrefix:"NATS_"`

	Otlp struct {
		Endpoint string `env:"ENDPOINT"`
	} `envPrefix:"OTLP_"`

	*Extra
}

type Extra struct {
	BaseUrl               string `yaml:"base_url"`
	IpAddress             string `yaml:"ip_address"`
	IpSecretId            string `yaml:"ip_secret_id"`
	IpSecretKey           string `yaml:"ip_secret_key"`
	SmsSecretId           string `yaml:"sms_secret_id"`
	SmsSecretKey          string `yaml:"sms_secret_key" secret:"*"`
	SmsAppId              string `yaml:"sms_app_id"`
	SmsRegion             string `yaml:"sms_region"`
	EmqxHost              string `yaml:"emqx_host"`
	EmqxApiKey            string `yaml:"emqx_api_key"`
	EmqxSecretKey         string `yaml:"emqx_secret_key" secret:"*"`
	*values.DynamicValues `yaml:"dynamic_values"`
}

func (x Values) IsRelease() bool {
	return x.Mode == "release"
}

func (x Values) Name(v ...string) string {
	return fmt.Sprintf(`%s:%s`, x.Namespace, strings.Join(v, ":"))
}

func (x Values) NameX(sep string, v ...string) string {
	elems := []string{x.Namespace}
	elems = append(elems, v...)
	return strings.Join(elems, sep)
}
