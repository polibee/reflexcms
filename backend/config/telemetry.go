package config

import (
	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("telemetry", map[string]any{
		"service": map[string]any{
			"name":        config.Env("APP_NAME", "goravel"),
			"version":     config.Env("APP_VERSION", ""),
			"environment": config.Env("APP_ENV", ""),
		},
		"resource":    map[string]any{},
		"propagators": config.Env("OTEL_PROPAGATORS", "tracecontext"),
		// Shutdown Timeout
		"shutdown_timeout": config.Env("OTEL_SHUTDOWN_TIMEOUT", "15s"),
		"traces": map[string]any{
			"exporter": config.Env("OTEL_TRACES_EXPORTER"),
			"sampler": map[string]any{
				"parent": config.Env("OTEL_TRACES_SAMPLER_PARENT", true),
				"type":   config.Env("OTEL_TRACES_SAMPLER_TYPE", "always_on"),
				"ratio":  config.Env("OTEL_TRACES_SAMPLER_RATIO", 0.05),
			},
			"processor": map[string]any{
				"type":     config.Env("OTEL_TRACE_PROCESSOR_TYPE", "batch"),
				"interval": config.Env("OTEL_TRACE_EXPORT_INTERVAL", "5s"),
				"timeout":  config.Env("OTEL_TRACE_EXPORT_TIMEOUT", "30s"),
			},
		},
		"metrics": map[string]any{
			"exporter": config.Env("OTEL_METRICS_EXPORTER"),
			"reader": map[string]any{
				"interval": config.Env("OTEL_METRIC_EXPORT_INTERVAL", "60s"),
				"timeout":  config.Env("OTEL_METRIC_EXPORT_TIMEOUT", "30s"),
			},
		},
		"logs": map[string]any{
			"exporter": config.Env("OTEL_LOGS_EXPORTER"),
			"processor": map[string]any{
				"type":     config.Env("OTEL_LOG_PROCESSOR_TYPE", "batch"),
				"interval": config.Env("OTEL_LOG_EXPORT_INTERVAL", "1s"),
				"timeout":  config.Env("OTEL_LOG_EXPORT_TIMEOUT", "30s"),
			},
		},
		"instrumentation": map[string]any{
			"http_server": map[string]any{
				"enabled":          config.Env("OTEL_HTTP_SERVER_ENABLED", true),
				"excluded_paths":   []string{"/healthz"},
				"excluded_methods": []string{"OPTIONS", "HEAD"},
			},
			"http_client": map[string]any{
				"enabled": config.Env("OTEL_HTTP_CLIENT_ENABLED", true),
			},
			"grpc_server": map[string]any{
				"enabled": config.Env("OTEL_GRPC_SERVER_ENABLED", true),
			},
			"grpc_client": map[string]any{
				"enabled": config.Env("OTEL_GRPC_CLIENT_ENABLED", true),
			},
			"database": map[string]any{
				"enabled": config.Env("OTEL_DATABASE_ENABLED", true),
			},
			"log": map[string]any{
				"enabled": config.Env("OTEL_LOG_ENABLED", true),
			},
		},
		"exporters": map[string]any{
			"otlptrace": map[string]any{
				"driver":      "otlp",
				"endpoint":    config.Env("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "localhost:4318"),
				"protocol":    config.Env("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL", "http/protobuf"),
				"insecure":    config.Env("OTEL_EXPORTER_OTLP_TRACES_INSECURE", true),
				"timeout":     config.Env("OTEL_EXPORTER_OTLP_TRACES_TIMEOUT", "10s"),
				"compression": config.Env("OTEL_EXPORTER_OTLP_TRACES_COMPRESSION", ""),
				"tls": map[string]any{
					"ca":   config.Env("OTEL_EXPORTER_OTLP_TRACES_CERTIFICATE", ""),
					"cert": config.Env("OTEL_EXPORTER_OTLP_TRACES_CLIENT_CERTIFICATE", ""),
					"key":  config.Env("OTEL_EXPORTER_OTLP_TRACES_CLIENT_KEY", ""),
				},
				"retry": map[string]any{
					"enabled":          true,
					"initial_interval": "5s",
					"max_interval":     "30s",
					"max_elapsed_time": "1m",
				},
			},
			"otlpmetric": map[string]any{
				"driver":              "otlp",
				"endpoint":            config.Env("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", "localhost:4318"),
				"protocol":            config.Env("OTEL_EXPORTER_OTLP_METRICS_PROTOCOL", "http/protobuf"),
				"insecure":            config.Env("OTEL_EXPORTER_OTLP_METRICS_INSECURE", true),
				"timeout":             config.Env("OTEL_EXPORTER_OTLP_METRICS_TIMEOUT", "10s"),
				"metric_temporality":  config.Env("OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY", "cumulative"),
				"compression":         config.Env("OTEL_EXPORTER_OTLP_METRICS_COMPRESSION", ""),
				"tls": map[string]any{
					"ca":   config.Env("OTEL_EXPORTER_OTLP_METRICS_CERTIFICATE", ""),
					"cert": config.Env("OTEL_EXPORTER_OTLP_METRICS_CLIENT_CERTIFICATE", ""),
					"key":  config.Env("OTEL_EXPORTER_OTLP_METRICS_CLIENT_KEY", ""),
				},
				"retry": map[string]any{
					"enabled":          true,
					"initial_interval": "5s",
					"max_interval":     "30s",
					"max_elapsed_time": "1m",
				},
			},
			"otlplog": map[string]any{
				"driver":      "otlp",
				"endpoint":    config.Env("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT", "localhost:4318"),
				"protocol":    config.Env("OTEL_EXPORTER_OTLP_LOGS_PROTOCOL", "http/protobuf"),
				"insecure":    config.Env("OTEL_EXPORTER_OTLP_LOGS_INSECURE", true),
				"timeout":     config.Env("OTEL_EXPORTER_OTLP_LOGS_TIMEOUT", "10s"),
				"compression": config.Env("OTEL_EXPORTER_OTLP_LOGS_COMPRESSION", ""),
				"tls": map[string]any{
					"ca":   config.Env("OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE", ""),
					"cert": config.Env("OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE", ""),
					"key":  config.Env("OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY", ""),
				},
				"retry": map[string]any{
					"enabled":          true,
					"initial_interval": "5s",
					"max_interval":     "30s",
					"max_elapsed_time": "1m",
				},
			},
			"console": map[string]any{
				"driver":       "console",
				"pretty_print": false,
			},
			"custom": map[string]any{
				"driver": "custom",
			},
		},
	})
}
