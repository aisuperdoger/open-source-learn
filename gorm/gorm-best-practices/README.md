# GORM最佳实践示例

本项目演示了在Go语言中使用GORM操作MySQL数据库的最佳实践。

## 项目结构

```
gorm-best-practices/
├── config/          # 数据库配置
├── models/          # 数据库结构和go语言结构映射关系
├── repository/      # 数据库表的各种操作都封装成类的函数
├── service/         # 调用repository实现对数据库的操作，从而实现业务逻辑
```

## GORM最佳实践要点

### 1. 连接池配置
```go
// 设置连接池参数
sqlDB, err := db.DB()
sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
```

### 2. 模型定义规范
- 使用结构体标签定义字段约束
- 合理使用索引提高查询性能
- 定义关联关系便于数据操作

### 3. 错误处理
- 检查并处理GORM返回的错误
- 区分不同类型的错误（如记录未找到）

### 4. 事务处理
- 使用事务保证数据一致性
- 正确处理事务的提交和回滚

### 5. 预加载关联数据
```go
// 预加载用户的文章
db.Preload("Posts").First(&user, userID)
```

### 6. 分页查询
```go
// 分页查询用户列表
offset := (page - 1) * pageSize
db.Offset(offset).Limit(pageSize).Find(&users)
```

## 运行项目

1. 确保MySQL服务正在运行
2. 创建数据库：
   ```sql
   CREATE DATABASE testdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```
3. 安装依赖：
   ```bash
   go mod tidy
   ```
4. 运行程序：
   ```bash
   go run main.go
   ```

## 数据库连接信息

项目使用以下MySQL连接信息：
- Host: 127.0.0.1
- Port: 3306
- User: root
- Password: zwl1819123
- Database: testdb

## 主要特性演示

1. **自动迁移**：GORM自动创建和更新数据库表结构
2. **CRUD操作**：创建、读取、更新、删除用户数据
3. **关联关系**：用户与订单、文章的关联关系
4. **事务处理**：保证数据一致性的事务操作
5. **预加载**：高效加载关联数据
6. **分页查询**：处理大量数据的分页查询
7. **原生SQL**：必要时使用原生SQL查询

## 注意事项

1. **生产环境密码安全**：不要在代码中硬编码密码，应使用环境变量或配置文件
2. **连接池调优**：根据实际负载调整连接池参数
3. **索引优化**：为常用查询字段添加合适的索引
4. **错误日志**：在生产环境中记录详细的错误日志