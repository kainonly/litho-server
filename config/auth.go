package config

type Lock struct {
	Limit        int64 `yaml:"limit"`
	RecoveryTime int64 `yaml:"recovery_time"`
}
