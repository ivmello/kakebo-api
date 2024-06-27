package transactions

import (
	"encoding/json"
	"io"
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

func (h *handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	var input CreateTransactionInput
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
	output, err := h.service.CreateTransaction(ctx, userId, input)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}

func (h *handler) ImportTransactionsFromFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.JSONErrorResponse(w, "Erro ao fazer o parse do form", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		utils.JSONErrorResponse(w, "Erro ao obter o arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()
	output, err := h.service.ImportTransactionsFromCSV(ctx, userId, file)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}

func (h *handler) ListAllUserTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	output, err := h.service.GetAllUserTransactions(ctx, userId)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}

func (h *handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	transactionId := r.PathValue("id")
	output, err := h.service.GetTransaction(ctx, userId, transactionId)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}

func (h *handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	transactionId := r.PathValue("id")
	err := h.service.DeleteTransaction(ctx, userId, transactionId)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, nil, http.StatusNoContent)
}
