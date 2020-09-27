package types

type Config struct {
	Listen string     `yaml:"listen"`
	Cors   CorsOption `yaml:"cors"`
}
