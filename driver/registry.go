package driver

import (
	"github.com/open-identity/roles-manager/driver/configuration"
	"github.com/open-identity/utils/dbal"
	"github.com/open-identity/utils/driver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Registry interface {
	Init() error

	WithConfig(c configuration.Provider) Registry
	WithLogger(l logrus.FieldLogger) Registry
	WithBuildInfo(version, hash, date string) Registry

	BuildVersion() string
	BuildDate() string
	BuildHash() string

	driver.RegistryLogger
	driver.RegistryWriter
	driver.RegistryHealth
	driver.RegistryMetrics
	driver.RegistryTracer
}

func NewRegistry(c configuration.Provider) (Registry, error) {
	d, err := dbal.GetDriverFor(c.DSN())
	if err != nil {
		return nil, err
	}

	registry, ok := d.(Registry)
	if !ok {
		return nil, errors.Errorf("driver of type %T does not implement interface Registry", d)
	}

	registry = registry.WithConfig(c)

	if err := registry.Init(); err != nil {
		return nil, err
	}

	return registry, nil
}

func CallRegistry(r Registry) {
	r.Health()
	r.Metrics()
	r.Tracer()
	r.Logger()
	r.Writer()
}
