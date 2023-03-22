package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/NenadRadulovic/go-api/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var defaultStopTimeout = time.Second * 30

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := e.handler(w, r); err != nil {
		logrus.WithError(err).Error("Could not process request")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}
}

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func NewAPIServer(addr string, storage *storage.Storage) (*APIServer, error) {
	if addr == "" {
		return nil, errors.New("addr cannot be blank")
	}

	return &APIServer{
		addr:    addr,
		storage: storage,
	}, nil
}

// Start starts a server with a stop channel
func (s *APIServer) Start(stop <-chan struct{}) error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.initRoutes(),
	}
	go func() {
		logrus.WithField("addr", srv.Addr).Info("starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), defaultStopTimeout)
	defer cancel()

	logrus.WithField("timeout", defaultStopTimeout).Info("stopping server")
	return srv.Shutdown(ctx)
}

func (s *APIServer) initRoutes() http.Handler {
	router := mux.NewRouter()
	router.Use(CommonMiddleware)

	router.HandleFunc("/", s.defaultRoute)
	s.initApiRoutes(router)
	s.initAuthRouter(router)
	return router
}

func (s *APIServer) defaultRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello GO WORLD"))
}

func (s *APIServer) getUrlParams(r *http.Request) map[string]string {
	vars := mux.Vars(r)
	return vars
}
