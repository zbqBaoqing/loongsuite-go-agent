![](../public/anim-logo.svg)

**龙蜥Go-Agent** 为希望利用OpenTelemetry实现有效可观察性的Golang应用程序提供了一个自动化的解决方案。目标应用程序无需更改代码，埋点在编译时完成。只需在`go build`前加上`otel`前缀即可开始使用 :rocket:

# 安装

### 预编译二进制文件

- [![下载](https://shields.io/badge/-Linux_AMD64-blue?logo=ubuntu)](https://github.com/alibaba/loongsuite-go-agent/releases/latest/download/otel-linux-amd64)
- [![下载](https://shields.io/badge/-Linux_ARM64-blue?logo=ubuntu)](https://github.com/alibaba/loongsuite-go-agent/releases/latest/download/otel-linux-arm64)
- [![下载](https://shields.io/badge/-MacOS_AMD64-blue?logo=apple)](https://github.com/alibaba/loongsuite-go-agent/releases/latest/download/otel-darwin-amd64)
- [![下载](https://shields.io/badge/-MacOS_ARM64-blue?logo=apple)](https://github.com/alibaba/loongsuite-go-agent/releases/latest/download/otel-darwin-arm64)
- [![下载](https://shields.io/badge/-Windows_AMD64-blue?logo=wine)](https://github.com/alibaba/loongsuite-go-agent/releases/latest/download/otel-windows-amd64.exe)

**这是安装该工具的推荐方法。**

### 通过Bash安装
对于Linux和MacOS用户，以下脚本默认会将`otel`安装在`/usr/local/bin/otel`：
```bash
$ sudo curl -fsSL https://cdn.jsdelivr.net/gh/alibaba/loongsuite-go-agent@main/install.sh | sudo bash
```

### 从源码构建

```bash
$ make         # 仅构建
$ make install # 构建并安装
```

# 快速开始

确保工具已安装：
```bash
$ # 你也可以使用 "otel-linux-amd64" 代替 "otel"
$ otel version
```

只需在`go build`前加上`otel`前缀来构建你的项目：

```bash
$ otel go build
$ otel go build -o app cmd/app
$ otel go build -gcflags="-m" cmd/app
```

整个过程就是这样！该工具将自动为你的代码注入OpenTelemetry，你就可以开始观察你的应用程序了。 :telescope:

`otel`工具的详细用法可以在[**用法**](./user/config.md)中找到。

> [!NOTE]
> 如果你在`go build`能正常工作的情况下发现任何编译失败，这很可能是一个bug。
> 请随时在[GitHub Issues](https://github.com/alibaba/loongsuite-go-agent/issues)上提交一个bug
> 来帮助我们改进这个项目。

# 示例

- [demo](https://github.com/alibaba/loongsuite-go-agent/tree/main/example/demo) - 带有OpenTelemetry追踪和指标的端到端示例
- [zap logging](https://github.com/alibaba/loongsuite-go-agent/tree/main/example/log) - `github.com/uber-go/zap`日志记录的自动埋点
- [benchmark](https://github.com/alibaba/loongsuite-go-agent/tree/main/example/benchmark) - 性能测试和开销测量
- [sql injection](https://github.com/alibaba/loongsuite-go-agent/tree/main/example/sqlinject) - 用于SQL注入检测的自定义代码注入
- [nethttp](https://github.com/alibaba/loongsuite-go-agent/tree/main/example/nethttp) - 带有请求/响应埋点的HTTP监控
- [kratos-demo](https://github.com/alibaba/loongsuite-go-agent/tree/main/example/kratos-demo) - 与Kratos框架的集成

# 支持的库
<details>
 <summary>支持的库列表</summary>

| 库 | 仓库地址 | 最低版本 | 最高版本 |
|---------------| ---------------------------------------------- |----------------------|-----------------------|
| database/sql | https://pkg.go.dev/database/sql | - | - |
| dubbo-go | https://github.com/apache/dubbo-go | v3.3.0 | - |
| echo | https://github.com/labstack/echo | v4.0.0 | v4.12.0 |
| eino | https://github.com/cloudwego/eino | v0.3.51 | - |
| elasticsearch | https://github.com/elastic/go-elasticsearch | v8.4.0 | v8.15.0 |
| fasthttp | https://github.com/valyala/fasthttp | v1.45.0 | v1.63.0 |
| fiber | https://github.com/gofiber/fiber | v2.43.0 | v2.52.8 |
| gin | https://github.com/gin-gonic/gin | v1.7.0 | v1.10.0 |
| go-redis | https://github.com/redis/go-redis | v9.0.5 | v9.5.1 |
| go-redis v8 | https://github.com/redis/go-redis | v8.11.0 | v8.11.5 |
| gomicro | https://github.com/micro/go-micro | v5.0.0 | v5.3.0 |
| gorestful | https://github.com/emicklei/go-restful | v3.7.0 | v3.12.1 |
| gorm | https://github.com/go-gorm/gorm | v1.22.0 | v1.25.9 |
| grpc | https://google.golang.org/grpc | v1.44.0 | - |
| hertz | https://github.com/cloudwego/hertz | v0.8.0 | - |
| iris | https://github.com/kataras/iris | v12.2.0 | v12.2.11 |
| client-go | https://github.com/kubernetes/client-go | v0.33.3 | - |
| kitex | https://github.com/cloudwego/kitex | v0.5.1 | v0.11.3 |
| kratos | https://github.com/go-kratos/kratos | v2.6.3 | v2.8.4 |
| langchaingo | https://github.com/tmc/langchaingo | v0.1.13 | v0.1.13 |
| log | https://pkg.go.dev/log | - | - |
| logrus | https://github.com/sirupsen/logrus | v1.5.0 | v1.9.3 |
| mongodb | https://github.com/mongodb/mongo-go-driver | v1.11.1 | v1.15.1 |
| mux | https://github.com/gorilla/mux | v1.3.0 | v1.8.1 |
| nacos | https://github.com/nacos-group/nacos-sdk-go/v2 | v2.0.0 | v2.2.7 |
| net/http | https://pkg.go.dev/net/http | - | - |
| ollama | https://github.com/ollama/ollama | v0.3.14 | - |
| redigo | https://github.com/gomodule/redigo | v1.9.0 | v1.9.2 |
| sentinel | https://github.com/alibaba/sentinel-golang | v1.0.4 | - |
| slog | https://pkg.go.dev/log/slog | - | - |
| trpc-go | https://github.com/trpc-group/trpc-go | v1.0.0 | v1.0.3 |
| zap | https://github.com/uber-go/zap | v1.20.0 | v1.27.0 |
| zerolog | https://github.com/rs/zerolog | v1.10.0 | v1.33.0 |
| go-kit/log | https://github.com/go-kit/log | v0.1.0 | v0.2.1 |
| pg | https://github.com/go-pg/pg | v1.10.0 | v1.14.0 |
| gocql | https://github.com/gocql/gocql | v1.3.0 | v1.7.0 |
| sqlx | https://github.com/jmoiron/sqlx | v1.3.0 | v1.4.0 |

</details>

我们正在逐步开源我们支持的库，<kbd>非常欢迎</kbd>您的贡献。

> [!IMPORTANT]
> 你期望的框架不在列表中？别担心，你可以轻松地将你的代码注入任何未正式支持的框架/库。
>
> 请参考[此文档](./dev/overview.md)开始。

# 社区

我们期待您的反馈和建议。您可以加入我们的[钉钉群](https://qr.dingtalk.com/action/joingroup?code=v1,k1,PBuICMTDvdh0En8MrVbHBYTGUcPXJ/NdT6JKCZ8CG+w=&_dt_no_comment=1&origin=11)与我们交流。

| 钉钉 | Star历史 |
| :---: | :---: |
| <img src="../public/dingtalk.png" height="200" /> | <img src="https://api.star-history.com/svg?repos=alibaba/loongsuite-go-agent&type=Date" height="200" /> |
