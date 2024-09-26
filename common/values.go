package common

type Values struct {
	Mode      string `env:"MODE" envDefault:"debug"`
	Hostname  string `env:"HOSTNAME"`
	Address   string `env:"ADDRESS"`
	Namespace string `env:"NAMESPACE"`
	Key       string `env:"KEY,required"`

	Database struct {
		Url   string `env:"URL,required"`
		Redis string `env:"REDIS,required"`
	} `envPrefix:"DATABASE_"`

	Nats struct {
		Hosts []string `env:"HOSTS,required" envSeparator:","`
		Pub   string   `env:"PUB,required"`
		Nkey  string   `env:"NKEY,required"`
	} `envPrefix:"NATS_"`

	Otlp struct {
		Enabled  *bool  `env:"ENABLED" envDefault:"true"`
		Endpoint string `env:"ENDPOINT"`
		Token    string `env:"TOKEN"`
	} `envPrefix:"OTLP_"`
}

func (x Values) IsRelease() bool {
	return x.Mode == "release"
}
