module lab-api

go 1.13

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/go-playground/validator/v10 v10.6.1 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/kainonly/gin-helper v0.0.1-beta // indirect
	github.com/kainonly/gin-planx v0.0.0-20210629081218-4978ca8748fc
	github.com/stretchr/testify v1.5.1 // indirect
	go.uber.org/fx v1.13.1
	go.uber.org/multierr v1.5.0 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
	golang.org/x/text v0.3.3 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
)

replace (
	github.com/kainonly/gin-helper => ./library/gin-helper
	github.com/kainonly/gin-planx => ./library/gin-planx
)
