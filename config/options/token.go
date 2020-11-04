package options

type TokenOption struct {
	Issuer   string   `yaml:"issuer"`
	Audience []string `yaml:"audience"`
	Expires  uint     `yaml:"expires"`
}
