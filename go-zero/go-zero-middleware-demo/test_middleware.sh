#!/bin/bash

# go-zero中间件功能测试脚本

echo "=== go-zero 中间件功能测试 ==="
echo ""

BASE_URL="http://localhost:8888"

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_test() {
    echo -e "${YELLOW}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# 检查服务是否启动
check_service() {
    echo "检查服务状态..."
    if curl -s -f "$BASE_URL/user/info/123" > /dev/null 2>&1; then
        print_success "服务已启动"
    else
        print_error "服务未启动，请先运行: go run user.go"
        exit 1
    fi
    echo ""
}

# 1. 测试用户登录（无中间件，仅全局中间件）
test_login() {
    print_test "测试用户登录 (仅全局中间件)"
    
    echo "1.1 正确的用户名密码:"
    curl -X POST "$BASE_URL/user/login" \
        -H "Content-Type: application/json" \
        -H "User-Agent: Test-Client/1.0" \
        -d '{"username":"user123","password":"password123"}' \
        -w "\nHTTP状态码: %{http_code}\n\n"
    
    echo "1.2 错误的密码:"
    curl -X POST "$BASE_URL/user/login" \
        -H "Content-Type: application/json" \
        -H "User-Agent: Test-Client/1.0" \
        -d '{"username":"user123","password":"wrong"}' \
        -w "\nHTTP状态码: %{http_code}\n\n"
}

# 2. 测试获取用户信息（UserAgent中间件）
test_userinfo() {
    print_test "测试获取用户信息 (UserAgent中间件)"
    
    echo "2.1 桌面浏览器请求:"
    curl -X GET "$BASE_URL/user/info/123" \
        -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36" \
        -w "\nHTTP状态码: %{http_code}\n\n"
    
    echo "2.2 移动端请求:"
    curl -X GET "$BASE_URL/user/info/456" \
        -H "User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) Mobile/15E148" \
        -w "\nHTTP状态码: %{http_code}\n\n"
    
    echo "2.3 爬虫请求:"
    curl -X GET "$BASE_URL/user/info/789" \
        -H "User-Agent: Googlebot/2.1 (+http://www.google.com/bot.html)" \
        -w "\nHTTP状态码: %{http_code}\n\n"
}

# 3. 测试认证中间件
test_auth() {
    print_test "测试认证中间件"
    
    echo "3.1 无Authorization头:"
    curl -X PUT "$BASE_URL/user/123" \
        -H "Content-Type: application/json" \
        -w "\nHTTP状态码: %{http_code}\n\n"
    
    echo "3.2 错误的Authorization格式:"
    curl -X PUT "$BASE_URL/user/123" \
        -H "Authorization: Token invalid-format" \
        -H "Content-Type: application/json" \
        -w "\nHTTP状态码: %{http_code}\n\n"
    
    echo "3.3 无效的token:"
    curl -X PUT "$BASE_URL/user/123" \
        -H "Authorization: Bearer invalid-token" \
        -H "Content-Type: application/json" \
        -w "\nHTTP状态码: %{http_code}\n\n"
    
    echo "3.4 有效的用户token:"
    curl -X PUT "$BASE_URL/user/123" \
        -H "Authorization: Bearer valid-token-123" \
        -H "Content-Type: application/json" \
        -w "\nHTTP状态码: %{http_code}\n\n"
    
    echo "3.5 管理员token访问其他用户:"
    curl -X PUT "$BASE_URL/user/999" \
        -H "Authorization: Bearer admin-token-456" \
        -H "Content-Type: application/json" \
        -w "\nHTTP状态码: %{http_code}\n\n"
}

# 4. 测试限流中间件
test_ratelimit() {
    print_test "测试限流中间件"
    
    echo "4.1 正常请求频率:"
    for i in {1..3}; do
        echo "请求 $i:"
        curl -X PUT "$BASE_URL/user/123" \
            -H "Authorization: Bearer valid-token-123" \
            -H "Content-Type: application/json" \
            -w " (状态码: %{http_code})\n"
        sleep 0.1
    done
    echo ""
    
    echo "4.2 高频请求测试（可能触发限流）:"
    echo "快速发送10个请求..."
    for i in {1..10}; do
        curl -s -X PUT "$BASE_URL/user/123" \
            -H "Authorization: Bearer valid-token-123" \
            -H "Content-Type: application/json" \
            -w "请求 $i: %{http_code}\n"
    done
    echo ""
}

# 5. 测试CORS（跨域）
test_cors() {
    print_test "测试CORS功能"
    
    echo "5.1 预检请求 (OPTIONS):"
    curl -X OPTIONS "$BASE_URL/user/info/123" \
        -H "Origin: http://example.com" \
        -H "Access-Control-Request-Method: GET" \
        -H "Access-Control-Request-Headers: Authorization" \
        -v 2>&1 | grep -E "(< HTTP|< Access-Control)"
    echo ""
    
    echo "5.2 带Origin的GET请求:"
    curl -X GET "$BASE_URL/user/info/123" \
        -H "Origin: http://example.com" \
        -H "User-Agent: Browser/1.0" \
        -v 2>&1 | grep -E "(< HTTP|< Access-Control|< X-)"
    echo ""
}

# 6. 测试安全头
test_security() {
    print_test "测试安全中间件"
    
    echo "6.1 检查安全响应头:"
    curl -I "$BASE_URL/user/info/123" \
        -H "User-Agent: Security-Test/1.0" 2>/dev/null | \
        grep -E "(X-Content-Type-Options|X-Frame-Options|X-XSS-Protection|Strict-Transport-Security|X-API-Version|X-Powered-By)"
    echo ""
    
    echo "6.2 测试请求体大小限制 (发送大请求):"
    # 创建一个大的JSON负载（超过10MB）
    large_data=$(printf '{"data":"%*s"}' 10485760 '')
    echo "$large_data" | curl -X POST "$BASE_URL/user/login" \
        -H "Content-Type: application/json" \
        -d @- \
        -w "\nHTTP状态码: %{http_code}\n\n"
}

# 7. 测试日志中间件（通过请求不同类型的API）
test_logging() {
    print_test "测试日志中间件"
    
    echo "7.1 GET请求（查看服务器日志中的请求记录）:"
    curl -X GET "$BASE_URL/user/info/123" \
        -H "User-Agent: LogTest-Client/1.0" \
        -w "\nHTTP状态码: %{http_code}\n"
    
    echo "7.2 POST请求（带请求体，查看服务器日志中的请求体记录）:"
    curl -X POST "$BASE_URL/user/login" \
        -H "Content-Type: application/json" \
        -H "User-Agent: LogTest-Client/1.0" \
        -d '{"username":"admin456","password":"admin123"}' \
        -w "\nHTTP状态码: %{http_code}\n\n"
}

# 主测试流程
main() {
    echo "开始测试 go-zero 中间件功能..."
    echo "请确保服务已在 http://localhost:8888 启动"
    echo ""
    
    # 检查服务状态
    check_service
    
    # 执行各项测试
    test_login
    sleep 1
    
    test_userinfo
    sleep 1
    
    test_auth
    sleep 1
    
    test_ratelimit
    sleep 1
    
    test_cors
    sleep 1
    
    test_security
    sleep 1
    
    test_logging
    
    echo ""
    print_success "所有测试完成！"
    echo ""
    echo "请查看服务器控制台输出，观察各个中间件的日志记录。"
    echo "您应该能看到："
    echo "  - [全局中间件] 的请求处理日志"
    echo "  - [安全中间件] 的安全检查日志"  
    echo "  - [UserAgent中间件] 的用户代理记录"
    echo "  - [认证中间件] 的认证过程日志"
    echo "  - [日志中间件] 的详细请求响应记录"
    echo "  - [限流中间件] 的限流处理日志"
}

# 如果脚本被直接执行，运行主函数
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi