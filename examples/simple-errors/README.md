## 运行示例

- 运行示例：`go run main.go`
- 在另一个终端请求测试：
```bash
# 请求英文响应
curl -X POST "http://localhost:8080/api/v1/user/login" \
     -H "Accept-Language: en-US" \
     -H "Content-Type: application/json" \
     -d '{"username": "test", "password": "wrong"}'

# 请求中文响应
curl -X POST "http://localhost:8080/api/v1/user/login" \
     -H "Accept-Language: zh-CN" \
     -H "Content-Type: application/json" \
     -d '{"username": "test", "password": "wrong"}'
```