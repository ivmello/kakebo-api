package webserver

import (
	"net/http"
	"slices"

	"github.com/ivmello/kakebo-go-api/internal/adapters/webserver/middlewares"
	"github.com/ivmello/kakebo-go-api/internal/core/auth"
	"github.com/ivmello/kakebo-go-api/internal/core/reports"
	"github.com/ivmello/kakebo-go-api/internal/core/transactions"
	"github.com/ivmello/kakebo-go-api/internal/core/users"
)

type (
	middleware middlewares.Middleware
	router     struct {
		*http.ServeMux
		chain []middleware
	}
)

func NewRouter(mux *http.ServeMux, mx ...middleware) *router {
	return &router{
		ServeMux: mux,
		chain:    mx,
	}
}

func (r *router) Use(middleware ...middleware) {
	r.chain = append(r.chain, middleware...)
}

func (r *router) Group(fn func(r *router)) {
	fn(&router{
		ServeMux: r.ServeMux,
		chain:    slices.Clone(r.chain),
	})
}

func (r *router) Get(path string, handler http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodGet, path, handler, mx)
}

func (r *router) Post(path string, handler http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodPost, path, handler, mx)
}

func (r *router) Put(path string, handler http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodPut, path, handler, mx)
}

func (r *router) Delete(path string, handler http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodDelete, path, handler, mx)
}

func (r *router) handle(method, path string, handler http.HandlerFunc, mx []middleware) {
	r.Handle(method+" "+path, r.wrap(handler, mx))
}

func (r *router) wrap(handler http.HandlerFunc, mx []middleware) (out http.Handler) {
	out, mx = http.Handler(handler), append(slices.Clone(r.chain), mx...)
	slices.Reverse(mx)
	for _, m := range mx {
		out = m(out)
	}
	return
}

func (w *webserver) registerRoutes(mux *http.ServeMux) *router {
	authHandler := auth.NewHandler(w.provider.GetAuthService())
	userHandler := users.NewHandler(w.provider.GetUserService())
	transactionHandler := transactions.NewHandler(w.provider.GetTransactionService())
	reportHandler := reports.NewHandler(w.provider.GetReportService())

	loggerMiddleware := middlewares.NewLoggerMiddleware()
	authMiddleware := middlewares.NewAuthMiddleware(w.provider)

	routes := NewRouter(mux)
	routes.Use(loggerMiddleware.Execute)

	// non authenticated routes
	routes.Group(func(r *router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		r.Post("/users", userHandler.CreateUser)
		r.Post("/auth/login", authHandler.Login)
	})

	// authenticated routes
	routes.Group(func(r *router) {
		r.Use(authMiddleware.Execute)
		r.Put("/users", userHandler.UpdateUser)
		r.Get("/user", userHandler.GetUser)
		r.Get("/transactions", transactionHandler.ListAllUserTransactions)
		r.Get("/transactions/{id}", transactionHandler.GetTransaction)
		r.Delete("/transactions/{id}", transactionHandler.DeleteTransaction)
		r.Post("/transactions", transactionHandler.CreateTransaction)
		r.Post("/transactions/import", transactionHandler.ImportTransactionsFromFile)
		r.Post("/reports/summarize", reportHandler.Summarize)
		r.Post("/reports/summarize/by-month", reportHandler.SummarizeByMonth)
	})

	return routes
}
