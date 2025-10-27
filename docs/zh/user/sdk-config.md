# OpenTelemetry SDK 配置

`otel`工具除了自动埋点外，还会注入配置代码，在应用启动时会初始化 OpenTelemetry SDK，使用以下环境变量可以改变 OpenTelemetry SDK 的行为。

- `OTEL_SERVICE_NAME`: 为您的应用指定服务名称。
- `OTEL_TRACES_EXPORTER`: 指定链路导出器。支持的值: `none`, `console`, `zipkin`, `otlp`。支持使用逗号分隔指定多个导出器（例如 `console,otlp`）。默认为 `otlp`。
- `OTEL_METRICS_EXPORTER`: 指定指标导出器。支持的值: `none`, `console`, `prometheus`, `otlp`。支持使用逗号分隔指定多个导出器（例如 `console,otlp`）。默认为 `otlp`。
- `OTEL_EXPORTER_OTLP_PROTOCOL`: 指定 OTLP 协议，用于链路和指标。支持的值: `http/protobuf` (默认), `grpc`。
- `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL`: 指定用于链路的 OTLP 协议，会覆盖 `OTEL_EXPORTER_OTLP_PROTOCOL` 的设置。支持的值: `http/protobuf` (默认), `grpc`。
- `OTEL_EXPORTER_OTLP_ENDPOINT`: 指定 OTLP 导出器的通用端点。
- `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`: 指定 OTLP 链路导出器的端点。
- `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT`: 指定 OTLP 指标导出器的端点。
- `OTEL_EXPORTER_OTLP_HEADERS`: 为所有 OTLP 导出器指定请求头 (例如, `key1=value1,key2=value2`)。
- `OTEL_EXPORTER_PROMETHEUS_PORT`: 当 `OTEL_METRICS_EXPORTER` 设置为 `prometheus` 时，指定 Prometheus 导出器的端口。默认为 `9464`。
- `OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY_PREFERENCE`: 指定指标的聚合时间性偏好（不区分大小写）。支持的值:
  - `cumulative` (默认): 所有指标类型都使用累积时间性
  - `delta`: Counter、Asynchronous Counter 和 Histogram 使用增量时间性；UpDownCounter 和 Asynchronous UpDownCounter 使用累积时间性
  - `lowmemory`: Synchronous Counter 和 Histogram 使用增量时间性；其他类型使用累积时间性（低内存模式）
- `OTEL_TRACE_SAMPLER`: 指定链路采样器。0.0 到 1.0 之间的浮点数会设置一个基于比率的采样器。小于等于 0 的值将永不采样，大于等于 1 的值将始终采样。默认是基于父级的采样器，并且始终采样。
