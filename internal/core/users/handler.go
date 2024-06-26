package users

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ivmello/kakebo-go-api/internal/core/users/dto"
	"github.com/ivmello/kakebo-go-api/internal/utils"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateUserInput
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	output, err := h.service.CreateUser(ctx, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	fmt.Println(userId)
	if userId == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var input dto.UpdateUserInput
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.service.UpdateUser(ctx, userId, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}
