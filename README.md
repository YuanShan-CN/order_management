# 📷 摄影订单管理系统

一个基于 Go + Gin + MariaDB 构建的摄影订单管理系统，提供订单管理、店铺维护、收入统计等功能。

## ✨ 功能特性

- **订单管理** - 创建、编辑、删除订单，支持订单状态管理
- **店铺维护** - 店铺信息管理，支持回收站功能
- **收入统计** - 多维度数据统计，支持表格、条形图、饼图展示
- **回收站** - 支持订单和店铺的软删除与恢复
- **数据导入** - 支持 CSV 文件批量导入订单数据

## 🛠️ 技术栈

- **后端**: Go 1.22 + Gin Framework
- **数据库**: MariaDB 11.5
- **ORM**: GORM v2
- **前端**: HTML5 + CSS3 + JavaScript
- **容器化**: Docker + Docker Compose

## 🚀 快速开始

### 环境要求

- Go 1.22+
- Docker & Docker Compose
- MariaDB 11.5+ (如果不使用 Docker)

### 使用 Docker 部署

```bash
# 1. 克隆项目
git clone https://github.com/YuanShan-CN/order_management.git
cd order_management

# 2. 复制环境变量模板
cp .env.example .env

# 3. 修改 .env 文件中的密码
# 编辑 .env 文件，设置 DB_ROOT_PASSWORD 和 DB_PASSWORD

# 4. 启动服务
docker-compose up -d

# 5. 访问系统
# 前端页面: http://localhost:8081
# 数据库端口: 3307
```

### 本地开发

```bash
# 1. 设置环境变量
export DB_HOST=127.0.0.1
export DB_PORT=3307
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=order_management

# 2. 安装依赖
go mod tidy

# 3. 运行服务
go run main.go

# 或指定配置文件
go run main.go -c config.yaml
```

## 📁 项目结构

```
order_management/
├── config/          # 配置管理
│   └── config.go    # 配置解析
├── controllers/     # 控制器
│   ├── order_controller.go
│   └── shop_controller.go
├── database/        # 数据库连接
│   └── database.go
├── models/          # 数据模型
│   ├── order.go
│   └── shop.go
├── scripts/         # 辅助脚本
│   ├── import_csv.go      # CSV 数据导入
│   └── fill_shops.go      # 店铺数据填充
├── templates/       # HTML 模板
│   ├── admin.html
│   ├── index.html
│   ├── orders.html
│   ├── shops.html
│   ├── stats.html
│   └── trash.html
├── .env             # 环境变量 (不提交)
├── .env.example     # 环境变量模板
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go          # 入口文件
└── README.md
```

## 🔧 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|-------|------|-------|
| DB_HOST | 数据库主机 | db |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | - |
| DB_NAME | 数据库名称 | order_management |

### 配置文件 (config.yaml)

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: order_management

server:
  port: 8080
```

## 📊 数据导入

### 导入 CSV 订单数据

```bash
# 设置环境变量
export DB_PASSWORD=your_password

# 运行导入脚本
cd scripts
go run import_csv.go /path/to/orders.csv
```

### 填充店铺数据

```bash
# 从订单表提取商家名称填充店铺表
cd scripts
go run fill_shops.go
```

## 🗂️ API 接口

### 订单接口

| 方法 | 路径 | 描述 |
|-----|------|------|
| GET | /api/orders | 获取订单列表 |
| GET | /api/orders/:id | 获取单个订单 |
| POST | /api/orders | 创建订单 |
| PUT | /api/orders/:id | 更新订单 |
| DELETE | /api/orders/:id | 软删除订单 |
| GET | /api/orders/trash | 获取回收站订单 |
| POST | /api/orders/:id/restore | 恢复订单 |
| DELETE | /api/orders/:id/force-delete | 永久删除订单 |

### 店铺接口

| 方法 | 路径 | 描述 |
|-----|------|------|
| GET | /api/shops | 获取店铺列表 |
| GET | /api/shops/:id | 获取单个店铺 |
| POST | /api/shops | 创建店铺 |
| PUT | /api/shops/:id | 更新店铺 |
| DELETE | /api/shops/:id | 软删除店铺 |
| GET | /api/shops/trash | 获取回收站店铺 |
| POST | /api/shops/:id/restore | 恢复店铺 |
| DELETE | /api/shops/:id/force-delete | 永久删除店铺 |

## 📄 前端页面

| 路径 | 描述 |
|-----|------|
| / | 首页 |
| /orders | 订单管理 |
| /shops | 店铺维护 |
| /stats | 收入统计 |
| /trash | 回收站 |

## 📝 开发说明

### 代码规范

- 使用 `go fmt` 格式化代码
- 使用 `go vet` 检查代码
- 使用 `golint` 进行代码静态分析

### 提交规范

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具

## 📄 License

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！