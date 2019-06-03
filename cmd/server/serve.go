package server

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/open-identity/roles-manager/api"
	"github.com/open-identity/roles-manager/driver"
	"github.com/open-identity/utils/logrusx"
	"github.com/open-identity/utils/metricsx"
	"github.com/ory/graceful"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
)

// RunServe runs the Customer Onboarding API HTTP server
func RunServe(
	logger *logrus.Logger,
	version, commit string, date string,
) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		d := driver.NewDefaultDriver(
			logrusx.New(),
			version, commit, date,
		)

		d.CallRegistry()

		router := xmux.New()
		api.RegisterWithMux(router, d.Registry())

		c := xhandler.Chain{}

		// Append a context-aware middleware handler
		c.UseC(xhandler.CloseHandler)

		// Another context-aware middleware handler
		c.UseC(xhandler.TimeoutHandler(10 * time.Second))

		configuration := d.Configuration()
		c.UseC(metricsx.InstrumentHandlerWithIncludePath(d.Registry().Metrics(), "/api"))

		n := negroni.New()
		n.Use(negronilogrus.NewMiddlewareFromLogger(logger, "custm"))

		if tracer := d.Registry().Tracer(); tracer.IsLoaded() {
			n.Use(tracer)
		}
		n.UseHandler(c.Handler(router))

		server := graceful.WithDefaults(&http.Server{
			Addr:    configuration.PublicListenOn(),
			Handler: n,
		})

		if d.Registry().Tracer().IsLoaded() {
			server.RegisterOnShutdown(d.Registry().Tracer().Close)
		}

		health := d.Registry().Health()

		if err := graceful.Graceful(func() error {
			logger.Println("Starting health check")
			if err := health.Start(); err != nil {
				return err
			}

			logger.Printf("Listening on http://%s", configuration.PublicListenOn())
			return server.ListenAndServe()
		}, func(ctx context.Context) error {
			var result error

			logger.Printf("Shutting down http server")
			if err := server.Shutdown(ctx); err != nil {
				result = multierror.Append(result, err)
			}

			logger.Println("Shutting down health check")
			if err := health.Stop(); err != nil {
				result = multierror.Append(result, err)
			}
			return result
		}); err != nil {
			logger.Fatalf("Unable to gracefully shutdown HTTP(s) server because %v", err)
			return
		}
	}
}
