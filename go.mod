module lab-api

go 1.15

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/emirpasic/gods v1.12.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v8 v8.3.3
	github.com/json-iterator/go v1.1.10
	github.com/kainonly/gin-curd v0.0.0-20201104073226-f975793ac46d
	github.com/kainonly/gin-extra v0.0.0-20201106142436-122606ab7729
	go.uber.org/fx v1.13.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	gorm.io/driver/postgres v1.0.6-0.20201120082907-566aa2e6ed74
	gorm.io/gorm v1.20.9-0.20201207023106-e1952924e2a8
)

replace (
	github.com/kainonly/gin-curd v0.0.0-20201104073226-f975793ac46d => ../gin-curd
	github.com/kainonly/gin-extra v0.0.0-20201106142436-122606ab7729 => ../gin-extra
)
