package common

type Values struct {
	Address  string   `yaml:"address"`
	Key      string   `yaml:"key"`
	Domain   string   `yaml:"domain"`
	Database Database `yaml:"database"`
}

type Cors struct {
	Origins  []string `yaml:"origins"`
	SameSite string   `yaml:"same_site,omitempty"`
}

type Database struct {
	Url   string `yaml:"url"`
	Name  string `yaml:"name"`
	Redis string `yaml:"redis"`
}
