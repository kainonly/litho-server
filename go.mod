module lab-api

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.2
	github.com/json-iterator/go v1.1.10
	github.com/kainonly/gin-bit v0.0.0-20210630045306-1050ca2c7fcd
	github.com/kainonly/gin-helper v0.0.0-20210630045306-6e8210fedda9
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
	github.com/kainonly/gin-bit v0.0.0-20210630045306-1050ca2c7fcd => ./library/gin-bit
	github.com/kainonly/gin-helper v0.0.0-20210630045306-6e8210fedda9 => ./library/gin-helper
)
