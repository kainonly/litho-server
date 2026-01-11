package common

import "fmt"

type Values struct {
	Mode      string   `yaml:"mode"`
	Address   string   `yaml:"address"`
	Namespace string   `yaml:"namespace"`
	Key       string   `yaml:"key"`
	Domain    string   `yaml:"domain"`
	Cors      Cors     `yaml:"cors"`
	Database  Database `yaml:"database"`
}

type Cors struct {
	Origins  []string `yaml:"origins"`
	SameSite string   `yaml:"same_site,omitempty"`
}

type Database struct {
	Debug bool   `yaml:"debug"`
	Url   string `yaml:"url"`
	Name  string `yaml:"name"`
	Redis string `yaml:"redis"`
}

type Nats struct {
	Hosts []string `yaml:"hosts"`
	Token string   `yaml:"token"`
}

type Cos struct {
	Bucket    string `yaml:"bucket"`
	Files     string `yaml:"files"`
	Region    string `yaml:"region"`
	SecretId  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}

type Sms struct {
	SecretId  string            `yaml:"secret_id"`
	SecretKey string            `yaml:"secret_key"`
	Sign      string            `yaml:"sign"`
	AppId     string            `yaml:"app_id"`
	Templates map[string]string `yaml:"templates"`
}

type Otlp struct {
	Name        string `yaml:"name"`
	Endpoint    string `yaml:"endpoint"`
	Token       string `yaml:"token"`
	Environment string `yaml:"environment"`
}

func (x Values) IsRelease() bool {
	return x.Mode == "release"
}

func (x Values) IsSqlDebug() bool {
	return x.Database.Debug
}

func (x Values) KeyName(i string) string {
	return fmt.Sprintf("%s:%s", x.Namespace, i)
}

func (x Values) LogName(i string) string {
	if x.IsRelease() {
		return fmt.Sprintf(`%s_logs`, i)
	}
	return fmt.Sprintf(`%s_logs_dev`, i)
}
