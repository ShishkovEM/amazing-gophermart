package assembly

import (
	"log"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/config"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/models"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/server"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/service"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/storage"
)

var gopherMartConfig config.Config

func StartApplication() {

	// Считывем конфигурацию приложения
	gopherMartConfig.Parse()

	// Инициализируем хранилище
	appStorage, dbErr := storage.NewStorage(gopherMartConfig.Database)

	// Инициализируем клинет
	client := service.NewProcessingClient(gopherMartConfig.AccrualSystem, "/api/orders")

	// Иницилизируем каналы для обработки заказов
	ordersToProcessingCh := make(chan string)
	ordersProcessedCh := make(chan models.ProcessingOrder)

	// Запускаем процессы обработки заказов
	go service.GetOrdersToProcessing(*appStorage, ordersToProcessingCh)
	go service.GetProcessedInfo(client, ordersToProcessingCh, ordersProcessedCh)
	go service.ApplyLoyalty(*appStorage, ordersProcessedCh)
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	// Проверяем подключение к БД
	ping := appStorage.Repo.Ping()
	log.Println(ping)

	// Запусаем приложение
	MainApp := server.NewServer(&gopherMartConfig, appStorage)
	if runErr := MainApp.Run(); runErr != nil {
		log.Printf("%s", runErr.Error())
	}
}
