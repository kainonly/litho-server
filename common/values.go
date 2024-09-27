package common

type Values struct {
	Mode      string `yaml:"mode"`
	Address   string `yaml:"address"`
	Namespace string `yaml:"namespace"`
	Key       string `yaml:"key"`
	Database  `yaml:"database"`
	Nats      `yaml:"nats"`
}

type Database struct {
	Url   string `yaml:"url"`
	Redis string `yaml:"redis"`
}

type Nats struct {
	Hosts []string `yaml:"hosts"`
	Nkey  string   `yaml:"nkey"`
}

func (x Values) IsRelease() bool {
	return x.Mode == "release"
}
