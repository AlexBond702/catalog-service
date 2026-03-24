package httph

import "net/http"

type Error struct {
	Message string `json:"error"`
}

func ErrorApply(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", HeaderContentType)
	w.WriteHeader(code)
	if err := EncodeJSON(w, Error{Message: message}); err != nil {
		return
	}
}
