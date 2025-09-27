package common

type Values struct {
	Address  string   `yaml:"address"`
	Key      string   `yaml:"key"`
	Domain   string   `yaml:"domain"`
	Database Database `yaml:"database"`
}

type Database struct {
	Url  string `yaml:"url"`
	Name string `yaml:"name"`
}
