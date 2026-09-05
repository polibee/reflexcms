#!/usr/bin/env bash
# =============================================================================
# ReflexCMS 一键部署脚本 / One-click deployment
# 用法: ./deploy/deploy.sh          (Linux/macOS, 需 Docker + Docker Compose)
# =============================================================================
set -euo pipefail

cd "$(dirname "$0")/.."
ROOT_DIR=$(pwd)

echo "=========================================="
echo " ReflexCMS 一键部署 / One-click deploy"
echo "=========================================="

# ---------- 1. 环境检查 ----------
command -v docker >/dev/null 2>&1 || { echo "✗ 未安装 Docker"; exit 1; }
docker compose version >/dev/null 2>&1 || { echo "✗ 未安装 Docker Compose v2"; exit 1; }
command -v go >/dev/null 2>&1 && GO_OK=1 || GO_OK=0
command -v node >/dev/null 2>&1 && NODE_OK=1 || NODE_OK=0
echo "✓ Docker 就绪"

# ---------- 2. 环境变量 ----------
if [ ! -f backend/.env ]; then
  cp backend/.env.example backend/.env
  echo "✓ 已生成 backend/.env（请按需修改数据库连接）"
fi
if [ ! -f admin/.env ]; then
  cat > admin/.env <<'EOF'
NUXT_PUBLIC_API_BASE=http://127.0.0.1:9000
cookieSecure=false
EOF
  echo "✓ 已生成 admin/.env"
fi

# ---------- 3. 启动基础设施 ----------
echo "→ 启动 PostgreSQL + Redis..."
docker compose -f deploy/docker-compose.yml up -d
sleep 5

# ---------- 4. 构建后端 ----------
echo "→ 构建后端..."
if [ "$GO_OK" = "1" ]; then
  (cd backend && go build -o tmp/main.exe .)
  (cd backend && ./tmp/main.exe artisan migrate) || true
  (cd backend && ./tmp/main.exe artisan db:seed) || true
  # Super-admin bring-up: on a fresh install this generates a strong random
  # password and prints it exactly once. Existing super-admins are untouched.
  echo "→ 配置超级管理员账号..."
  (cd backend && ./tmp/main.exe artisan admin:bootstrap) || true
  echo "✓ 后端构建完成（生产建议: go build -o reflexcms-api 且用 systemd/容器托管）"
else
  echo "! 本机无 Go，跳过本地构建——请使用 Dockerfile 构建镜像"
fi

# ---------- 5. 构建前端 ----------
echo "→ 构建前端..."
if [ "$NODE_OK" = "1" ]; then
  (cd admin && npm install --no-audit --no-fund && npm run build)
  echo "✓ 前端构建完成: admin/.output/ （Node 服务: node admin/.output/server/index.mjs）"
else
  echo "! 本机无 Node，跳过前端构建"
fi

echo ""
echo "=========================================="
echo " 部署完成！"
echo "  后端 API : http://127.0.0.1:9000/healthz"
echo "  前端站点 : http://localhost:3000"
echo "  管理后台 : http://localhost:3000/admin"
echo "  默认账号 : admin@reflexcms.dev / ReflexCMS@2026"
echo "=========================================="
