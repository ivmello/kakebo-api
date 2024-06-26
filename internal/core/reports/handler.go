package reports

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ivmello/kakebo-go-api/internal/core/transactions"
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

func (h *handler) Summarize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input transactions.TransactionFilter
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &input)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	output, err := h.service.Summarize(ctx, userId, input)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}
