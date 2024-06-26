package reports

import (
	"context"
	"time"

	"github.com/ivmello/kakebo-go-api/internal/core/transactions"
	"github.com/ivmello/kakebo-go-api/internal/utils"
)

type Service interface {
	Summarize(ctx context.Context, userId int, input transactions.TransactionFilter) (SummarizeOutput, error)
	SummarizeByMonth(ctx context.Context, userId int, input transactions.TransactionFilter) (SummarizeByMonthOutput, error)
}

type service struct {
	repo transactions.Repository
}

func NewService(repo transactions.Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) Summarize(ctx context.Context, userId int, input transactions.TransactionFilter) (SummarizeOutput, error) {
	transactionsList, err := s.repo.GetAllUserTransactionsByFilter(ctx, userId, input)
	if err != nil {
		return SummarizeOutput{}, err
	}
	var total, debits, credits int
	for _, t := range transactionsList {
		if t.TransactionType == transactions.DEBIT {
			debits += t.Amount
		} else {
			credits += t.Amount
		}
	}
	total = credits - debits
	return SummarizeOutput{
		Total:           total,
		TotalFormated:   utils.FormatMoney(total, ".", ","),
		Debits:          debits,
		DebitsFormated:  utils.FormatMoney(debits, ".", ","),
		Credits:         credits,
		CreditsFormated: utils.FormatMoney(credits, ".", ","),
	}, nil
}

func (s *service) SummarizeByMonth(ctx context.Context, userId int, input transactions.TransactionFilter) (SummarizeByMonthOutput, error) {
	transactionsList, err := s.repo.GetAllUserTransactionsByFilter(ctx, userId, input)
	if err != nil {
		return SummarizeByMonthOutput{}, err
	}
	var total, debits, credits int
	for _, t := range transactionsList {
		if t.TransactionType == transactions.DEBIT {
			debits += t.Amount
		} else {
			credits += t.Amount
		}
	}
	total = credits - debits
	year := input.Year
	if year == 0 {
		year = time.Now().Year()
	}
	return SummarizeByMonthOutput{
		Month:           input.Month,
		Year:            year,
		Total:           total,
		TotalFormated:   utils.FormatMoney(total, ".", ","),
		Debits:          debits,
		DebitsFormated:  utils.FormatMoney(debits, ".", ","),
		Credits:         credits,
		CreditsFormated: utils.FormatMoney(credits, ".", ","),
	}, nil
}
