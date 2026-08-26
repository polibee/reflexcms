# ReflexCMS 工程常用命令（Git Bash / 任意 POSIX shell）
# Windows 若未安装 make，可直接执行各目标中的命令。

.PHONY: infra down backend admin lint test build migrate seed contract-test

## 启动基础设施（PostgreSQL 16 + Redis 7，需 Docker）
infra:
	docker compose -f deploy/docker-compose.yml up -d

down:
	docker compose -f deploy/docker-compose.yml down

## 开发服务：backend 与 admin 各开一个终端分别运行 make backend / make admin
backend:
	cd backend && go run .

admin:
	cd admin && npm run dev

lint:
	cd backend && go vet ./...
	cd admin && npm run lint

test:
	cd backend && go test ./...
	cd admin && npm test

build:
	cd backend && go build ./...
	cd admin && npm run build

migrate:
	cd backend && go run . artisan migrate

seed:
	cd backend && go run . artisan db:seed

contract-test:
	node scripts/contract-test.mjs
