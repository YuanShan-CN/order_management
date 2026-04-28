# 📷 摄影订单管理系统

一个基于 Go + Gin + MariaDB 构建的摄影订单管理系统，提供订单管理、店铺维护、收入统计等功能，支持多用户登录和数据隔离。

## ✨ 功能特性

- **订单管理** - 创建、编辑、删除订单，支持订单状态管理
- **店铺维护** - 店铺信息管理，支持回收站功能
- **收入统计** - 多维度数据统计，支持表格、条形图、饼图展示
- **回收站** - 支持订单和店铺的软删除与恢复
- **数据导入** - 支持 CSV 文件批量导入订单数据
- **数据导出** - 支持将订单数据导出为 CSV 文件
- **多用户支持** - 用户登录、数据隔离、密码管理
- **JWT 认证** - 安全的 token 认证机制

## 🛠️ 技术栈

- **后端**: Go 1.26 + Gin Framework
- **数据库**: MariaDB 11.5
- **ORM**: GORM v2
- **认证**: JWT (JSON Web Token)
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

# 5. 初始化用户
docker exec -it order_management_app_1 ./init_user -username admin -password your_password

# 6. 访问系统
# 登录页面: http://localhost:8080/login
# 数据库端口: 3306
```

### 本地开发

```bash
# 1. 设置环境变量
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=order_management

# 2. 安装依赖
go mod tidy

# 3. 编译项目
go build -o order_management .

# 4. 初始化用户
./order_management init -username admin -password your_password

# 5. 运行服务
./order_management -c config.yaml
```

## 🔐 用户管理

### 创建初始用户

```bash
# 使用初始化脚本创建用户
go run scripts/init_user.go -username admin -password 123456

# 或使用编译后的二进制
./init_user -username admin -password 123456
```

### 重置用户密码

```bash
go run scripts/reset_password.go -username admin -password new_password
```

## 📁 项目结构

```
order_management/
├── src/             # 源代码目录
│   ├── config/      # 配置管理
│   │   └── config.go
│   ├── controllers/ # 控制器
│   │   ├── auth_controller.go   # 认证控制器
│   │   ├── order_controller.go
│   │   └── shop_controller.go
│   ├── database/    # 数据库连接
│   │   └── database.go
│   ├── middleware/  # 中间件
│   │   ├── auth.go  # JWT 认证中间件
│   │   └── cors.go
│   ├── models/      # 数据模型
│   │   ├── order.go
│   │   ├── shop.go
│   │   └── user.go  # 用户模型
│   └── router/      # 路由配置
│       └── router.go
├── scripts/         # 辅助脚本
│   ├── import_csv.go      # CSV 数据导入
│   ├── fill_shops.go      # 店铺数据填充
│   ├── init_user.go       # 用户初始化
│   └── reset_password.go  # 密码重置
├── migrations/      # 数据库迁移
│   └── 20260428_add_user_tables.sql
├── templates/       # HTML 模板
│   ├── login.html   # 登录页面
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
├── config.yaml      # 应用配置
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

jwt:
  secret: your-secret-key-change-in-production
  expire_hour: 24
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

## 📤 数据导出

### 导出订单为 CSV

通过前端页面或 API 导出订单数据：

**前端导出：**
- 登录后访问订单管理页面 (`/orders`) 或收入统计页面 (`/stats`)
- 设置筛选条件（可选）
- 点击橙色的 **📥 导出CSV** 按钮

**API 导出（需要认证）：**

```bash
# 获取登录 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your_password"}' | jq -r '.token')

# 导出所有订单
curl -H "Authorization: Bearer $TOKEN" \
  -o orders.csv http://localhost:8080/api/orders/export

# 带筛选条件导出
curl -H "Authorization: Bearer $TOKEN" \
  -o orders.csv "http://localhost:8080/api/orders/export?year=2026&settled=settled"
```

**导出参数：**

| 参数 | 说明 | 示例 |
|-----|------|------|
| search | 搜索关键词 | search=故宫 |
| settled | 结算状态 | settled 或 unsettled |
| year | 年份筛选 | year=2026 |
| month | 月份筛选 | month=4 |
| shop | 店铺筛选 | shop=某店铺 |

## 🗂️ API 接口

### 认证接口

| 方法 | 路径 | 描述 | 需要认证 |
|-----|------|------|---------|
| POST | /api/auth/login | 用户登录 | 否 |
| POST | /api/auth/logout | 用户登出 | 否 |
| GET | /api/auth/current | 获取当前用户 | 是 |
| PUT | /api/auth/password | 修改密码 | 是 |
| DELETE | /api/auth/account | 删除账号 | 是 |

### 订单接口

| 方法 | 路径 | 描述 | 需要认证 |
|-----|------|------|---------|
| GET | /api/orders | 获取订单列表 | 是 |
| GET | /api/orders/:id | 获取单个订单 | 是 |
| POST | /api/orders | 创建订单 | 是 |
| PUT | /api/orders/:id | 更新订单 | 是 |
| DELETE | /api/orders/:id | 软删除订单 | 是 |
| GET | /api/orders/trash | 获取回收站订单 | 是 |
| POST | /api/orders/:id/restore | 恢复订单 | 是 |
| DELETE | /api/orders/:id/force-delete | 永久删除订单 | 是 |
| GET | /api/orders/export | 导出订单为 CSV | 是 |
| GET | /api/orders/stats/export | 导出统计数据为 CSV | 是 |

### 店铺接口

| 方法 | 路径 | 描述 | 需要认证 |
|-----|------|------|---------|
| GET | /api/shops | 获取店铺列表 | 是 |
| GET | /api/shops/:id | 获取单个店铺 | 是 |
| POST | /api/shops | 创建店铺 | 是 |
| PUT | /api/shops/:id | 更新店铺 | 是 |
| DELETE | /api/shops/:id | 软删除店铺 | 是 |
| GET | /api/shops/trash | 获取回收站店铺 | 是 |
| POST | /api/shops/:id/restore | 恢复店铺 | 是 |
| DELETE | /api/shops/:id/force-delete | 永久删除店铺 | 是 |

## 📄 前端页面

| 路径 | 描述 | 需要登录 |
|-----|------|---------|
| /login | 登录页面 | 否 |
| / | 首页 | 否 |
| /orders | 订单管理 | 是 |
| /shops | 店铺维护 | 是 |
| /stats | 收入统计 | 是 |
| /trash | 回收站 | 是 |

## 🔑 认证说明

系统使用 JWT (JSON Web Token) 进行身份认证：

1. 用户登录成功后获取 token
2. token 有效期为 24 小时（可配置）
3. 访问受保护接口时需在请求头中携带 `Authorization: Bearer <token>`
4. 前端自动处理 token 的存储和携带

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