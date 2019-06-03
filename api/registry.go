package api

import "github.com/rs/xmux"

type Registry interface {
}

type apHandler interface {
	Register(mux *xmux.Mux)
}

func RegisterWithMux(mux *xmux.Mux, registry Registry) {
	for _, s := range []apHandler{} {
		s.Register(mux)
	}
}
