package controllers

import (
	"github.com/ShishkovEM/amazing-gophermart/internal/app/controllers/handlers"
	"log"
	"net/http"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(storage *storage.Storage, secretKey []byte, tokenLifetime string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/plain", "application/x-gzip"))
	r.Use(middleware.Compress(5, handlers.GzipContentTypes))
	r.Mount("/debug", middleware.Profiler())

	uc := NewUserController(storage, secretKey, tokenLifetime)

	r.Post("/api/user/register", uc.Register)
	r.Post("/api/user/login", uc.Login)
	r.Post("/api/user/orders", uc.PostOrder)
	r.Get("/api/user/orders", uc.GetOrders)
	r.Get("/api/user/balance", uc.GetBalance)
	r.Post("/api/user/balance/withdraw", uc.Withdraw)
	r.Get("/api/user/withdrawals", uc.GetAllWithdrawals)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, nfErr := w.Write([]byte("route does not exist"))
		if nfErr != nil {
			log.Println(nfErr)
		}
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, naErr := w.Write([]byte("sorry, only GET and POST methods are supported."))
		if naErr != nil {
			log.Println(naErr)
		}
	})

	return r
}
