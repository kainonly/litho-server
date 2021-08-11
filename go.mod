module lab-api

go 1.13

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/go-redis/redis/v8 v8.11.0
	github.com/google/wire v0.5.0
	github.com/json-iterator/go v1.1.10
	github.com/kainonly/gin-bit v0.0.0-20210630045306-1050ca2c7fcd
	github.com/kainonly/gin-helper v0.0.0-20210630045306-6e8210fedda9
	github.com/lib/pq v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.11
)

replace (
	github.com/kainonly/gin-bit v0.0.0-20210630045306-1050ca2c7fcd => ./library/gin-bit
	github.com/kainonly/gin-helper v0.0.0-20210630045306-6e8210fedda9 => ./library/gin-helper
)
