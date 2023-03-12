package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/exceptions"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/models"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/storage"
)

func PostOrder(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("text/plain", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		// Проверка авторизации по токену
		userID, tokenErr := GetToken(r)
		if tokenErr != nil {
			messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		body, bodyErr := io.ReadAll(b)
		if bodyErr != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		orderNum, convErr := strconv.Atoi(string(body))
		if convErr != nil {
			messageResponse(w, "invalid order number format", "application/json", http.StatusUnprocessableEntity)
			return
		}

		if !valid(orderNum) {
			messageResponse(w, "invalid order number format", "application/json", http.StatusUnprocessableEntity)
			return
		}

		orderNumStr := fmt.Sprintf("%d", orderNum)
		orderDB, orderDBErr := storage.Repo.CheckOrder(orderNumStr)
		if orderDBErr != nil {
			messageResponse(w, "Internal Server Error: "+orderDBErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		if orderDB.OrderNum == orderNumStr && orderDB.UserID != userID {
			messageResponse(w, "the order number has already been uploaded by another user", "application/json", http.StatusConflict)
			return
		}

		if orderDB.OrderNum == orderNumStr && orderDB.UserID == userID {
			messageResponse(w, "order number has already been uploaded by this user", "application/json", http.StatusOK)
			return
		}

		var order models.Order
		order.UserID, order.OrderNum = userID, orderNumStr

		insertErr := storage.Repo.CreateOrder(&order)

		if insertErr != nil {
			messageResponse(w, "Internal Server Error: "+insertErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		messageResponse(w, "new order number accepted for processing", "application/json", http.StatusAccepted)
	}
}

func GetOrders(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Length")
		if len(headerContentType) != 0 {
			messageResponse(w, "Content-Length is not equal 0", "application/json", http.StatusBadRequest)
			return
		}

		// Проверка авторизации по токену
		userID, tokenErr := GetToken(r)
		if tokenErr != nil {
			messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		orders, ordersErr := storage.Repo.ReadOrders(userID)
		if ordersErr != nil {
			if ordersErr == exceptions.ErrNoValues {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			messageResponse(w, "Internal Server Error: "+ordersErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		ordersList, ordersListErr := json.Marshal(orders)
		if ordersListErr != nil {
			panic(ordersListErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(ordersList)
	}
}
