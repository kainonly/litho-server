module lab-api

go 1.13

require (
	github.com/gin-gonic/gin v1.6.3
	go.uber.org/fx v1.13.1
	gopkg.in/yaml.v2 v2.2.8
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.9
)

replace (
	github.com/kainonly/gin-helper => ./library/gin-helper
	github.com/kainonly/gin-planx => ./library/gin-planx
)
