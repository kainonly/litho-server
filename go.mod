module lab-api

go 1.13

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kainonly/gin-helper v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.5.1 // indirect
	go.uber.org/fx v1.13.1
	go.uber.org/multierr v1.5.0 // indirect
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
	golang.org/x/text v0.3.3 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
)

replace (
	github.com/kainonly/gin-helper => ./library/gin-helper
	github.com/kainonly/gin-planx => ./library/gin-planx
)
