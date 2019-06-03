package driver

import "github.com/open-identity/roles-manager/driver/configuration"

type Driver interface {
	Configuration() configuration.Provider
	Registry() Registry
	CallRegistry() Driver
}
