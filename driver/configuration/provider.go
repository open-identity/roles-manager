package configuration

import "github.com/open-identity/utils/tracing"

type Provider interface {
	DSN() string

	PublicListenOn() string

	TracingServiceName() string
	TracingProvider() string
	TracingJaegerConfig() *tracing.JaegerConfig
}
