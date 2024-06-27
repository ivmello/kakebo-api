package provider

import (
	"github.com/ivmello/kakebo-go-api/internal/config"
	"github.com/ivmello/kakebo-go-api/internal/core/auth"
	"github.com/ivmello/kakebo-go-api/internal/core/reports"
	"github.com/ivmello/kakebo-go-api/internal/core/transactions"
	"github.com/ivmello/kakebo-go-api/internal/core/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Provider struct {
	cfg  *config.Config
	conn *pgxpool.Pool
}

func New(cfg *config.Config, conn *pgxpool.Pool) *Provider {
	return &Provider{
		cfg,
		conn,
	}
}

func (p *Provider) GetDB() *pgxpool.Pool {
	return p.conn
}

func (p *Provider) GetUserRepository() users.Repository {
	return users.NewRepository(p.conn)
}

func (p *Provider) GetUserService() users.Service {
	return users.NewService(p.GetUserRepository())
}

func (p *Provider) GetTransactionRepository() transactions.Repository {
	return transactions.NewRepository(p.conn)
}

func (p *Provider) GetTransactionService() transactions.Service {
	return transactions.NewService(p.GetTransactionRepository())
}

func (p *Provider) GetAuthService() auth.Service {
	return auth.NewService(p.cfg.JWTSecret, p.GetUserRepository())
}

func (p *Provider) GetReportService() reports.Service {
	return reports.NewService(p.GetTransactionRepository())
}
