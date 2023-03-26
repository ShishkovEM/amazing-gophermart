package controllers

import (
	"github.com/ShishkovEM/amazing-gophermart/internal/app/controllers/handlers"
	"net/http"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/storage"
)

type UserController struct {
	storage       *storage.Storage
	secretKey     []byte
	tokenLifetime string
}

func NewUserController(storage *storage.Storage, secretKey []byte, tokenLifetime string) *UserController {
	return &UserController{storage, secretKey, tokenLifetime}
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
	handlers.UserRegistration(w, r, uc.storage, uc.secretKey, uc.tokenLifetime)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	handlers.UserAuthentication(w, r, uc.storage)
}

func (uc *UserController) PostOrder(w http.ResponseWriter, r *http.Request) {
	handlers.PostOrder(w, r, uc.storage, uc.secretKey)
}

func (uc *UserController) GetOrders(w http.ResponseWriter, r *http.Request) {
	handlers.GetOrders(w, r, uc.storage, uc.secretKey)
}

func (uc *UserController) GetBalance(w http.ResponseWriter, r *http.Request) {
	handlers.GetBalance(w, r, uc.storage, uc.secretKey)
}

func (uc *UserController) Withdraw(w http.ResponseWriter, r *http.Request) {
	handlers.Withdraw(w, r, uc.storage, uc.secretKey)
}

func (uc *UserController) GetAllWithdrawals(w http.ResponseWriter, r *http.Request) {
	handlers.GetAllWithdrawals(w, r, uc.storage, uc.secretKey)
}
