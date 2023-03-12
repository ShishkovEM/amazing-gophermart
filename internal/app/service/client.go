package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/logger"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/models"

	"gopkg.in/eapache/go-resiliency.v1/retrier"
	"gopkg.in/h2non/gentleman-retry.v2"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
)

var ErrInternalServer = errors.New("ErrInternalServer")
var ErrEmptyOrder = errors.New("empty order")

type ProcessingClient struct {
	Client *gentleman.Client
}

func NewProcessingClient(serviceAddress, basicURL string) *ProcessingClient {
	log.Println("LoyalityServer: ", serviceAddress+basicURL)
	cli := gentleman.New()
	cli.Use(logger.New(os.Stdout))
	cli.Use(timeout.Request(60 * time.Second))
	cli.Use(retry.New(retrier.New(retrier.ExponentialBackoff(5, 100*time.Millisecond), nil)))
	cli.URL(serviceAddress + basicURL)
	return &ProcessingClient{
		Client: cli,
	}
}

func (pc *ProcessingClient) GetOrder(orderNum string) (models.ProcessingOrder, error) {
	req := pc.Client.Request()
	req.Method("GET")
	req.AddPath(fmt.Sprintf("/%s", orderNum))
	res, err := req.Send()
	var order models.ProcessingOrder
	if err != nil {
		return order, err
	}

	switch res.StatusCode {
	case http.StatusInternalServerError:
		log.Printf("Internal server error: %d\n", res.StatusCode)
		return order, ErrInternalServer
	case http.StatusTooManyRequests:
		log.Printf("Too Many Requests: %d\n", res.StatusCode)
		time.Sleep(time.Second * 60)
	case http.StatusOK:
		if UnmarshErr := json.Unmarshal(res.Bytes(), &order); UnmarshErr != nil {
			return order, UnmarshErr
		}
	}

	emptyOrder := models.ProcessingOrder{
		OrderNum: "",
		Status:   "",
		Accrual:  nil,
	}

	if order == emptyOrder {
		return order, ErrEmptyOrder
	}

	if order.OrderNum == "" {
		return order, ErrEmptyOrder
	}

	return order, nil
}
