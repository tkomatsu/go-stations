package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.createHandler(w, r)
	case "PUT":
		h.updateHandler(w, r)
	case "GET":
		h.readHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *TODOHandler) createHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody model.CreateTODORequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&reqBody); err != nil {
		http.Error(w, fmt.Errorf("json decode: %v", err).Error(), http.StatusInternalServerError)
		return
	}

	ret, err := h.Create(r.Context(), &reqBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(ret); err != nil {
		http.Error(w, fmt.Sprintf("json encode: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, buf.String())
}

func (h *TODOHandler) updateHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody model.UpdateTODORequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&reqBody); err != nil {
		http.Error(w, fmt.Errorf("json decode: %v", err).Error(), http.StatusInternalServerError)
		return
	}

	if reqBody.ID == 0 || reqBody.Subject == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ret, err := h.Update(r.Context(), &reqBody)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(ret); err != nil {
		log.Print("json encode: ", err)
		http.Error(w, fmt.Sprintf("json encode: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, buf.String())
}

func (h *TODOHandler) readHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var reqBody model.ReadTODORequest
	if params.Get("prev_id") != "" {
		prevId, err := strconv.Atoi(params.Get("prev_id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("get prev_id: %v", err), http.StatusInternalServerError)
		}
		reqBody.PrevID = int64(prevId)
	}
	if params.Get("size") != "" {
		size, err := strconv.Atoi(params.Get("size"))
		if err != nil {
			http.Error(w, fmt.Sprintf("get size: %v", err), http.StatusInternalServerError)
		}
		reqBody.Size = int64(size)
	}

	ret, err := h.Read(r.Context(), &reqBody)
	if err != nil {
		log.Print(err)
		http.NotFound(w, r)
		return
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(ret); err != nil {
		log.Print("json encode: ", err)
		http.Error(w, fmt.Sprintf("json encode: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, buf.String())
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	if req.Subject == "" {
		return nil, errors.New("subject empty")
	}

	ret, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: ret}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	ret, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: ret}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	ret, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: ret}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
