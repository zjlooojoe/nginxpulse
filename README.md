<p align="center">
  <img src="docs/brand-mark.svg" alt="NginxPulse Logo" width="120" height="120">
</p>

# NginxPulse

轻量级 Nginx 访问日志分析与可视化面板，提供实时统计、PV 过滤、IP 归属地与客户端解析。

![demo-img-1.png](docs/demo-img-1.png)

![demo-img-2.png](docs/demo-img-2.png)
## 目录
- [项目开发技术栈](#项目开发技术栈)
- [IP 归属地查询策略](#ip-归属地查询策略)
- [如何使用项目](#如何使用项目)
- [二次开发注意事项](#二次开发注意事项)
- [目录结构与主要文件](#目录结构与主要文件)

## 项目开发技术栈
- **后端**：`Go 1.23.x` · `Gin` · `Logrus`
- **数据**：`SQLite (modernc.org/sqlite)`
- **IP 归属地**：`ip2region`（本地库） + `ip-api.com`（远程批量）
- **前端**：`Vue 3` · `Vite` · `TypeScript` · `PrimeVue` · `ECharts/Chart.js` · `Scss`
- **容器**：`Docker / Docker Compose` · `Nginx`（前端静态部署）

## IP 归属地查询策略
1. **快速过滤**：空值/本地/回环地址返回“本地”，内网地址返回“内网/本地网络”。
2. **缓存优先**：内存缓存命中直接返回（最多缓存 50,000 条）。
3. **远程优先**：调用 `ip-api.com/batch` 批量查询，超时 1.2s，单批最多 100 个。
4. **本地兜底**：远程失败或结果为“未知”时，IPv4 使用内置 ip2region 数据库本地查询（50ms 超时）。
5. **IPv6 处理**：仅走远程查询，远程失败则返回“未知”。

> 本地数据库 `ip2region.xdb` 内嵌在二进制中，首次启动会自动解压到 `./var/nginxpulse_data/ip2region.xdb`，并尝试加载向量索引提升查询性能。

> 本项目会访问外网 IP 归属地 API（`ip-api.com`），部署环境需放行该域名的出站访问。

## 如何使用项目

### 1) Docker
单镜像（前端 Nginx + 后端服务）：

使用远程镜像（Docker Hub）：

```bash
docker run -d --name nginxpulse \
  -p 8088:8088 \
  -p 8089:8089 \
  -e WEBSITES='[{"name":"主站","logPath":"/var/log/nginx/access.log","domains":["example.com","www.example.com"]}]' \
  -e SERVER_PORT=":8089" \
  -v /var/log/nginx:/var/log/nginx:ro \
  -v "$(pwd)/var/nginxpulse_data:/app/var/nginxpulse_data" \
  magiccoders/nginxpulse:latest
```

本地构建运行：

```bash
docker build -t nginxpulse:local .
docker run -d --name nginxpulse \
  -p 8088:8088 \
  -p 8089:8089 \
  -e WEBSITES='[{"name":"主站","logPath":"/var/log/nginx/access.log","domains":["example.com","www.example.com"]}]' \
  -v /var/log/nginx:/var/log/nginx:ro \
  -v "$(pwd)/var/nginxpulse_data:/app/var/nginxpulse_data" \
  nginxpulse:local
```

> 如果更偏好配置文件方式，可将 `configs/nginxpulse_config.json` 挂载到容器内的 `/app/configs/nginxpulse_config.json`。

### 2) Docker Compose
使用远程镜像（Docker Hub）：将 `docker-compose.yml` 改为下方远程镜像版本，然后执行：

```bash
docker compose up -d
```

本地构建运行（基于源码构建镜像）：保持仓库自带的 `docker-compose.yml`，执行：

```bash
docker compose up -d --build
```

示例 `docker-compose.yml`（远程镜像）：

```yml
version: "3.8"
services:
  nginxpulse:
    image: magiccoders/nginxpulse:latest
    container_name: nginxpulse
    ports:
      - "8088:8088"
      - "8089:8089"
    environment:
      WEBSITES: '[{"name":"主站","logPath":"/var/log/nginx/access.log","domains":["example.com","www.example.com"]}]'
    volumes:
      - /var/log/nginx:/var/log/nginx:ro
      - ./var/nginxpulse_data:/app/var/nginxpulse_data
      - /etc/localtime:/etc/localtime:ro
    restart: unless-stopped
```

示例 `docker-compose.yml`（本地构建）：

```yml
version: "3.8"
services:
  nginxpulse:
    image: nginxpulse:local
    build:
      context: .
    container_name: nginxpulse
    ports:
      - "8088:8088"
      - "8089:8089"
    environment:
      WEBSITES: '[{"name":"主站","logPath":"/var/log/nginx/access.log","domains":["example.com","www.example.com"]}]'
    volumes:
      - /var/log/nginx:/var/log/nginx:ro
      - ./var/nginxpulse_data:/app/var/nginxpulse_data
      - /etc/localtime:/etc/localtime:ro
    restart: unless-stopped
```

说明：
- `logPath` 必须是容器内路径，确保与挂载目录一致。
- `var/nginxpulse_data` 挂载用于持久化数据库和解析缓存，推荐保留。

参数说明（环境变量）：
- `WEBSITES`（必填，无配置文件时）
  - 网站列表 JSON 数组，字段：`name`、`logPath`、`domains`（可选）。
  - `domains` 用于将 referer 归类为“站内访问”，不影响日志解析与 PV 过滤。
- `CONFIG_JSON`（可选）
  - 完整配置 JSON 字符串（等同于 `configs/nginxpulse_config.json` 内容）。
  - 设置后会忽略本地配置文件，其他环境变量仍可覆盖其中字段。
- `LOG_DEST`（可选，默认：`file`）
  - 日志输出位置：`file` 或 `stdout`。
- `TASK_INTERVAL`（可选，默认：`1m`）
  - 扫描间隔，支持 `5m`、`25s` 等 Go duration 格式。
- `SERVER_PORT`（可选，默认：`:8089`）
  - 服务监听地址，可传 `:8089` 或 `8089`，不带冒号会自动补上。
- `PV_STATUS_CODES`（可选，默认：`[200]`）
  - 统计 PV 的状态码列表，可用 JSON 数组或逗号分隔值。
- `PV_EXCLUDE_PATTERNS`（可选，默认内置规则）
  - 全局 URL 排除正则数组（JSON 数组）。
- `PV_EXCLUDE_IPS`（可选，默认：空或配置文件）
  - 排除 IP 列表（JSON 数组或逗号分隔）。

访问：
- 前端：`http://localhost:8088`
- 后端：`http://localhost:8089`

> PV_EXCLUDE_PATTERNS和PV_EXCLUDE_IPS的具体格式请参考[nginxpulse_config.json](configs/nginxpulse_config.json)

### 3) 手动构建（前端、后端）
前端构建：

```bash
cd webapp
npm install
npm run build
```

后端构建：

```bash
go mod download
go build -o bin/nginxpulse ./cmd/nginxpulse/main.go
```

本地开发（前后端一起跑）：

```bash
./scripts/dev_local.sh
```

> 前端开发服务默认端口 8088，并会将 `/api` 代理到 `http://127.0.0.1:8089`。
> 本地开发前请准备好日志文件，放在 `var/log/` 下（或确保 `configs/nginxpulse_config.json` 的 `logPath` 指向对应文件）。

## 二次开发注意事项

### 环境依赖
- Go 1.23.x（与 `go.mod` 保持一致）
- Node.js 20+ / npm
- Docker（可选，用于容器化）

### 配置与数据目录
- 配置文件：`configs/nginxpulse_config.json`
- 数据目录：`var/nginxpulse_data/`
  - `nginxpulse.db`：SQLite 数据库
  - `nginx_scan_state.json`：日志扫描游标
  - `ip2region.xdb`：IP 本地库
- 环境变量覆盖：
  - `CONFIG_JSON` / `WEBSITES`
  - `LOG_DEST` / `TASK_INTERVAL` / `SERVER_PORT`
  - `PV_STATUS_CODES` / `PV_EXCLUDE_PATTERNS` / `PV_EXCLUDE_IPS`
  
### Nginx 日志格式
默认解析模式基于典型的 access log 格式：
```
<ip> - <user> [time] "METHOD /path HTTP/1.x" status bytes "referer" "ua"
```
如果你的 Nginx 使用自定义 `log_format`，需要同步调整 `internal/ingest/log_parser.go` 中的正则。

## 目录结构与主要文件

```
.
├── cmd/
│   └── nginxpulse/
│       └── main.go                 # 程序入口
├── internal/                       # 核心逻辑（解析、统计、存储、API）
│   ├── app/
│   │   └── app.go                  # 初始化、依赖装配、任务调度
│   ├── analytics/                  # 统计口径与聚合
│   ├── enrich/
│   │   ├── ip_geo.go               # IP 归属地（远程+本地）与缓存
│   │   └── pv_filter.go            # PV 过滤规则
│   ├── ingest/
│   │   └── log_parser.go           # 日志扫描、解析与入库
│   ├── server/
│   │   └── http.go                 # HTTP 服务与中间件
│   ├── store/
│   │   └── repository.go           # SQLite 结构与写入
│   └── web/
│       └── handler.go              # API 路由
├── webapp/
│   └── src/
│       └── main.ts                 # 前端入口
├── configs/
│   ├── nginxpulse_config.json      # 核心配置入口
│   └── nginx_frontend.conf         # 内置 Nginx 配置
├── scripts/
│   └── dev_local.sh                # 本地一键启动
├── var/                            # 数据目录（运行时生成/挂载）
├── Dockerfile
└── docker-compose.yml
```

---

如需更详细的统计口径或 API 扩展，建议从 `internal/analytics/` 与 `internal/web/handler.go` 开始。
