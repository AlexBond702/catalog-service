package rhealth

import (
	"io"
	"log"
	"net/http"

	"github.com/AlexBond702/catalog-service/internal/app/handler"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	if _, err := io.ReadAll(r.Body); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	if _, err := w.Write([]byte("ok")); err != nil {
		log.Println("Failed to write")
	}
}
