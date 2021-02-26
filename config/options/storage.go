package options

type Storage struct {
	Type   string                 `yaml:"type"`
	Option map[string]interface{} `yaml:"option"`
}
