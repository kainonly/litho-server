module api

go 1.15

replace github.com/weplanx/go v0.0.2 => ./library

require (
	github.com/fatih/color v1.13.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-contrib/requestid v0.0.3
	github.com/gin-contrib/zap v0.0.2
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/google/wire v0.5.0
	github.com/json-iterator/go v1.1.11
	github.com/nats-io/nats.go v1.13.1-0.20220308171302-2f2f6968e98d
	github.com/nats-io/nkeys v0.3.0
	github.com/qri-io/jsonschema v0.2.1
	github.com/speps/go-hashids/v2 v2.0.1
	github.com/spf13/cobra v1.4.0
	github.com/tencentyun/cos-go-sdk-v5 v0.7.33
	github.com/thoas/go-funk v0.9.1
	github.com/weplanx/go v0.0.2
	github.com/weplanx/openapi v0.1.2
	github.com/weplanx/transfer v1.2.0 // indirect
	go.mongodb.org/mongo-driver v1.8.1
	go.uber.org/zap v1.21.0
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
