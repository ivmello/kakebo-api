package auth

import (
	"encoding/json"
	"net/http"

	"github.com/ivmello/kakebo-go-api/internal/utils"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service,
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	input := LoginInput{}
	json.NewDecoder(r.Body).Decode(&input)
	output, err := h.service.Login(input)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}
