param(
  [string]$ProjectName = "antblog",
  [int]$AppPort = 8080,
  [string]$MySQLRootPassword = "root123456",
  [string]$MySQLDatabase = "antblog",
  [string]$RedisPassword = "foofoo"
)

$ErrorActionPreference = "Stop"

$projectRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$sqlInitFile = (Resolve-Path (Join-Path $projectRoot "docs\sqls\sqls.sql")).Path
$dockerSqlInitFile = $sqlInitFile -replace "\\", "/"
$composeFile = Join-Path $PSScriptRoot "docker-compose.generated.yml"

$composeYaml = @"
name: $ProjectName
services:
  mysql:
    image: mysql:8.4
    container_name: ${ProjectName}-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: $MySQLRootPassword
      MYSQL_DATABASE: $MySQLDatabase
    volumes:
      - ${ProjectName}-mysql-data:/var/lib/mysql
      - $dockerSqlInitFile:/docker-entrypoint-initdb.d/001_schema.sql:ro
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-p$MySQLRootPassword"]
      interval: 10s
      timeout: 5s
      retries: 20

  redis:
    image: redis:7-alpine
    container_name: ${ProjectName}-redis
    restart: unless-stopped
    command: ["redis-server", "--requirepass", "$RedisPassword"]
    volumes:
      - ${ProjectName}-redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "$RedisPassword", "ping"]
      interval: 10s
      timeout: 5s
      retries: 20

  server:
    build:
      context: $projectRoot
      dockerfile: Dockerfile
    image: ${ProjectName}-server:latest
    container_name: ${ProjectName}-server
    restart: unless-stopped
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    ports:
      - "$AppPort:8080"
    environment:
      ANTBLOG_SERVER_HOST: 0.0.0.0
      ANTBLOG_SERVER_PORT: 8080
      ANTBLOG_DATABASE_DRIVER: mysql
      ANTBLOG_DATABASE_DSN: root:$MySQLRootPassword@tcp(${ProjectName}-mysql:3306)/$MySQLDatabase?charset=utf8mb4&parseTime=True&loc=Local
      ANTBLOG_REDIS_ADDR: ${ProjectName}-redis:6379
      ANTBLOG_REDIS_PASSWORD: $RedisPassword
      ANTBLOG_UPLOAD_LOCAL_PATH: /app/uploads
      ANTBLOG_UPLOAD_BASE_URL: http://localhost:$AppPort/uploads
    volumes:
      - ${ProjectName}-uploads:/app/uploads

volumes:
  ${ProjectName}-mysql-data:
  ${ProjectName}-redis-data:
  ${ProjectName}-uploads:
"@

Set-Content -Path $composeFile -Value $composeYaml -Encoding UTF8
docker compose -f $composeFile up -d --build
docker compose -f $composeFile ps
Write-Output "http://localhost:$AppPort"
