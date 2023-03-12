package controllers

import (
	"log"
	"net/http"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(storage *storage.Storage) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/plain", "application/x-gzip"))
	r.Use(middleware.Compress(5, gzipContentTypes))
	r.Mount("/debug", middleware.Profiler())
	r.Post("/api/user/register", UserRegistration(storage))
	r.Post("/api/user/login", UserAuthentication(storage))
	r.Post("/api/user/orders", PostOrder(storage))
	r.Get("/api/user/orders", GetOrders(storage))
	r.Get("/api/user/balance", GetBalance(storage))
	r.Post("/api/user/balance/withdraw", Withdraw(storage))
	r.Get("/api/user/withdrawals", GetAllWithdrawals(storage))

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
