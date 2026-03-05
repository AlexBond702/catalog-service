package rprocessor

import (
	"net/http"

	"github.com/gorilla/mux"

	rhandler "github.com/AlexBond702/catalog-service/internal/app/handler"
)

func vGenericRegHealthCheck(r *mux.Router, h rhandler.Health) {
	reg(r, http.MethodGet, "/health", http.HandlerFunc(h.LastCheck))
}

func handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
