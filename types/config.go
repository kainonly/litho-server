package types

type Config struct {
	Listen string                 `yaml:"listen"`
	App    AppOption              `yaml:"app"`
	Cors   CorsOption             `yaml:"cors"`
	Mysql  MysqlOption            `yaml:"mysql"`
	Redis  RedisOption            `yaml:"redis"`
	Token  map[string]TokenOption `yaml:"token"`
}
