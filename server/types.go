package server

import (
	"net/http"

	"github.com/NenadRadulovic/go-api/storage"
)

type APIServer struct {
	addr    string
	storage *storage.Storage
}

type Endpoint struct {
	handler EndpointFunc
}

type EndpointFunc func(w http.ResponseWriter, r *http.Request) error
