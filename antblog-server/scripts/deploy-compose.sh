#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="${PROJECT_NAME:-antblog}"
APP_PORT="${APP_PORT:-8080}"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-root123456}"
MYSQL_DATABASE="${MYSQL_DATABASE:-antblog}"
REDIS_PASSWORD="${REDIS_PASSWORD:-foofoo}"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SQL_INIT_FILE="${ROOT_DIR}/docs/sqls/sqls.sql"
COMPOSE_FILE="${ROOT_DIR}/scripts/docker-compose.generated.yml"

if [[ ! -f "${SQL_INIT_FILE}" ]]; then
  echo "missing sql file: ${SQL_INIT_FILE}" >&2
  exit 1
fi

cat > "${COMPOSE_FILE}" <<EOF
name: ${PROJECT_NAME}
services:
  mysql:
    image: mysql:8.4
    container_name: ${PROJECT_NAME}-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: ${MYSQL_DATABASE}
    volumes:
      - ${PROJECT_NAME}-mysql-data:/var/lib/mysql
      - ${SQL_INIT_FILE}:/docker-entrypoint-initdb.d/001_schema.sql:ro
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-p${MYSQL_ROOT_PASSWORD}"]
      interval: 10s
      timeout: 5s
      retries: 20

  redis:
    image: redis:7-alpine
    container_name: ${PROJECT_NAME}-redis
    restart: unless-stopped
    command: ["redis-server", "--requirepass", "${REDIS_PASSWORD}"]
    volumes:
      - ${PROJECT_NAME}-redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "${REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 5s
      retries: 20

  server:
    build:
      context: ${ROOT_DIR}
      dockerfile: Dockerfile
    image: ${PROJECT_NAME}-server:latest
    container_name: ${PROJECT_NAME}-server
    restart: unless-stopped
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    ports:
      - "${APP_PORT}:8080"
    environment:
      ANTBLOG_SERVER_HOST: 0.0.0.0
      ANTBLOG_SERVER_PORT: 8080
      ANTBLOG_DATABASE_DRIVER: mysql
      ANTBLOG_DATABASE_DSN: "root:${MYSQL_ROOT_PASSWORD}@tcp(${PROJECT_NAME}-mysql:3306)/${MYSQL_DATABASE}?charset=utf8mb4&parseTime=True&loc=Local"
      ANTBLOG_REDIS_ADDR: ${PROJECT_NAME}-redis:6379
      ANTBLOG_REDIS_PASSWORD: ${REDIS_PASSWORD}
      ANTBLOG_UPLOAD_LOCAL_PATH: /app/uploads
      ANTBLOG_UPLOAD_BASE_URL: "http://localhost:${APP_PORT}/uploads"
    volumes:
      - ${PROJECT_NAME}-uploads:/app/uploads

volumes:
  ${PROJECT_NAME}-mysql-data:
  ${PROJECT_NAME}-redis-data:
  ${PROJECT_NAME}-uploads:
EOF

docker compose -f "${COMPOSE_FILE}" up -d --build
docker compose -f "${COMPOSE_FILE}" ps
echo "http://localhost:${APP_PORT}"
