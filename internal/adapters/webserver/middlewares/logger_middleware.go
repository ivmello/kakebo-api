package middlewares

import (
	"fmt"
	"net/http"
)

type loggerMiddleware struct{}

func NewLoggerMiddleware() *loggerMiddleware {
	return &loggerMiddleware{}
}

func (m *loggerMiddleware) Execute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method
		fmt.Println(method, path)
		next.ServeHTTP(w, r)
	})
}
