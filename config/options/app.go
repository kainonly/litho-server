package options

type AppOption struct {
	Name  string `yaml:"name"`
	Key   string `yaml:"key"`
	Debug bool   `yaml:"debug"`
}
