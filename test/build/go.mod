module build

go 1.23.0

toolchain go1.23.3

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240822170219-fc7c04adadcd

replace github.com/alibaba/loongsuite-go-agent => ../../

replace github.com/alibaba/loongsuite-go-agent/test/verifier => ../../test/verifier

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.35.0

replace go.opentelemetry.io/contrib/instrumentation/runtime => go.opentelemetry.io/contrib/instrumentation/runtime v0.60.0

replace go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.35.0

replace go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v1.35.0

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.35.0

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.35.0

replace go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.35.0

replace go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v1.35.0

replace go.opentelemetry.io/otel/exporters/prometheus => go.opentelemetry.io/otel/exporters/prometheus v0.57.0

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.35.0

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.35.0

replace go.opentelemetry.io/otel/exporters/zipkin => go.opentelemetry.io/otel/exporters/zipkin v1.35.0

replace go.opentelemetry.io/otel/sdk/metric => go.opentelemetry.io/otel/sdk/metric v1.35.0

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.35.0

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp => go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.35.0

replace google.golang.org/protobuf => google.golang.org/protobuf v1.35.2
