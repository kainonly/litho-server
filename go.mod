module laboratory

go 1.15

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/google/uuid v1.3.0
	github.com/json-iterator/go v1.1.11
	github.com/lib/pq v1.10.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/weplanx/support v1.1.0
	go.uber.org/fx v1.14.2
	golang.org/x/tools v0.1.5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/postgres v1.1.1
	gorm.io/gorm v1.21.15
)

replace github.com/weplanx/support v1.1.0 => ./support
