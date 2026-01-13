# TestFlight Monitor

<p align="center">
  <img src="https://developer.apple.com/assets/elements/icons/testflight/testflight-96x96_2x.png" width="96" height="96" alt="TestFlight">
</p>

<p align="center">
  <strong>自动监控 TestFlight 名额，有位立即通知</strong>
</p>

<p align="center">
  <a href="#功能特性">功能</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#docker-部署">Docker 部署</a> •
  <a href="#配置说明">配置</a> •
  <a href="README_EN.md">English</a>
</p>

---

## 功能特性

- 🔍 **自动检测** - 定时检测 TestFlight 是否有空位
- 📱 **Telegram 通知** - 有位时立即推送通知
- ⏱️ **自定义间隔** - 检测间隔可调（最低 10 秒）
- ♾️ **永久监控** - 支持永久监控，不过期
- 🌐 **代理支持** - 支持 HTTP/SOCKS5 代理
- 🌍 **中英双语** - 界面支持简体中文和 English
- 🐳 **Docker 部署** - 一键部署，开箱即用

## 截图预览

监控界面简洁直观：
- 左侧添加监控、配置 Telegram
- 右侧卡片展示监控状态
- 支持暂停/恢复/编辑/删除

## 快速开始

### 方式一：Docker Compose（推荐）

```bash
# 克隆项目
git clone https://github.com/your-username/tf-monitor.git
cd tf-monitor

# 启动服务
docker-compose up -d

# 访问 http://localhost:8080
```

### 方式二：手动编译

```bash
# 编译前端
cd web && npm install && npm run build && cd ..

# 编译后端
go build -o tf-monitor ./cmd/server/main.go

# 运行
./tf-monitor
```

## Docker 部署

### docker-compose.yml

```yaml
version: "3.8"

services:
  tf-monitor:
    build: .
    container_name: tf-monitor
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - TZ=Asia/Shanghai
      - PROXY_ENABLED=false
      - PROXY_URL=
```

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `SERVER_PORT` | 8080 | 服务端口 |
| `DB_PATH` | data/tf-monitor.db | 数据库路径 |
| `PROXY_ENABLED` | false | 是否启用代理 |
| `PROXY_URL` | - | 代理地址，如 `http://127.0.0.1:7890` |

## 配置说明

### Telegram 通知配置

1. 向 [@BotFather](https://t.me/BotFather) 发送 `/newbot` 创建机器人
2. 获取 Bot Token（格式：`123456789:ABCdefGHI...`）
3. 向 [@userinfobot](https://t.me/userinfobot) 发送消息获取 Chat ID
4. 在设置中填入 Bot Token 和 Chat ID
5. 点击「测试发送」验证配置

### 代理配置

国内访问 TestFlight 可能需要代理：

```bash
# HTTP 代理
PROXY_URL=http://127.0.0.1:7890

# SOCKS5 代理
PROXY_URL=socks5://127.0.0.1:7890
```

## 使用说明

### 添加监控

1. 在左侧输入 TestFlight 链接（支持批量，每行一个）
2. 设置检测间隔（建议 60 秒）
3. 选择监控时长（2h/8h/12h/24h/永久）
4. 点击「添加监控」

### 通知模式

| 模式 | 说明 |
|------|------|
| 循环推送 | 每次检测到有位都通知 |
| 仅一次 | 通知一次后停止 |
| 状态变化 | 仅当状态从「已满」变为「有位」时通知 |

### 卡片操作

- **暂停/恢复** - 暂停或恢复监控
- **编辑** - 修改检测间隔和监控时长
- **删除** - 删除监控

## API 文档

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/monitors | 获取监控列表 |
| POST | /api/monitors | 添加监控 |
| PUT | /api/monitors/:id | 更新监控 |
| DELETE | /api/monitors/:id | 删除监控 |
| POST | /api/monitors/:id/toggle | 暂停/恢复监控 |
| GET | /api/telegram | 获取 Telegram 配置 |
| PUT | /api/telegram | 更新 Telegram 配置 |
| POST | /api/telegram/test | 测试 Telegram 通知 |
| GET | /api/status | 获取服务状态 |

## 技术栈

- **后端**: Go + Gin + GORM + SQLite
- **前端**: Vue 3 + TypeScript + Vite
- **部署**: Docker

## 许可证

MIT License

## 致谢

- [Apple TestFlight](https://developer.apple.com/testflight/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Vue.js](https://vuejs.org/)
