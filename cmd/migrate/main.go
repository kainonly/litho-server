package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config 配置结构
type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	DSN       string           `yaml:"dsn"`
	SSHTunnel *SSHTunnelConfig `yaml:"ssh_tunnel,omitempty"`
}

type SSHTunnelConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Key        string `yaml:"key"`
	RemoteHost string `yaml:"remote_host"`
	RemotePort int    `yaml:"remote_port"`
	LocalPort  int    `yaml:"local_port"`
}

// 命令定义
type Command struct {
	Desc  string
	Args  []string
	NoEnv bool
}

var commands = map[string]Command{
	"diff": {
		Desc: "生成迁移文件",
		Args: []string{"migrate", "diff"},
	},
	"apply": {
		Desc: "应用迁移",
		Args: []string{"migrate", "apply"},
	},
	"status": {
		Desc: "查看迁移状态",
		Args: []string{"migrate", "status"},
	},
	"inspect": {
		Desc: "查看数据库 schema",
		Args: []string{"schema", "inspect"},
	},
	"push": {
		Desc: "直接推送 schema（开发用）",
		Args: []string{"schema", "apply", "--auto-approve"},
	},
	"hash": {
		Desc:  "重新生成哈希",
		Args:  []string{"migrate", "hash"},
		NoEnv: true,
	},
}

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/values.yml", "配置文件路径")
	help := flag.Bool("help", false, "显示帮助")
	flag.Parse()

	if *help || flag.NArg() == 0 {
		printUsage()
		return
	}

	cmdName := flag.Arg(0)
	extraArgs := flag.Args()[1:]

	// 验证命令
	cmd, ok := commands[cmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "错误: 未知命令 '%s'\n\n", cmdName)
		printUsage()
		os.Exit(1)
	}

	// 加载配置
	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 加载配置失败: %v\n", err)
		fmt.Fprintln(os.Stderr, "\n请确保配置文件存在且包含 database.dsn 配置")
		os.Exit(1)
	}

	// 转换 DSN 为 URL
	dbURL, err := dsnToURL(config.Database.DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 解析 DSN 失败: %v\n", err)
		os.Exit(1)
	}

	// 构建 Atlas 命令参数
	args := append([]string{}, cmd.Args...)
	if !cmd.NoEnv {
		args = append(args, "--env", "local")
	}
	args = append(args, extraArgs...)

	// 执行命令
	exitCode := runAtlas(args, dbURL)
	os.Exit(exitCode)
}

func loadConfig(path string) (*Config, error) {
	// 获取项目根目录
	execPath, err := os.Executable()
	if err != nil {
		execPath, _ = os.Getwd()
	}

	// 尝试多个可能的根目录
	var configFile string
	candidates := []string{
		filepath.Join(filepath.Dir(execPath), "..", "..", path),
		filepath.Join(filepath.Dir(execPath), path),
		path,
	}

	for _, candidate := range candidates {
		absPath, _ := filepath.Abs(candidate)
		if _, err := os.Stat(absPath); err == nil {
			configFile = absPath
			break
		}
	}

	if configFile == "" {
		return nil, fmt.Errorf("配置文件不存在: %s", path)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Database.DSN == "" {
		return nil, fmt.Errorf("配置文件缺少 database.dsn")
	}

	return &config, nil
}

func dsnToURL(dsn string) (string, error) {
	// 解析 GORM DSN 格式: host=xxx user=xxx password=xxx dbname=xxx port=xxx ...
	params := make(map[string]string)
	re := regexp.MustCompile(`(\w+)=([^\s]+|"[^"]*")`)
	matches := re.FindAllStringSubmatch(dsn, -1)

	for _, match := range matches {
		key := match[1]
		value := strings.Trim(match[2], `"`)
		params[key] = value
	}

	host := getOrDefault(params, "host", "localhost")
	port := getOrDefault(params, "port", "5432")
	user := getOrDefault(params, "user", "postgres")
	password := params["password"]
	dbname := getOrDefault(params, "dbname", "postgres")
	sslmode := getOrDefault(params, "sslmode", "disable")

	var userInfo string
	if password != "" {
		userInfo = fmt.Sprintf("%s:%s", user, url.QueryEscape(password))
	} else {
		userInfo = user
	}

	return fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s", userInfo, host, port, dbname, sslmode), nil
}

func getOrDefault(m map[string]string, key, defaultVal string) string {
	if v, ok := m[key]; ok && v != "" {
		return v
	}
	return defaultVal
}

func runAtlas(args []string, dbURL string) int {
	cmd := exec.Command("atlas", args...)
	cmd.Env = append(os.Environ(), "DATABASE_URL="+dbURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("执行: atlas %s\n\n", strings.Join(args, " "))

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		if strings.Contains(err.Error(), "executable file not found") {
			fmt.Fprintln(os.Stderr, "错误: 未找到 atlas 命令，请确保已安装 Atlas CLI")
			fmt.Fprintln(os.Stderr, "安装方法: curl -sSf https://atlasgo.sh | sh")
		} else {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		}
		return 1
	}
	return 0
}

func printUsage() {
	fmt.Println("Atlas 数据库迁移工具")
	fmt.Println()
	fmt.Println("用法: migrate [选项] <命令> [参数...]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -config string  配置文件路径 (默认: config/values.yml)")
	fmt.Println("  -help           显示帮助")
	fmt.Println()
	fmt.Println("命令:")
	for name, cmd := range commands {
		fmt.Printf("  %-8s  %s\n", name, cmd.Desc)
	}
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  migrate diff                           # 生成迁移")
	fmt.Println("  migrate apply                          # 应用迁移")
	fmt.Println("  migrate -config config/prod.yml apply  # 使用指定配置")
}
