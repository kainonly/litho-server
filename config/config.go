package config

type Config struct {
	App      App      `yaml:"app"`
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
}

type App struct {
	Debug bool   `yaml:"debug"`
	Name  string `yaml:"name"`
	Key   string `yaml:"key"`
	Lock  Lock   `yaml:"lock"`
}

type Lock struct {
	Limit        int64 `yaml:"limit"`
	RecoveryTime int64 `yaml:"recovery_time"`
}

type Database struct {
	Dsn             string `yaml:"dsn"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	TablePrefix     string `yaml:"table_prefix"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
