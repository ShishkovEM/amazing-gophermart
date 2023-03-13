package repository

import (
	"github.com/ShishkovEM/amazing-gophermart/internal/app/models"

	"github.com/google/uuid"
)

type CommonRepository interface {
	Ping() bool
	CreateUser(user *models.User) error
	ReadUser(username string) (*models.User, error)
	CheckOrder(orderNum string) (*models.Order, error)
	CreateOrder(order *models.Order) error
	ReadOrders(userID uuid.UUID) ([]*models.OrderDB, error)
	ReadBalance(userID uuid.UUID) (*models.Balance, error)
	CreateWithdrawal(withdraw *models.Withdraw) error
	ReadAllWithdrawals(userID uuid.UUID) ([]*models.WithdrawDB, error)
	ReadOrdersForProcessing() ([]string, error)
	UpdateOrder(order models.ProcessingOrder)
}
