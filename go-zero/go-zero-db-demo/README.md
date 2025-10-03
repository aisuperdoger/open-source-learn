# go-zero数据库操作示例

这是一个使用go-zero框架操作数据库的简单示例，演示了如何使用go-zero进行数据库的增删改查操作。

## 项目结构

```
.
├── go.mod
├── main.go
├── model/
│   ├── user.sql
│   ├── usermodel.go
│   └── error.go
└── README.md
```

## 功能特性

1. 使用go-zero的sqlx库进行数据库操作
2. 集成Redis缓存
3. 实现用户表的增删改查操作
4. 使用接口设计，便于测试和扩展

## 数据表结构

```sql
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `mobile` varchar(255) NOT NULL DEFAULT '',
  `gender` varchar(255) NOT NULL DEFAULT '',
  `nickname` varchar(255) NOT NULL DEFAULT '',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 使用方法

1. 确保已安装Go环境(1.16+)
2. 安装MySQL和Redis服务
3. 创建数据库和数据表
4. 修改main.go中的数据库和Redis连接配置
5. 运行示例程序：

```bash
# 克隆项目
git clone <repository-url>

# 进入项目目录
cd go-zero-db-demo

# 下载依赖
go mod tidy

# 编译
go build -o go-zero-db-demo .

# 运行
./go-zero-db-demo
```

## 核心代码说明

### 数据模型定义

在[model/usermodel.go](file:///home/zwl/projects/projects/open-source-learn/go-zero-db-demo/model/usermodel.go)中定义了User结构体和UserModel接口：

```go
type User struct {
    Id         int64     `db:"id"`
    Name       string    `db:"name"`
    Password   string    `db:"password"`
    Mobile     string    `db:"mobile"`
    Gender     string    `db:"gender"`
    Nickname   string    `db:"nickname"`
    CreateTime time.Time `db:"create_time"`
    UpdateTime time.Time `db:"update_time"`
}

type UserModel interface {
    Insert(data User) (int64, error)
    FindOne(id int64) (*User, error)
    FindOneByName(name string) (*User, error)
    FindOneByMobile(mobile string) (*User, error)
    Update(data User) error
    Delete(id int64) error
}
```

### 主要操作示例

1. **插入数据**：
```go
user := model.User{
    Name:     "user_123",
    Password: "password123",
    Mobile:   "13800000000",
    Gender:   "male",
    Nickname: "test_user",
}
id, err := userModel.Insert(user)
```

2. **查询数据**：
```go
// 根据ID查询
user, err := userModel.FindOne(id)

// 根据用户名查询
user, err := userModel.FindOneByName("user_123")
```

3. **更新数据**：
```go
user.Nickname = "updated_user"
err := userModel.Update(user)
```

4. **删除数据**：
```go
err := userModel.Delete(id)
```

## 依赖

- github.com/zeromicro/go-zero v1.6.0
- MySQL数据库
- Redis缓存（可选）

## 注意事项

1. 请根据实际环境修改数据库和Redis连接配置
2. 确保MySQL和Redis服务正常运行
3. 首次运行前请创建相应的数据库和数据表