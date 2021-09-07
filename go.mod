module lab-api

go 1.17

require (
	github.com/caarlos0/env/v6 v6.6.2
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/wire v0.5.0
	github.com/json-iterator/go v1.1.11
	github.com/kainonly/go-bit v1.0.1-beta.0.20210813060118-d167d0efebac
	github.com/lib/pq v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.12
)

require (
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.6.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.8.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.0.6 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.7.0 // indirect
	github.com/jackc/pgx/v4 v4.11.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/speps/go-hashids/v2 v2.0.1 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/kainonly/go-bit v1.0.1-beta.0.20210813060118-d167d0efebac => ./library/go-bit
