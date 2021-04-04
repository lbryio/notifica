package action

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lbryio/notifica/app/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// StartNotifica starts up Notifica!
func StartNotifica(port int) {
	routes := GetRoutes()

	httpServeMux := http.NewServeMux()
	httpServeMux.Handle(promPath, promBasicAuthWrapper(promhttp.Handler()))
	routes.Each(func(pattern string, handler http.Handler) {
		httpServeMux.Handle(pattern, handler)
	})

	mux := http.Handler(httpServeMux)

	for _, middleware := range []func(h http.Handler) http.Handler{
		promRequestHandler,
	} {
		mux = middleware(mux)
	}

	configureAPIServer()

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
		//https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		//https://blog.cloudflare.com/exposing-go-on-the-internet/
		ReadTimeout: 5 * time.Second,
		//WriteTimeout: 10 * time.Second, // cant use this yet, since some of our responses take a long time (e.g. sending emails)
		IdleTimeout: 120 * time.Second,
	}
	logrus.Infof("Listening on port %v", port)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			//Normal graceful shutdown error
			if err.Error() == "http: Server closed" {
				logrus.Info(err)
			} else {
				logrus.Fatal(err)
			}
		}
	}()
	//Wait for shutdown signal, then shutdown api server. This will wait for all connections to finish.
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interruptChan
	logrus.Debug("Shutting down API server...")
	err := srv.Shutdown(context.Background())
	if err != nil {
		logrus.Error("Error shutting down server: ", err)
	}
	logrus.Debug("Rick Reports is shutting down...")

}

const promPath = "/metrics"

func promRequestHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimLeft(r.URL.Path, "/")
		if path != strings.TrimLeft(promPath, "/") {
			metrics.UserLoadOverall.Inc()
			defer metrics.UserLoadOverall.Dec()
			metrics.UserLoadByAPI.WithLabelValues(path).Inc()
			defer metrics.UserLoadByAPI.WithLabelValues(path).Dec()
			apiStart := time.Now()
			h.ServeHTTP(w, r)
			duration := time.Since(apiStart).Seconds()
			metrics.Durations.WithLabelValues(path).Observe(duration)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func promBasicAuthWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "authentication required", http.StatusBadRequest)
			return
		}
		if user == "prom" && pass == "prom-notifica-access" {
			h.ServeHTTP(w, r)
		} else {
			http.Error(w, "invalid username or password", http.StatusForbidden)
		}
	})
}
