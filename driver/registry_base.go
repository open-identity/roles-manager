package driver

import (
	"github.com/InVisionApp/go-health"
	"github.com/open-identity/roles-manager/driver/configuration"
	"github.com/open-identity/utils/healthx"
	"github.com/open-identity/utils/logrusx"
	"github.com/open-identity/utils/metricsx"
	"github.com/open-identity/utils/tracing"
	"github.com/ory/herodot"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type RegistryBase struct {
	l            logrus.FieldLogger
	c            configuration.Provider
	writer       herodot.Writer
	buildVersion string
	buildHash    string
	buildDate    string
	r            Registry
	trc          *tracing.Tracer
	mr           *prometheus.Registry
	health       health.IHealth
}

func (m *RegistryBase) with(r Registry) *RegistryBase {
	m.r = r
	return m
}

func (m *RegistryBase) WithBuildInfo(version, hash, date string) Registry {
	m.buildVersion = version
	m.buildHash = hash
	m.buildDate = date
	return m.r
}

func (m *RegistryBase) BuildVersion() string {
	return m.buildVersion
}

func (m *RegistryBase) BuildDate() string {
	return m.buildDate
}

func (m *RegistryBase) BuildHash() string {
	return m.buildHash
}

func (m *RegistryBase) WithConfig(c configuration.Provider) Registry {
	m.c = c
	return m.r
}

func (m *RegistryBase) WithLogger(l logrus.FieldLogger) Registry {
	m.l = l
	return m.r
}

func (m *RegistryBase) Writer() herodot.Writer {
	if m.writer == nil {
		h := herodot.NewJSONWriter(m.Logger())
		h.ErrorEnhancer = nil
		m.writer = h
	}
	return m.writer
}

func (m *RegistryBase) Logger() logrus.FieldLogger {
	if m.l == nil {
		m.l = logrusx.New()
	}
	return m.l
}

func (m *RegistryBase) Health() health.IHealth {
	if m.health == nil {
		h := health.New()
		h.Logger = healthx.NewShim(m.Logger())
		m.health = h
	}

	return m.health
}

func (m *RegistryBase) Metrics() *prometheus.Registry {
	if m.mr == nil {
		m.mr = prometheus.NewRegistry()
		metricsx.RegisterSystemMetrics(m.mr)
	}
	return m.mr
}

func (m *RegistryBase) Tracer() *tracing.Tracer {
	if m.trc == nil {
		m.trc = &tracing.Tracer{
			ServiceName:  m.c.TracingServiceName(),
			JaegerConfig: m.c.TracingJaegerConfig(),
			Provider:     m.c.TracingProvider(),
			Logger:       m.Logger(),
		}

		if err := m.trc.Setup(); err != nil {
			m.Logger().WithError(err).Fatalf("Unable to initialize Tracer.")
		}
	}

	return m.trc
}
