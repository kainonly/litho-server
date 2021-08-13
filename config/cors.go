package config

type CorsOption struct {

	// Matches the request origin
	Origin []string `yaml:"origin"`

	// Matches the request method
	Method []string `yaml:"method"`

	// Sets the Access-Control-Allow-Headers response header
	AllowHeader []string `yaml:"allow_header"`

	// Sets the Access-Control-Expose-Headers response header
	ExposedHeader []string `yaml:"exposed_header"`

	// Sets the Access-Control-Max-Age response header
	MaxAge int `yaml:"max_age"`

	// Sets the Access-Control-Allow-Credentials header
	Credentials bool `yaml:"credentials"`
}
