package common

import (
	"fmt"
	"regexp"
	"strings"
)

type Values struct {
	Mode      string          `yaml:"mode"`
	Address   string          `yaml:"address"`
	Ip        string          `yaml:"ip"`
	Namespace string          `yaml:"namespace"`
	Key       string          `yaml:"key"`
	Domain    string          `yaml:"domain"`
	Cors      Cors            `yaml:"cors"`
	Database  Database        `yaml:"database"`
	Nats      Nats            `yaml:"nats"`
	Cos       Cos             `yaml:"cos"`
	Sms       Sms             `yaml:"sms"`
	Extend    Extend          `yaml:"extend"`
	Otlp      Otlp            `yaml:"otlp"`
	Resources ResourceCatalog `yaml:"-"`
}

var resourcePath = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type ResourceAction struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Resource struct {
	Name    string           `json:"name"`
	Key     string           `json:"key"`
	Actions []ResourceAction `json:"actions"`
}

type ResourceModule struct {
	Key     string            `yaml:"key"`
	Label   string            `yaml:"label"`
	Actions []string          `yaml:"actions"`
	Labels  map[string]string `yaml:"labels"`
}

type ResourceCatalog struct {
	Actions map[string]string `yaml:"actions"`
	Modules []ResourceModule  `yaml:"modules"`
	list    []Resource
	paths   map[string]struct{}
	allowed map[string]map[string]struct{}
}

func (x *ResourceCatalog) Prepare() error {
	if len(x.Modules) == 0 {
		return fmt.Errorf("resources.yml 缺少模块")
	}
	x.list = make([]Resource, 0, len(x.Modules))
	x.paths = make(map[string]struct{})
	x.allowed = make(map[string]map[string]struct{})
	for _, module := range x.Modules {
		if !resourcePath.MatchString(module.Key) || strings.TrimSpace(module.Label) == "" {
			return fmt.Errorf("resources.yml 模块 key 或 label 无效: %q", module.Key)
		}
		if _, exists := x.allowed[module.Key]; exists {
			return fmt.Errorf("resources.yml 模块重复: %s", module.Key)
		}
		if len(module.Actions) == 0 {
			return fmt.Errorf("resources.yml 模块 %s 缺少操作", module.Key)
		}
		x.allowed[module.Key] = make(map[string]struct{}, len(module.Actions))
		resource := Resource{Name: module.Label, Key: module.Key, Actions: make([]ResourceAction, 0, len(module.Actions))}
		for _, path := range module.Actions {
			if !resourcePath.MatchString(path) || path == "r" {
				return fmt.Errorf("resources.yml 模块 %s 的操作无效: %q", module.Key, path)
			}
			if _, exists := x.allowed[module.Key][path]; exists {
				return fmt.Errorf("resources.yml 操作重复: %s/%s", module.Key, path)
			}
			name := module.Labels[path]
			if name == "" {
				name = x.Actions[path]
			}
			if strings.TrimSpace(name) == "" {
				return fmt.Errorf("resources.yml 操作缺少名称: %s/%s", module.Key, path)
			}
			x.allowed[module.Key][path] = struct{}{}
			x.paths[module.Key+"/"+path] = struct{}{}
			resource.Actions = append(resource.Actions, ResourceAction{Name: name, Path: path})
		}
		for path := range module.Labels {
			if _, exists := x.allowed[module.Key][path]; !exists {
				return fmt.Errorf("resources.yml 存在未使用的名称: %s/%s", module.Key, path)
			}
		}
		x.list = append(x.list, resource)
	}
	return nil
}

func (x *ResourceCatalog) List() []Resource {
	return x.list
}

func (x *ResourceCatalog) ValidateRoutes(paths []string) error {
	actual := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		if _, exists := actual[path]; exists {
			return fmt.Errorf("重复注册的写入路由: %s", path)
		}
		actual[path] = struct{}{}
		if _, exists := x.paths[path]; !exists {
			return fmt.Errorf("resources.yml 缺少写入路由: %s", path)
		}
	}
	for path := range x.paths {
		if _, exists := actual[path]; !exists {
			return fmt.Errorf("resources.yml 中的操作没有对应写入路由: %s", path)
		}
	}
	return nil
}

// NormalizeStrategy checks resource permissions and stores a fully selected module as *.
func (x *ResourceCatalog) NormalizeStrategy(strategy map[string]any) (map[string]any, error) {
	result := make(map[string]any, len(strategy))
	for key, raw := range strategy {
		if key == "resources" {
			continue
		}
		allowed, exists := x.allowed[key]
		if !exists {
			return nil, fmt.Errorf("未知的资源模块: %s", key)
		}
		values, ok := raw.([]any)
		if !ok {
			return nil, fmt.Errorf("资源 %s 的操作必须是数组", key)
		}
		actions := make([]string, 0, len(values))
		seen := make(map[string]struct{}, len(values))
		for _, value := range values {
			path, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("资源 %s 的操作必须是字符串", key)
			}
			if _, exists := allowed[path]; !exists && path != "r" && path != "*" {
				return nil, fmt.Errorf("未知的资源操作: %s/%s", key, path)
			}
			if _, exists := seen[path]; !exists {
				actions = append(actions, path)
				seen[path] = struct{}{}
			}
		}
		if _, all := seen["*"]; all {
			result[key] = []string{"*"}
		} else if _, read := seen["r"]; read && len(seen) == len(allowed)+1 {
			result[key] = []string{"*"}
		} else if len(actions) > 0 {
			result[key] = actions
		}
	}
	return result, nil
}

type Cors struct {
	Origins  []string `yaml:"origins"`
	SameSite string   `yaml:"same_site,omitempty"` // Cookie SameSite 设置
}

type Database struct {
	Url        string  `yaml:"url"`
	Name       string  `yaml:"name"`
	Redis      string  `yaml:"redis"`
	Mongodb    Mongodb `yaml:"mongodb"`
	Clickhouse string  `yaml:"clickhouse"`
	Debug      bool    `yaml:"debug"`
}

type Mongodb struct {
	Url      string `yaml:"url"`
	Database string `yaml:"database"`
}

type Nats struct {
	Hosts []string `yaml:"hosts"`
	Token string   `yaml:"token"`
}

type Cos struct {
	Bucket    string `yaml:"bucket"`
	Files     string `yaml:"files"`
	Region    string `yaml:"region"`
	SecretId  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}

type Sms struct {
	SecretId  string            `yaml:"secret_id"`
	SecretKey string            `yaml:"secret_key"`
	Sign      string            `yaml:"sign"`
	AppId     string            `yaml:"app_id"`
	Templates map[string]string `yaml:"templates"`
}

type Otlp struct {
	Name        string `yaml:"name"`
	Endpoint    string `yaml:"endpoint"`
	Token       string `yaml:"token"`
	Environment string `yaml:"environment"`
}

type Extend struct {
	BaseUrl string `yaml:"base_url"`
}

func (x Values) IsRelease() bool {
	return x.Mode == "release"
}

func (x Values) IsSqlDebug() bool {
	return x.Database.Debug
}

// Name Redis Key 生成
func (x Values) Name(key string) string {
	return fmt.Sprintf("%s:%s", x.Namespace, key)
}
