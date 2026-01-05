#!/usr/bin/env python3
"""
Atlas 数据库迁移脚本
用法: python scripts/migrate.py [--env local|prod] <command>

配置文件: .env (不提交到 git)
"""

import os
import subprocess
import sys
from contextlib import contextmanager
from dataclasses import dataclass
from pathlib import Path

# ============================================================================
# 配置
# ============================================================================


@dataclass
class SSHTunnelConfig:
    """SSH 隧道配置"""

    ssh_host: str  # SSH 跳板机地址
    ssh_port: int  # SSH 端口
    ssh_user: str  # SSH 用户名
    ssh_key: str  # SSH 私钥路径
    remote_host: str  # 远程数据库地址
    remote_port: int  # 远程数据库端口
    local_port: int  # 本地映射端口


@dataclass
class EnvConfig:
    """环境配置"""

    atlas_env: str  # Atlas 环境名称
    database_url: str  # 数据库连接 URL
    ssh_tunnel: SSHTunnelConfig | None = None  # SSH 隧道配置（可选）


def load_env_file():
    """
    加载 .env 文件到环境变量

    支持的格式:
    - KEY=value
    - KEY="value with spaces"
    - # 注释行
    """
    # 查找项目根目录的 .env 文件
    script_dir = Path(__file__).parent
    project_root = script_dir.parent
    env_file = project_root / ".env"

    if not env_file.exists():
        return

    with open(env_file) as f:
        for line in f:
            line = line.strip()

            # 跳过空行和注释
            if not line or line.startswith("#"):
                continue

            # 解析 KEY=value
            if "=" in line:
                key, _, value = line.partition("=")
                key = key.strip()
                value = value.strip()

                # 移除引号
                if (value.startswith('"') and value.endswith('"')) or (
                    value.startswith("'") and value.endswith("'")
                ):
                    value = value[1:-1]

                # 只设置未定义的环境变量（允许系统环境变量覆盖）
                if key not in os.environ:
                    os.environ[key] = value


def get_env(key: str, default: str = "") -> str:
    """获取环境变量"""
    return os.environ.get(key, default)


def get_env_int(key: str, default: int = 0) -> int:
    """获取整数类型的环境变量"""
    value = os.environ.get(key, "")
    if value:
        try:
            return int(value)
        except ValueError:
            pass
    return default


def build_environments() -> dict[str, EnvConfig]:
    """从环境变量构建环境配置"""
    envs = {}

    # 本地环境
    local_url = get_env("LOCAL_DATABASE_URL")
    if local_url:
        envs["local"] = EnvConfig(
            atlas_env="local",
            database_url=local_url,
        )

    # 生产环境
    prod_url = get_env("PROD_DATABASE_URL")
    if prod_url:
        # 检查是否配置了 SSH 隧道
        ssh_tunnel = None
        ssh_host = get_env("PROD_SSH_HOST")

        if ssh_host:
            ssh_tunnel = SSHTunnelConfig(
                ssh_host=ssh_host,
                ssh_port=get_env_int("PROD_SSH_PORT", 22),
                ssh_user=get_env("PROD_SSH_USER", "root"),
                ssh_key=get_env("PROD_SSH_KEY", "~/.ssh/id_rsa"),
                remote_host=get_env("PROD_DB_HOST", "localhost"),
                remote_port=get_env_int("PROD_DB_PORT", 5432),
                local_port=get_env_int("PROD_LOCAL_PORT", 15432),
            )

        envs["prod"] = EnvConfig(
            atlas_env="prod",
            database_url=prod_url,
            ssh_tunnel=ssh_tunnel,
        )

    return envs


# 命令定义
COMMANDS = {
    "diff": {
        "desc": "生成迁移文件",
        "args": ["migrate", "diff"],
    },
    "apply": {
        "desc": "应用迁移",
        "args": ["migrate", "apply"],
    },
    "status": {
        "desc": "查看迁移状态",
        "args": ["migrate", "status"],
    },
    "inspect": {
        "desc": "查看数据库 schema",
        "args": ["schema", "inspect"],
    },
    "push": {
        "desc": "直接推送 schema（开发用）",
        "args": ["schema", "apply", "--auto-approve"],
    },
    "hash": {
        "desc": "重新生成哈希",
        "args": ["migrate", "hash"],
        "no_env": True,  # 不需要 --env 参数
    },
}


# ============================================================================
# SSH 隧道管理
# ============================================================================


@contextmanager
def ssh_tunnel(config: SSHTunnelConfig):
    """
    建立 SSH 隧道的上下文管理器

    使用 sshtunnel 库建立隧道，退出时自动关闭
    """
    try:
        from sshtunnel import SSHTunnelForwarder
    except ImportError:
        print("错误: 需要安装 sshtunnel 库")
        print("安装方法: pip install sshtunnel")
        sys.exit(1)

    ssh_key_path = os.path.expanduser(config.ssh_key)

    print(f"正在建立 SSH 隧道: {config.ssh_user}@{config.ssh_host}...")
    print(f"  远程: {config.remote_host}:{config.remote_port}")
    print(f"  本地: localhost:{config.local_port}")

    tunnel = SSHTunnelForwarder(
        (config.ssh_host, config.ssh_port),
        ssh_username=config.ssh_user,
        ssh_pkey=ssh_key_path,
        remote_bind_address=(config.remote_host, config.remote_port),
        local_bind_address=("localhost", config.local_port),
    )

    try:
        tunnel.start()
        print("SSH 隧道已建立 ✓")
        print()
        yield tunnel
    finally:
        tunnel.stop()
        print()
        print("SSH 隧道已关闭")


@contextmanager
def no_tunnel():
    """空的上下文管理器，用于不需要隧道的情况"""
    yield None


# ============================================================================
# 命令执行
# ============================================================================


def run_atlas(args: list[str], env_config: EnvConfig) -> int:
    """执行 Atlas 命令"""
    os.environ["DATABASE_URL"] = env_config.database_url

    cmd = ["atlas"] + args
    print(f"执行: {' '.join(cmd)}")
    print()

    try:
        result = subprocess.run(cmd, check=False)
        return result.returncode
    except FileNotFoundError:
        print("错误: 未找到 atlas 命令，请确保已安装 Atlas CLI")
        print("安装方法: curl -sSf https://atlasgo.sh | sh")
        return 1
    except KeyboardInterrupt:
        print("\n已取消")
        return 130


def print_usage(environments: dict[str, EnvConfig]):
    """打印使用说明"""
    print(f"用法: {sys.argv[0]} [--env <environment>] <command>")
    print()
    print("环境:")
    if environments:
        for name in environments:
            print(f"  {name}")
    else:
        print("  (未配置，请检查 .env 文件)")
    print()
    print("命令:")
    for name, info in COMMANDS.items():
        print(f"  {name:8} - {info['desc']}")
    print()
    print("示例:")
    print(f"  {sys.argv[0]} diff              # 本地环境生成迁移")
    print(f"  {sys.argv[0]} --env prod apply  # 生产环境应用迁移")
    print()
    print("配置文件: .env (项目根目录)")


def parse_args() -> tuple[str, str, list[str]]:
    """
    解析命令行参数

    返回: (环境名, 命令名, 额外参数)
    """
    args = sys.argv[1:]
    env_name = "local"
    command = None
    extra_args = []

    i = 0
    while i < len(args):
        arg = args[i]

        if arg == "--env" and i + 1 < len(args):
            env_name = args[i + 1]
            i += 2
            continue

        if arg in ("-h", "--help", "help"):
            return env_name, "help", []

        if command is None:
            command = arg
        else:
            extra_args.append(arg)

        i += 1

    if command is None:
        return env_name, "help", []

    return env_name, command, extra_args


def main():
    # 加载 .env 文件
    load_env_file()

    # 构建环境配置
    environments = build_environments()

    # 解析参数
    env_name, command, extra_args = parse_args()

    # 帮助命令
    if command == "help":
        print_usage(environments)
        sys.exit(0)

    # 验证环境
    if not environments:
        print("错误: 未找到任何环境配置")
        print()
        print("请在项目根目录创建 .env 文件，参考 .env.example")
        sys.exit(1)

    if env_name not in environments:
        print(f"错误: 未知环境 '{env_name}'")
        print(f"可用环境: {', '.join(environments.keys())}")
        sys.exit(1)

    # 验证命令
    if command not in COMMANDS:
        print(f"错误: 未知命令 '{command}'")
        print()
        print_usage(environments)
        sys.exit(1)

    env_config = environments[env_name]
    cmd_info = COMMANDS[command]

    # 构建 Atlas 命令参数
    args = cmd_info["args"].copy()
    if not cmd_info.get("no_env"):
        args.extend(["--env", env_config.atlas_env])
    args.extend(extra_args)

    # 选择是否使用 SSH 隧道
    tunnel_ctx = ssh_tunnel(env_config.ssh_tunnel) if env_config.ssh_tunnel else no_tunnel()

    # 执行命令
    with tunnel_ctx:
        exit_code = run_atlas(args, env_config)

    sys.exit(exit_code)


if __name__ == "__main__":
    main()
