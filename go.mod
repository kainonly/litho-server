module lab-api

go 1.15

require (
	entgo.io/ent v0.9.1 // indirect
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/google/uuid v1.3.0
	github.com/json-iterator/go v1.1.11
	github.com/kainonly/go-bit v1.0.1-beta.0.20210813060118-d167d0efebac
	go.mongodb.org/mongo-driver v1.7.2
	go.uber.org/fx v1.14.2
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/gorm v1.21.12
)

replace github.com/kainonly/go-bit v1.0.1-beta.0.20210813060118-d167d0efebac => ./library/go-bit
