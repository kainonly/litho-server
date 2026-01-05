// Atlas 配置文件
// 定义数据库迁移的环境和数据源

// 使用 GORM 模型作为 schema 源
data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas-loader",
  ]
}

// 本地开发环境
env "local" {
  // 使用 GORM 模型生成的 schema
  src = data.external_schema.gorm.url

  // 目标数据库 URL（可通过环境变量覆盖）
  url = getenv("DATABASE_URL")

  // 迁移文件目录
  migration {
    dir = "file://migrations"
  }

  // 开发数据库（用于计算 schema diff）
  dev = "docker+postgres://percona/percona-distribution-postgresql:18.1/dev?search_path=public"
}

// 生产环境
env "prod" {
  src = data.external_schema.gorm.url
  url = getenv("DATABASE_URL")

  migration {
    dir = "file://migrations"
  }

  // 生产环境也需要一个干净的数据库来计算 diff
  dev = "docker+postgres://percona/percona-distribution-postgresql:18.1/dev?search_path=public"
}
