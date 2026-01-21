# ginmetrics

Gin 框架的指标监控中间件，自动收集 HTTP 请求的各项指标数据。

## 安装

```bash
go get github.com/conflux-fans/ginmetrics
```

## 快速开始

```go
package main

import (
    "github.com/conflux-fans/ginmetrics"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    // 获取监控实例并应用到路由
    monitor := ginmetrics.GetMonitor()
    monitor.Use(r)
    
    // 你的路由
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    r.Run(":8080")
}
```

## 配置选项

### 设置慢请求阈值

默认慢请求阈值为 200 毫秒，超过此时间的请求会被计入慢请求指标。

```go
monitor := ginmetrics.GetMonitor()
monitor.SetSlowTime(500) // 设置为 500 毫秒
monitor.Use(r)
```

### 自定义指标注册表

```go
monitor := ginmetrics.GetMonitor()
monitor.SetRegistry(customRegistry)
monitor.Use(r)
```

## 收集的指标

中间件会自动收集以下指标：

- **gin_request_total** - 请求总数（按方法和路径分类）
- **gin_request_fail_total** - 失败请求数（HTTP 状态码 >= 400）
- **gin_request_uv_total** - 独立访客数（使用布隆过滤器统计 IP）
- **gin_request_body_total** - 请求体总大小（字节）
- **gin_response_body_total** - 响应体总大小（字节）
- **gin_request_duration** - 请求处理耗时（按方法和路径分类）
- **gin_slow_request_total** - 慢请求计数

## 数据上报

### 上报到 InfluxDB

```go
monitor := ginmetrics.GetMonitor()
monitor.Use(r)

// 每 10 秒上报一次到 InfluxDB
monitor.ReportToInfluxDB(
    10*time.Second,
    "http://localhost:8086",  // InfluxDB 地址
    "mydb",                   // 数据库名
    "username",                // 用户名
    "password",                // 密码
    "gin",                     // 命名空间
)
```

### 输出到日志

```go
import (
    "log"
    "os"
    "github.com/ethereum/go-ethereum/metrics"
)

monitor := ginmetrics.GetMonitor()
monitor.Use(r)

// 每 10 秒输出一次指标到日志
logger := log.New(os.Stderr, "metrics: ", log.Lmicroseconds)
monitor.ReportToLogger(10*time.Second, logger)
```

## 完整示例

```go
package main

import (
    "log"
    "os"
    "time"
    
    "github.com/conflux-fans/ginmetrics"
    "github.com/gin-gonic/gin"
    "github.com/ethereum/go-ethereum/metrics"
)

func main() {
    r := gin.Default()
    
    // 初始化监控
    monitor := ginmetrics.GetMonitor()
    monitor.SetSlowTime(300) // 设置慢请求阈值为 300ms
    monitor.Use(r)
    
    // 设置路由
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    r.GET("/slow", func(c *gin.Context) {
        time.Sleep(500 * time.Millisecond)
        c.JSON(200, gin.H{"message": "slow response"})
    })
    
    // 上报指标到日志（可选）
    logger := log.New(os.Stderr, "metrics: ", log.Lmicroseconds)
    monitor.ReportToLogger(10*time.Second, logger)
    
    // 或者上报到 InfluxDB（可选）
    // monitor.ReportToInfluxDB(
    //     10*time.Second,
    //     "http://localhost:8086",
    //     "mydb",
    //     "username",
    //     "password",
    //     "gin",
    // )
    
    r.Run(":8080")
}
```
