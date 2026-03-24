package httph

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("failed to encode")
	}
}

func SendEmpty(w http.ResponseWriter, status int) {
	SendRaw(w, status, "", nil)
}

func SendError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
		log.Println("failed to encode")
	}
}

func SendEncodeWithMIME(w http.ResponseWriter, r *http.Request, statusCode int, mimeType string, obj any) {
	SendRaw(w, statusCode, mimeType, nil)
	if err := EncodeJSON(w, obj); err != nil {
		ErrorApply(w, statusCode, err.Error())
		return
	}
}

func SendEncoded(w http.ResponseWriter, r *http.Request, statusCode int, obj any) {
	SendEncodeWithMIME(w, r, statusCode, MIMEApplicationJSONCharsetUTF8, obj)
}
