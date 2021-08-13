package config

type App struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
	Lock Lock   `yaml:"lock"`
}

type Lock struct {
	Limit        int64 `yaml:"limit"`
	RecoveryTime int64 `yaml:"recovery_time"`
}
