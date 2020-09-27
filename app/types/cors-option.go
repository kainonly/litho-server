package types

type CorsOption struct {
	Origin        []string `yaml:"origin"`
	Method        []string `yaml:"method"`
	AllowHeader   []string `yaml:"allow_header"`
	ExposedHeader []string `yaml:"exposed_header"`
	MaxAge        int64    `yaml:"max_age"`
	Credentials   bool     `yaml:"credentials"`
}
