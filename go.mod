module lab-api

go 1.15

require (
	github.com/gin-gonic/gin v1.6.3
	go.uber.org/fx v1.13.1
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.9
)

replace (
	github.com/kainonly/gin-curd v0.0.0-20201104073226-f975793ac46d => ../../kainonly/gin-curd
	github.com/kainonly/gin-extra v0.0.0-20201106142436-122606ab7729 => ../../kainonly/gin-extra
)
