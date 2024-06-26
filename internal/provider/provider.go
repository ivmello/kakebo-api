package provider

import (
	"github.com/ivmello/kakebo-go-api/internal/adapters/database"
	"github.com/ivmello/kakebo-go-api/internal/config"
	"github.com/ivmello/kakebo-go-api/internal/core/auth"
	"github.com/ivmello/kakebo-go-api/internal/core/transactions"
	"github.com/ivmello/kakebo-go-api/internal/core/users"
)

type Provider struct {
	cfg  *config.Config
	conn database.Connection
}

func New(cfg *config.Config, conn database.Connection) *Provider {
	return &Provider{
		cfg,
		conn,
	}
}

func (p *Provider) GetDB() database.Connection {
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
