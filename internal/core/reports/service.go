package reports

import (
	"context"
	"fmt"

	"github.com/ivmello/kakebo-go-api/internal/core/reports/dto"
	"github.com/ivmello/kakebo-go-api/internal/core/transactions"
	"github.com/ivmello/kakebo-go-api/internal/core/transactions/entity"
	"github.com/ivmello/kakebo-go-api/internal/utils"
)

type Service interface {
	Summarize(ctx context.Context, userId int) (dto.SummarizeOutput, error)
}

type service struct {
	repo transactions.Repository
}

func NewService(repo transactions.Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) Summarize(ctx context.Context, userId int) (dto.SummarizeOutput, error) {
	transactionsList, err := s.repo.GetAllUserTransactions(ctx, userId)
	fmt.Println(transactionsList, userId)
	if err != nil {
		return dto.SummarizeOutput{}, err
	}
	var total, debits, credits int
	for _, t := range transactionsList {
		if t.TransactionType == entity.DEBIT {
			debits += t.Amount
		} else {
			credits += t.Amount
		}
	}
	total = credits - debits
	return dto.SummarizeOutput{
		Total:           total,
		TotalFormated:   utils.FormatMoney(total, ".", ","),
		Debits:          debits,
		DebitsFormated:  utils.FormatMoney(debits, ".", ","),
		Credits:         credits,
		CreditsFormated: utils.FormatMoney(credits, ".", ","),
	}, nil
}
