# 📷 摄影订单管理系统

一个基于 Go + Gin + MariaDB 构建的摄影订单管理系统，提供订单管理、店铺维护、收入统计等功能，支持多用户登录和数据隔离。

## ✨ 功能特性

- **订单管理** - 创建、编辑、删除订单，支持订单状态管理
- **店铺维护** - 店铺信息管理，支持回收站功能
- **收入统计** - 多维度数据统计，支持表格、条形图、饼图展示，后端聚合查询优化性能
- **回收站** - 支持订单和店铺的软删除与恢复
- **数据导入** - 支持 CSV 文件批量导入订单数据，自动创建缺失的店铺，友好的错误提示
- **数据导出** - 支持将订单数据导出为 CSV 文件，日期只显示年月日（YYYY-MM-DD）
- **定时任务** - 每日自动导出 CSV，每小时检查，按用户时区独立处理
- **多用户支持** - 用户登录、数据隔离、密码管理、用户时区配置
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

# 2. 复制环境变量和配置模板
cp .env.example .env
cp config.yaml.example config.yaml

# 3. 修改配置文件
# 编辑 .env 文件，设置 DB_ROOT_PASSWORD 和 DB_PASSWORD
# 编辑 config.yaml 文件，根据需要修改数据库连接信息和 JWT 密钥

# 4. 启动服务
docker-compose build --no-cache
docker-compose up -d

# 5. 等待服务启动 (约5秒)
sleep 5

# 6. 创建初始用户
docker exec -it order_management_app ./init_user -username admin -password your_password

# 7. 导入订单数据 (可选)
# 登录系统后，在订单管理页面点击"📤 导入CSV"按钮
# 支持 orders.csv 和 stats.csv 两种格式

# 8. 访问系统
# 登录页面: http://localhost:8081/login
# 数据库端口: 3306
```

### 本地开发

```bash
# 1. 复制配置文件模板
cp config.yaml.example config.yaml

# 2. 修改配置文件
# 编辑 config.yaml，设置数据库连接信息和 JWT 密钥

# 3. 设置环境变量 (可选，如果使用 Docker 数据库)
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=order_management

# 4. 安装依赖
go mod tidy

# 5. 编译项目
go build -o order_management .

# 6. 初始化用户（使用单独的脚本）
go run scripts/init_user.go -username admin -password your_password

# 7. 运行服务
./order_management -c config.yaml
```

### Docker 本地重新部署

```bash
# 停止并删除容器
# docker-compose down

# 重新构建镜像（不使用缓存）
# docker-compose build --no-cache

# 启动服务（后台运行）
docker-compose up -d --build

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

### 环境变量 (.env)

| 变量名          | 说明     | 默认值               |
| ------------ | ------ | ----------------- |
| DB\_HOST     | 数据库主机  | db                |
| DB\_PORT     | 数据库端口  | 3306              |
| DB\_USER     | 数据库用户名 | root              |
| DB\_PASSWORD | 数据库密码  | -                 |
| DB\_NAME     | 数据库名称  | order\_management |

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

**配置说明：**

- `database`: 数据库连接信息
- `server`: HTTP 服务器配置
- `jwt`: JWT 认证配置，`secret` 请在生产环境修改为强密钥

**注意：**

- `config.yaml` 和 `.env` 都已在 `.gitignore` 中，不会被提交
- 首次部署时，请从模板文件复制：
  ```bash
  cp config.yaml.example config.yaml
  cp .env.example .env
  ```
- Docker 部署时，config.yaml 会被挂载到容器内

## 📊 数据导入

### 导入 CSV 订单数据

**通过前端页面导入（推荐）：**

- 登录后访问订单管理页面 (`/orders`)
- 点击紫色的 **📤 导入CSV** 按钮
- 选择 CSV 文件并提交
- 系统会自动创建缺失的店铺

**CSV 格式要求：**

- 必须包含列：日期、店铺、定金、尾款、已结算
- 可选列：地点、内容
- 已结算列填写"是"/"否" 或 "true"/"false" 或 "1"/"0"
- 日期格式：YYYY-MM-DD 或 YYYY-MM-DDTHH:MM:SS±HH:MM
- 支持导出的 orders.csv 和 stats.csv 格式

**命令行导入（备用）：**

```bash
# 设置环境变量
export DB_PASSWORD=your_password

# 运行导入脚本
cd scripts
go run import_csv.go /path/to/orders.csv
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

| 参数      | 说明    | 示例                  |
| ------- | ----- | ------------------- |
| search  | 搜索关键词 | search=旅拍           |
| settled | 结算状态  | settled 或 unsettled |
| year    | 年份筛选  | year=2026           |
| month   | 月份筛选  | month=4             |
| shop    | 店铺筛选  | shop=某店铺            |

### 导入 CSV 订单数据

```bash
# 获取登录 token
TOKEN=$(curl -s -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your_password"}' | jq -r '.token')

# 导入 CSV 文件
curl -H "Authorization: Bearer $TOKEN" \
  -X POST \
  -F "file=@/path/to/orders.csv" \
  http://localhost:8081/api/orders/import
```

**CSV 文件格式要求：**

**支持的 CSV 格式：**

- `orders.csv` - 标准订单格式（日期、店铺、地点、内容、定金、尾款、已结算）
- `stats.csv` - 统计导出格式（日期、店铺、地点、内容、定金、尾款、已结算、收入）

**导入响应示例：**

```json
{
  "importedCount": 220,
  "newShopCount": 15,
  "failedCount": 0,
  "errors": [],
  "message": "成功导入 220 条数据，新建 15 个店铺"
}
```

**错误响应示例：**

```json
{
  "error": "CSV格式不正确，缺少必要的列：日期、店铺\n期望的CSV格式示例：\n日期,店铺,地点,内容,定金,尾款,已结算\n2026-04-29,示例店,北京,购买商品,500.00,500.00,否"
}
```

## 🗂️ API 接口

### 认证接口

| 方法     | 路径                 | 描述     | 需要认证 |
| ------ | ------------------ | ------ | ---- |
| POST   | /api/auth/login    | 用户登录   | 否    |
| POST   | /api/auth/logout   | 用户登出   | 否    |
| GET    | /api/auth/current  | 获取当前用户 | 是    |
| PUT    | /api/auth/password | 修改密码   | 是    |
| DELETE | /api/auth/account  | 删除账号   | 是    |

### 订单接口

| 方法     | 路径                             | 描述                 | 需要认证 |
| ------ | ------------------------------ | ------------------ | ---- |
| GET    | /api/orders                    | 获取订单列表             | 是    |
| GET    | /api/orders/:id                | 获取单个订单             | 是    |
| POST   | /api/orders                    | 创建订单               | 是    |
| PUT    | /api/orders/:id                | 更新订单               | 是    |
| DELETE | /api/orders/:id                | 软删除订单              | 是    |
| GET    | /api/orders/trash              | 获取回收站订单            | 是    |
| POST   | /api/orders/:id/restore        | 恢复订单               | 是    |
| DELETE | /api/orders/:id/force-delete   | 永久删除订单             | 是    |
| GET    | /api/orders/export             | 导出订单为 CSV          | 是    |
| GET    | /api/orders/stats/export       | 导出统计数据为 CSV        | 是    |
| GET    | /api/orders/stats              | 获取统一统计数据（摘要+月度+店铺） | 是    |
| GET    | /api/orders/stats/summary      | 获取统计摘要             | 是    |
| GET    | /api/orders/stats/monthly      | 获取月度收入趋势           | 是    |
| GET    | /api/orders/stats/shops        | 获取店铺收入统计           | 是    |
| GET    | /api/orders/stats/shop-details | 获取店铺订单详情（分页）       | 是    |

### 店铺接口

| 方法     | 路径                          | 描述      | 需要认证 |
| ------ | --------------------------- | ------- | ---- |
| GET    | /api/shops                  | 获取店铺列表  | 是    |
| GET    | /api/shops/:id              | 获取单个店铺  | 是    |
| POST   | /api/shops                  | 创建店铺    | 是    |
| PUT    | /api/shops/:id              | 更新店铺    | 是    |
| DELETE | /api/shops/:id              | 软删除店铺   | 是    |
| GET    | /api/shops/trash            | 获取回收站店铺 | 是    |
| POST   | /api/shops/:id/restore      | 恢复店铺    | 是    |
| DELETE | /api/shops/:id/force-delete | 永久删除店铺  | 是    |

## 📄 前端页面

| 路径      | 描述   | 需要登录 |
| ------- | ---- | ---- |
| /login  | 登录页面 | 否    |
| /       | 首页   | 否    |
| /orders | 订单管理 | 是    |
| /shops  | 店铺维护 | 是    |
| /stats  | 收入统计 | 是    |
| /trash  | 回收站  | 是    |

## 🔑 认证说明

系统使用 JWT (JSON Web Token) 进行身份认证：

1. 用户登录成功后获取 token
2. token 有效期为 24 小时（可配置）
3. 访问受保护接口时需在请求头中携带 `Authorization: Bearer <token>`
4. 前端自动处理 token 的存储和携带

## 📋 数据格式说明

### 日期格式

系统使用 `YYYY-MM-DD` 格式处理所有日期：

- 导入/导出 CSV 时只显示年月日
- 前端日期选择器只选择年月日
- 数据库使用 `DATE` 类型存储

### 时区安全设计

系统采用以下最佳实践处理时区：

1. **核心策略**：**后端全用 UTC，前端自己格式化显示**
   - 数据库连接使用 `loc=UTC`
   - JWT 时间戳使用 UTC
   - 所有时间处理使用 UTC
2. **用户时区**：
   - `User.timezone` 字段存储每个用户的时区（如 `Asia/Shanghai`）
   - 定时任务按用户时区判断日期
   - 导出文件名使用用户时区的日期
3. **订单日期**：
   - 使用 `DATE` 类型（仅存储年月日，无时间信息）
   - Go 层使用 `string` 类型而非 `time.Time`
   - 不涉及时区转换
4. **API 通信**：建议使用时间戳或 ISO 8601 UTC 字符串，前端按用户本地时区展示

### 定时自动导出

系统支持每日自动导出 CSV，具有以下特性：

- **调度策略**：每小时检查一次（按 UTC 时间）
- **用户独立时区**：每个用户可设置自己的时区，按用户时区判断"今天"
- **文件名日期**：使用用户时区的日期（如 `admin_2026-05-11.csv`）
- **容错设计**：Mac 睡眠后唤醒立即补导，不遗漏任务
- **调试接口**：`GET /api/scheduler/run-now` 可手动触发导出
- **容器环境**：仅在 Docker 容器中启用定时任务，本地开发禁用
- **存储位置**：Docker 中 `/app/data/csv`，本地 `./data/csv`

## 📊 数据库迁移

项目使用 GORM AutoMigrate 自动处理数据库结构变更，但对于复杂的迁移，手动执行 SQL 文件：

### 迁移文件列表

迁移文件位于 `migrations/` 目录：

| 迁移文件                          | 说明              |
| ------------------------------ | --------------- |
| 20260426_remove_edited_column.sql | 移除 edited 字段   |
| 20260428_add_user_tables.sql       | 添加用户表          |
| 20260429_add_deposit_balance.sql | 添加定金尾款字段       |
| 20260511_add_user_timezone.sql   | 添加用户时区字段        |

### 执行迁移

**方式 1：通过 docker-compose（推荐）**

```bash
# 启动服务时 GORM AutoMigrate 会自动处理简单的表结构变更
docker-compose up -d

# 对于复杂迁移（如添加字段），如果有需要可手动执行 SQL
docker exec -i my_mariadb mysql -u root -p order_management < migrations/20260511_add_user_timezone.sql
```

**方式 2：本地开发**

```bash
# 使用 mysql 命令行
mysql -u root -p order_management < migrations/20260511_add_user_timezone.sql
```

### 迁移说明

- 大部分情况下，GORM AutoMigrate 会自动处理（如添加新列、修改字段等）
- 需要手动迁移的情况：删除/重命名列、修改数据等破坏性操作
- 每个迁移文件有时间戳前缀，按顺序执行

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
