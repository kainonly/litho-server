module lab-api

go 1.15

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/google/uuid v1.3.0
	github.com/json-iterator/go v1.1.10
	github.com/kainonly/go-bit v1.0.1-beta.0.20210813060118-d167d0efebac
	github.com/lib/pq v1.6.0 // indirect
	go.uber.org/fx v1.13.1
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.11
)

replace github.com/kainonly/go-bit v1.0.1-beta.0.20210813060118-d167d0efebac => ./library/go-bit
