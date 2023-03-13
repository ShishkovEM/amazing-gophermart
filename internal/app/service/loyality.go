package service

import (
	"time"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/models"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/storage"
)

func GetOrdersToProcessing(storage storage.Storage, ordersCh chan string) {
	for {
		loyal, _ := storage.Repo.ReadOrdersForProcessing()
		if len(loyal) > 0 {
			for _, order := range loyal {
				ordersCh <- order
			}
		} else {
			time.Sleep(time.Second * 10)
		}
		time.Sleep(time.Second * 2)
	}
}

func GetProcessedInfo(client *ProcessingClient, ordersCh chan string, procesedCh chan models.ProcessingOrder) {
	for order := range ordersCh {
		if orderInfo, orderInfoErr := client.GetOrder(order); orderInfoErr == nil {
			procesedCh <- orderInfo
		}
	}
}

func ApplyLoyalty(storage storage.Storage, procesedCh chan models.ProcessingOrder) {
	for order := range procesedCh {
		storage.Repo.UpdateOrder(order)
	}
}
