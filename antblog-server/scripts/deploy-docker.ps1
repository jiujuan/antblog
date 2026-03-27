param(
  [string]$ProjectName = "antblog",
  [string]$ImageName = "antblog-server",
  [string]$ImageTag = "latest",
  [int]$AppPort = 8080,
  [string]$MySQLRootPassword = "root123456",
  [string]$MySQLDatabase = "antblog",
  [string]$RedisPassword = "foofoo"
)

$ErrorActionPreference = "Stop"

$projectRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$sqlInitFile = (Resolve-Path (Join-Path $projectRoot "docs\sqls\sqls.sql")).Path
$dockerSqlInitFile = $sqlInitFile -replace "\\", "/"
$networkName = "$ProjectName-net"
$mysqlContainer = "$ProjectName-mysql"
$redisContainer = "$ProjectName-redis"
$appContainer = "$ProjectName-server"
$mysqlDataVolume = "$ProjectName-mysql-data"
$redisDataVolume = "$ProjectName-redis-data"
$uploadVolume = "$ProjectName-uploads"

$null = docker network create $networkName 2>$null
$null = docker volume create $mysqlDataVolume
$null = docker volume create $redisDataVolume
$null = docker volume create $uploadVolume

docker rm -f $mysqlContainer $redisContainer $appContainer 2>$null | Out-Null

docker build -t "${ImageName}:${ImageTag}" $projectRoot

docker run -d `
  --name $mysqlContainer `
  --network $networkName `
  -e MYSQL_ROOT_PASSWORD=$MySQLRootPassword `
  -e MYSQL_DATABASE=$MySQLDatabase `
  -v "${mysqlDataVolume}:/var/lib/mysql" `
  -v "${dockerSqlInitFile}:/docker-entrypoint-initdb.d/001_schema.sql:ro" `
  mysql:8.4

docker run -d `
  --name $redisContainer `
  --network $networkName `
  -v "${redisDataVolume}:/data" `
  redis:7-alpine redis-server --requirepass $RedisPassword

Start-Sleep -Seconds 10

$dsn = "root:$MySQLRootPassword@tcp($mysqlContainer:3306)/$MySQLDatabase?charset=utf8mb4&parseTime=True&loc=Local"

docker run -d `
  --name $appContainer `
  --network $networkName `
  -p "${AppPort}:8080" `
  -e ANTBLOG_SERVER_HOST=0.0.0.0 `
  -e ANTBLOG_SERVER_PORT=8080 `
  -e ANTBLOG_DATABASE_DRIVER=mysql `
  -e ANTBLOG_DATABASE_DSN="$dsn" `
  -e ANTBLOG_REDIS_ADDR="${redisContainer}:6379" `
  -e ANTBLOG_REDIS_PASSWORD="$RedisPassword" `
  -e ANTBLOG_UPLOAD_LOCAL_PATH=/app/uploads `
  -e ANTBLOG_UPLOAD_BASE_URL="http://localhost:$AppPort/uploads" `
  -v "${uploadVolume}:/app/uploads" `
  "${ImageName}:${ImageTag}"

docker ps --filter "name=$ProjectName"
Write-Output "http://localhost:$AppPort"
