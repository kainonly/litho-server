module laboratory

go 1.15

require (
	github.com/alexedwards/argon2id v0.0.0-20210511081203-7d35d68092b8
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/google/wire v0.5.0
	github.com/thoas/go-funk v0.9.1
	github.com/weplanx/support v1.1.0
	go.mongodb.org/mongo-driver v1.7.3
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/weplanx/support v1.1.0 => ./support
