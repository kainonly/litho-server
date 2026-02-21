package common

import (
    "fmt"
    "time"
)

type Values struct {
    App      AppValue      `yaml:"app"`
    Network  NetworkValue  `yaml:"network"`
    Cors     CorsValue     `yaml:"cors"`
    Database DatabaseValue `yaml:"database"`
    Nats     NatsValue     `yaml:"nats"`
    Dynamic  DynamicValues `yaml:"-"`
}

type AppValue struct {
    Mode      string `yaml:"mode"`
    Namespace string `yaml:"namespace"`
    Address   string `yaml:"address"`
    Key       string `yaml:"key"`
}

type NetworkValue struct {
    IP     string `yaml:"ip"`
    Domain string `yaml:"domain"`
}

type CorsValue struct {
    Enabled          bool          `yaml:"enabled"`
    Origins          []string      `yaml:"origins"`
    Methods          []string      `yaml:"methods"`
    Headers          []string      `yaml:"headers"`
    ExposeHeaders    []string      `yaml:"expose_headers"`
    AllowCredentials bool          `yaml:"allow_credentials"`
    MaxAge           time.Duration `yaml:"max_age"`
}

type DatabaseValue struct {
    Debug bool      `yaml:"debug"`
    DSN   string    `yaml:"dsn"`
    Name  string    `yaml:"name"`
    Gorm  GormValue `yaml:"gorm"`
    Pool  PoolValue `yaml:"pool"`
    Redis string    `yaml:"redis"`
}

type GormValue struct {
    LogLevel               string `yaml:"log_level"` // silent / error / warn / info
    PrepareStmt            bool   `yaml:"prepare_stmt"`
    SkipDefaultTransaction bool   `yaml:"skip_default_transaction"`
}

type PoolValue struct {
    MaxIdleConns int           `yaml:"max_idle_conns"`
    MaxOpenConns int           `yaml:"max_open_conns"`
    ConnMaxLife  time.Duration `yaml:"conn_max_life"`
}

type NatsValue struct {
    Hosts []string `yaml:"hosts"`
    Token string   `yaml:"token"`
}

type OtlpValue struct {
    Enabled     bool              `yaml:"enabled"`
    ServiceName string            `yaml:"service_name"`
    Environment string            `yaml:"environment"`
    Endpoint    string            `yaml:"endpoint"`
    Protocol    string            `yaml:"protocol"` // grpc / http
    Insecure    bool              `yaml:"insecure"`
    Headers     map[string]string `yaml:"headers,omitempty"`
}

type DynamicValues struct {
    Storage StorageValue `json:"storage"`
    SMS     SMSValue     `json:"sms"`
}

type StorageValue struct {
    Provider string   `json:"provider"` // cos
    Cos      CosValue `json:"cos"`
}

type CosValue struct {
    Bucket    string `json:"bucket"`
    Files     string `json:"files"`
    Region    string `json:"region"`
    SecretId  string `json:"secret_id"`
    SecretKey string `json:"secret_key"`
}

type SMSValue struct {
    Provider string           `json:"provider"` // tencent
    Tencent  *TencentSmsValue `json:"tencent,omitempty"`
}

type TencentSmsValue struct {
    SecretId  string            `json:"secret_id"`
    SecretKey string            `json:"secret_key"`
    AppId     string            `json:"app_id"`
    Sign      string            `json:"sign"`
    Templates map[string]string `json:"templates"`
    Region    string            `json:"region,omitempty"`
}

func (x Values) IsRelease() bool {
    return x.App.Mode == "release"
}

func (x Values) IsSqlDebug() bool {
    return x.Database.Debug
}

func (x Values) Name(key string) string {
    return fmt.Sprintf("%s:%s", x.App.Namespace, key)
}

func (x Values) LogName(key string) string {
    if x.IsRelease() {
        return fmt.Sprintf(`%s_logs`, key)
    }
    return fmt.Sprintf(`%s_logs_dev`, key)
}
