package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw := &model.HealthzResponse{Message:"OK"}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(raw); err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, buf.String())
}
