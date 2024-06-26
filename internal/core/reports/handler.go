package reports

import (
	"net/http"

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
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	output, err := h.service.Summarize(ctx, userId)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}
