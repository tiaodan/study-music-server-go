# study-api-autotest-go

基于 Go 语言开发的音乐网站后端 API 服务，是 Java 版本 [study-music-server-java](https://github.com/example/study-music-server-java) 的 Go 语言移植版本。

## 技术栈

- **Web 框架**: Gin v1.9.1
- **ORM**: Gorm v1.25.5
- **数据库**: MySQL
- **缓存**: Redis

## 项目结构

```
study-api-autotest-go/
├── main.go                    # 程序入口
├── go.mod                     # Go 模块依赖
├── config.yaml                # 配置文件
│
├── common/                    # 公共模块
│   ├── response.go            # 统一响应结构
│   └── constants.go           # 常量定义
│
├── config/                    # 配置模块
│   └── config.go              # 配置加载
│
├── models/                    # 数据模型
│   ├── consumer.go            # 用户模型
│   ├── song.go                # 歌曲模型
│   ├── singer.go              # 歌手模型
│   ├── songlist.go            # 歌单模型
│   ├── collect.go             # 收藏模型
│   ├── comment.go             # 评论模型
│   ├── ranklist.go            # 排行榜模型
│   ├── banner.go              # 轮播图模型
│   ├── admin.go               # 管理员模型
│   ├── listsong.go            # 歌单歌曲模型
│   ├── usersupport.go         # 用户点赞模型
│   └── request.go             # 请求参数模型
│
├── mapper/                    # 数据访问层
│   ├── db.go                  # 数据库连接
│   └── *_mapper.go            # 各模型的数据操作
│
├── service/                   # 业务逻辑层
│   └── *_service.go           # 各模块业务逻辑
│
├── controller/                # 控制器层
│   └── *_controller.go        # API 接口处理
│
├── routes/                    # 路由配置
│   └── routes.go              # 路由定义
│
├── middleware/                # 中间件
│   └── middleware.go          # CORS、日志等中间件
│
└── utils/                     # 工具类
    ├── file.go                # 文件处理
    └── random.go              # 随机数生成
```

## 功能模块

| 模块 | 功能说明 |
|------|---------|
| 用户管理 | 用户注册、登录、密码重置、头像更新 |
| 歌手管理 | 歌手的增删改查 |
| 歌曲管理 | 歌曲的增删改查、按歌手查询 |
| 歌单管理 | 歌单的增删改查 |
| 收藏功能 | 用户收藏歌曲、取消收藏 |
| 评论功能 | 歌曲评论、歌单评论 |
| 排行榜 | 歌单评分功能 |
| Banner | 轮播图管理 |
| 管理员 | 管理员登录 |

## 快速开始

### 1. 安装 Go 依赖

```bash
cd study-api-autotest-go
go mod tidy
```

### 2. 配置数据库

编辑 `config.yaml` 文件，配置 MySQL、Redis 等连接信息：

```yaml
server:
  port: 8003

database:
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: tp_music

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0

minio:
  endpoint: http://localhost:9000
  access-key: root
  secret-key: 123456789
  bucket-name: user01

mail:
  host: smtp.163.com
  port: 465
  username: ""
  password: ""
  address: ""
```

### 3. 创建数据库

执行 SQL 脚本创建数据库（参考原 Java 项目的 sql/tp_music.sql）

### 4. 启动服务

```bash
go run main.go
```

服务启动后默认运行在 `http://localhost:8003`

## API 接口

### 用户接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /user/add | 用户注册 |
| POST | /user/login/status | 用户登录 |
| POST | /user/email/status | 邮箱登录 |
| GET | /user | 获取所有用户 |
| GET | /user/detail | 获取指定用户 |
| GET | /user/delete | 删除用户 |
| POST | /user/update | 更新用户信息 |
| POST | /user/updatePassword | 更新密码 |
| POST | /user/avatar/update | 更新头像 |

### 歌手接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /singer/add | 添加歌手 |
| POST | /singer/update | 更新歌手 |
| GET | /singer/delete | 删除歌手 |
| GET | /singer/detail | 获取歌手详情 |
| GET | /singer/name/detail | 按名字查询歌手 |
| GET | /singer | 获取所有歌手 |

### 歌曲接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /song/add | 添加歌曲 |
| POST | /song/update | 更新歌曲 |
| GET | /song/delete | 删除歌曲 |
| GET | /song/detail | 获取歌曲详情 |
| GET | /song/singer/detail | 按歌手查询歌曲 |
| GET | /song/name/detail | 按名字查询歌曲 |
| GET | /song | 获取所有歌曲 |

### 歌单接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /songList/add | 添加歌单 |
| POST | /songList/update | 更新歌单 |
| GET | /songList/delete | 删除歌单 |
| GET | /songList/detail | 获取歌单详情 |
| GET | /songList/name/detail | 按标题查询歌单 |
| GET | /songList | 获取所有歌单 |

### 收藏接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /collect/add | 添加收藏 |
| GET | /collect/delete | 取消收藏 |
| GET | /collect/detail | 获取用户收藏 |

### 评论接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /comment/add | 添加评论 |
| GET | /comment/delete | 删除评论 |
| GET | /comment/song/detail | 获取歌曲评论 |
| GET | /comment/songList/detail | 获取歌单评论 |

### 排行榜接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /rankList/add | 添加评分 |
| GET | /rankList/detail | 获取歌单评分 |

### 其他接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /banner | 获取轮播图 |
| POST | /admin/login | 管理员登录 |
| POST | /listSong/add | 添加歌曲到歌单 |
| GET | /listSong/detail | 获取歌单歌曲列表 |
| POST | /userSupport/add | 点赞评论 |
| GET | /userSupport/delete | 取消点赞 |

## 注意事项

1. **数据库迁移**: 项目启动时会自动创建表（使用 Gorm AutoMigrate），但建议先执行原始 SQL 脚本确保数据完整性

2. **配置文件**: 首次运行前请务必修改 `config.yaml` 中的数据库连接信息

3. **端口占用**: 默认使用 8003 端口，如需修改请编辑 `config.yaml`

4. **静态文件**: 图片上传功能需要确保 `img` 目录存在

5. **密码加密**: 密码使用 MD5 + 盐值加密，与原 Java 项目保持一致

6. **Redis**: 验证码功能依赖 Redis，请确保 Redis 服务正常运行

## 与原 Java 项目对比

| 特性 | Java (Spring Boot) | Go (Gin) |
|------|-------------------|----------|
| 框架 | Spring Boot 2.6.2 | Gin 1.9.1 |
| ORM | MyBatis Plus | Gorm |
| 数据库 | MySQL | MySQL |
| 缓存 | Redis | Redis |
| 对象存储 | MinIO | MinIO |

## 许可证

MIT License
