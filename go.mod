module api

go 1.15

replace github.com/weplanx/go v0.0.0-20211109121132-7a8d66264652 => ./library

require (
	github.com/apache/pulsar-client-go v0.7.1-0.20220210221528-5daa17b02bff
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-contrib/requestid v0.0.3
	github.com/gin-contrib/zap v0.0.2
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/validator/v10 v10.6.1
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/google/wire v0.5.0
	github.com/json-iterator/go v1.1.11
	github.com/speps/go-hashids/v2 v2.0.1
	github.com/tencentyun/cos-go-sdk-v5 v0.7.33
	github.com/thoas/go-funk v0.9.1
	github.com/weplanx/go v0.0.0-20211109121132-7a8d66264652
	go.mongodb.org/mongo-driver v1.8.1
	go.uber.org/zap v1.21.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
