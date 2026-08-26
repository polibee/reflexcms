#!/usr/bin/env bash
# ReflexCMS PostgreSQL backup (plan §M5).
# Usage:
#   PG_BIN=/d/laragon/bin/postgresql/pgsql/bin ./deploy/backup.sh
# Environment:
#   PG_BIN      psql/pg_dump directory   (default auto-detect Laragon)
#   DB_HOST/DB_PORT/DB_DATABASE/DB_USERNAME/PGPASSWORD   same as backend/.env
#   BACKUP_DIR  output directory         (default ../backups relative to repo)
#   RETAIN_DAYS delete dumps older than N days (default 14)
set -euo pipefail

REPO_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PG_BIN="${PG_BIN:-/d/laragon/bin/postgresql/pgsql/bin}"
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-5433}"
DB_DATABASE="${DB_DATABASE:-reflexcms}"
DB_USERNAME="${DB_USERNAME:-reflexcms}"
export PGPASSWORD="${PGPASSWORD:-secret}"

BACKUP_DIR="${BACKUP_DIR:-$REPO_DIR/backups}"
RETAIN_DAYS="${RETAIN_DAYS:-14}"
STAMP="$(date +%Y%m%d_%H%M%S)"
TARGET="$BACKUP_DIR/${DB_DATABASE}_${STAMP}.dump"

mkdir -p "$BACKUP_DIR"

echo "==> pg_dump -> $TARGET"
"$PG_BIN/pg_dump" \
  --host "$DB_HOST" --port "$DB_PORT" \
  --username "$DB_USERNAME" \
  --format custom --compress 6 \
  --file "$TARGET" "$DB_DATABASE"

SIZE=$(du -h "$TARGET" | cut -f1)
echo "==> done ($SIZE)"

echo "==> pruning dumps older than ${RETAIN_DAYS} days"
find "$BACKUP_DIR" -name "${DB_DATABASE}_*.dump" -mtime +"$RETAIN_DAYS" -print -delete

echo "==> current backups:"
ls -lh "$BACKUP_DIR" || true
