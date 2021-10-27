module laboratory

go 1.15

require (
	github.com/alexedwards/argon2id v0.0.0-20210511081203-7d35d68092b8
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/thoas/go-funk v0.9.1
	github.com/weplanx/support v1.1.0
	go.mongodb.org/mongo-driver v1.7.3
	go.uber.org/fx v1.14.2
	golang.org/x/tools v0.1.5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/weplanx/support v1.1.0 => ./support
