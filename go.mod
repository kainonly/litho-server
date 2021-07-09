module lab-api

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/emirpasic/gods v1.12.0
	github.com/gin-gonic/gin v1.7.2
	github.com/go-redis/redis/v8 v8.11.0
	github.com/json-iterator/go v1.1.10
	github.com/kainonly/gin-bit v0.0.0-20210630045306-1050ca2c7fcd
	github.com/kainonly/gin-helper v0.0.0-20210630045306-6e8210fedda9
	github.com/speps/go-hashids/v2 v2.0.1 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
	go.uber.org/fx v1.13.1
	go.uber.org/multierr v1.5.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
)

replace (
	github.com/kainonly/gin-bit v0.0.0-20210630045306-1050ca2c7fcd => ./library/gin-bit
	github.com/kainonly/gin-helper v0.0.0-20210630045306-6e8210fedda9 => ./library/gin-helper
)
