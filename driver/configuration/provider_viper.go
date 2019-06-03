package configuration

import (
	"fmt"
	"strings"

	"github.com/open-identity/utils/tracing"
	"github.com/open-identity/utils/viperx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ViperProvider struct {
	l logrus.FieldLogger
}

const (
	ViperKeyDSN                                  = "dsn"
	ViperKeyPublicListenOnHost                   = "serve.public.host"
	ViperKeyPublicListenOnPort                   = "serve.public.port"
	ViperTracingServiceName                      = "tracing.service_name"
	ViperTracingProvider                         = "tracing.provider"
	ViperTracingProvidersJaegerAgentAddress      = "tracing.providers.jaeger.local_agent_address"
	ViperTracingProvidersJaegerSamplingType      = "tracing.providers.jaeger.sampling.type"
	ViperTracingProvidersJaegerSamplingValue     = "tracing.providers.jaeger.sampling.value"
	ViperTracingProvidersJaegerSamplingServerUrl = "tracing.providers.jaeger.sampling.server_url"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func NewViperProvider(l logrus.FieldLogger) Provider {

	return &ViperProvider{
		l: l,
	}
}
func (v *ViperProvider) getAddress(address string, port int) string {
	if strings.HasPrefix(address, "unix:") {
		return address
	}
	return fmt.Sprintf("%s:%d", address, port)
}

func (v *ViperProvider) publicHost() string {
	return viperx.GetString(v.l, ViperKeyPublicListenOnHost, "")
}

func (v *ViperProvider) publicPort() int {
	return viperx.GetInt(v.l, ViperKeyPublicListenOnPort, 80)
}

func (v *ViperProvider) DSN() string {
	return viperx.GetString(v.l, ViperKeyDSN, "")
}

func (v *ViperProvider) PublicListenOn() string {
	return v.getAddress(v.publicHost(), v.publicPort())
}

func (v *ViperProvider) TracingServiceName() string {
	return viperx.GetString(v.l, ViperTracingServiceName, "Config Inventory")
}

func (v *ViperProvider) TracingProvider() string {
	return viperx.GetString(v.l, ViperTracingProvider, "")
}

func (v *ViperProvider) TracingJaegerConfig() *tracing.JaegerConfig {
	return &tracing.JaegerConfig{
		LocalAgentHostPort: viperx.GetString(v.l, ViperTracingProvidersJaegerAgentAddress, ""),
		SamplerType:        viperx.GetString(v.l, ViperTracingProvidersJaegerSamplingType, "const"),
		SamplerValue:       viperx.GetFloat64(v.l, ViperTracingProvidersJaegerSamplingValue, float64(1)),
		SamplerServerURL:   viperx.GetString(v.l, ViperTracingProvidersJaegerSamplingServerUrl, ""),
	}
}
