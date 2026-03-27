#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="${PROJECT_NAME:-antblog}"
IMAGE_NAME="${IMAGE_NAME:-antblog-server}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
APP_PORT="${APP_PORT:-8080}"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-root123456}"
MYSQL_DATABASE="${MYSQL_DATABASE:-antblog}"
REDIS_PASSWORD="${REDIS_PASSWORD:-foofoo}"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SQL_INIT_FILE="${ROOT_DIR}/docs/sqls/sqls.sql"

if [[ ! -f "${SQL_INIT_FILE}" ]]; then
  echo "missing sql file: ${SQL_INIT_FILE}" >&2
  exit 1
fi

NETWORK_NAME="${PROJECT_NAME}-net"
MYSQL_CONTAINER="${PROJECT_NAME}-mysql"
REDIS_CONTAINER="${PROJECT_NAME}-redis"
APP_CONTAINER="${PROJECT_NAME}-server"
MYSQL_DATA_VOLUME="${PROJECT_NAME}-mysql-data"
REDIS_DATA_VOLUME="${PROJECT_NAME}-redis-data"
UPLOAD_VOLUME="${PROJECT_NAME}-uploads"

docker network create "${NETWORK_NAME}" >/dev/null 2>&1 || true
docker volume create "${MYSQL_DATA_VOLUME}" >/dev/null
docker volume create "${REDIS_DATA_VOLUME}" >/dev/null
docker volume create "${UPLOAD_VOLUME}" >/dev/null
docker rm -f "${MYSQL_CONTAINER}" "${REDIS_CONTAINER}" "${APP_CONTAINER}" >/dev/null 2>&1 || true

docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" "${ROOT_DIR}"

docker run -d \
  --name "${MYSQL_CONTAINER}" \
  --network "${NETWORK_NAME}" \
  --health-cmd="mysqladmin ping -h localhost -p${MYSQL_ROOT_PASSWORD}" \
  --health-interval=10s \
  --health-timeout=5s \
  --health-retries=20 \
  -e MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD}" \
  -e MYSQL_DATABASE="${MYSQL_DATABASE}" \
  -v "${MYSQL_DATA_VOLUME}:/var/lib/mysql" \
  -v "${SQL_INIT_FILE}:/docker-entrypoint-initdb.d/001_schema.sql:ro" \
  mysql:8.4 >/dev/null

docker run -d \
  --name "${REDIS_CONTAINER}" \
  --network "${NETWORK_NAME}" \
  --health-cmd="redis-cli -a ${REDIS_PASSWORD} ping" \
  --health-interval=10s \
  --health-timeout=5s \
  --health-retries=20 \
  -v "${REDIS_DATA_VOLUME}:/data" \
  redis:7-alpine redis-server --requirepass "${REDIS_PASSWORD}" >/dev/null

for i in $(seq 1 30); do
  mysql_state="$(docker inspect --format='{{.State.Health.Status}}' "${MYSQL_CONTAINER}" 2>/dev/null || true)"
  redis_state="$(docker inspect --format='{{.State.Health.Status}}' "${REDIS_CONTAINER}" 2>/dev/null || true)"
  if [[ "${mysql_state}" == "healthy" && "${redis_state}" == "healthy" ]]; then
    break
  fi
  sleep 2
done

DSN="root:${MYSQL_ROOT_PASSWORD}@tcp(${MYSQL_CONTAINER}:3306)/${MYSQL_DATABASE}?charset=utf8mb4&parseTime=True&loc=Local"

docker run -d \
  --name "${APP_CONTAINER}" \
  --network "${NETWORK_NAME}" \
  -p "${APP_PORT}:8080" \
  -e ANTBLOG_SERVER_HOST=0.0.0.0 \
  -e ANTBLOG_SERVER_PORT=8080 \
  -e ANTBLOG_DATABASE_DRIVER=mysql \
  -e ANTBLOG_DATABASE_DSN="${DSN}" \
  -e ANTBLOG_REDIS_ADDR="${REDIS_CONTAINER}:6379" \
  -e ANTBLOG_REDIS_PASSWORD="${REDIS_PASSWORD}" \
  -e ANTBLOG_UPLOAD_LOCAL_PATH=/app/uploads \
  -e ANTBLOG_UPLOAD_BASE_URL="http://localhost:${APP_PORT}/uploads" \
  -v "${UPLOAD_VOLUME}:/app/uploads" \
  "${IMAGE_NAME}:${IMAGE_TAG}" >/dev/null

docker ps --filter "name=${PROJECT_NAME}"
echo "http://localhost:${APP_PORT}"
