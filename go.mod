module lab-api

go 1.15

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kainonly/go-bit v1.1.0
	github.com/lib/pq v1.10.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/fx v1.14.2
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/tools v0.1.5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.12
)

replace github.com/kainonly/go-bit v1.1.0 => ./library/go-bit
