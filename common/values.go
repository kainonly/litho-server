package common

import "fmt"

type Values struct {
	Mode      string   `yaml:"mode"`
	Address   string   `yaml:"address"`
	Ip        string   `yaml:"ip"`
	Namespace string   `yaml:"namespace"`
	Key       string   `yaml:"key"`
	Domain    string   `yaml:"domain"`
	Database  Database `yaml:"database"`
	Nats      Nats     `yaml:"nats"`
	Cos       Cos      `yaml:"cos"`
	Otlp      Otlp     `yaml:"otlp"`
}

type Database struct {
	Url   string `yaml:"url"`
	Name  string `yaml:"name"`
	Redis string `yaml:"redis"`
	Debug bool   `yaml:"debug"`
}

type Nats struct {
	Hosts []string `yaml:"hosts"`
	Token string   `yaml:"token"`
}

type Cos struct {
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	SecretId  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
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

func (x Values) LogName(key string) string {
	if x.IsRelease() {
		return fmt.Sprintf(`%s_logs`, key)
	}
	return fmt.Sprintf(`%s_logs_dev`, key)
}
