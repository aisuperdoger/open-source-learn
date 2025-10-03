# Go标准库net/http简单示例

这个项目展示了如何使用Go语言的标准库`net/http`创建一个简单的HTTP服务器。

## 功能介绍

- 根路径(`/`)返回简单的问候信息
- `/info`路径返回JSON格式的信息

## 如何运行

### 前提条件
- 已安装Go语言环境

### 运行步骤

1. 进入项目目录
```bash
cd /home/zwl/projects/projects/open-source-learn/net-http-demo
```

2. 编译并运行程序
```bash
go run main.go
```

3. 打开浏览器访问
   - http://localhost:8080 - 查看问候信息
   - http://localhost:8080/info - 查看JSON格式信息

4. 按`Ctrl+C`停止服务器

## 代码说明

- 使用`http.HandleFunc()`注册URL路径对应的处理函数
- 使用`http.ListenAndServe()`启动HTTP服务器
- 处理函数接收`http.ResponseWriter`和`*http.Request`参数
- 可以通过`http.Error()`返回错误信息
- 可以通过`w.Header().Set()`设置响应头
- 可以使用`fmt.Fprintf()`输出响应内容