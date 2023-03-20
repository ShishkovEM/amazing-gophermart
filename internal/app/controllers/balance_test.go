package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/models"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGetBalance(t *testing.T) {
	database, dbErr := storage.NewStorage("postgres://junvlkns:BKHdP45va97hTKwWld-6fg85etq62rP8@trumpet.db.elephantsql.com/junvlkns", migrationsDir)
	var secretKey = []byte("G0pher")

	userID := uuid.New()
	cookie, cookieExpires := GenerateCookie(userID, secretKey, "10h")
	hashedPassword, bcrypteErr := bcrypt.GenerateFromPassword([]byte("123"), 4)
	if bcrypteErr != nil {
		log.Println(bcrypteErr)
	}
	user := models.User{
		ID:            userID,
		Username:      "test",
		Password:      string(hashedPassword),
		Cookie:        cookie.String(),
		CookieExpires: cookieExpires,
	}

	newUserErr := database.Repo.CreateUser(&user)
	if newUserErr != nil {
		log.Println("New User Error", newUserErr)
	}

	if dbErr != nil {
		log.Fatal(dbErr)
	}
	type want struct {
		code            int
		location        string
		contentType     string
		contentEncoding string
		responseFormat  bool
		response        string
	}

	tests := []struct {
		name                   string
		request                string
		requestPath            string
		requestMethod          string
		requestBody            string
		requestCompressBody    []byte
		requestContentType     string
		requestAcceptEncoding  string
		requestContentEncoding string
		requestCookie          string
		requestToken           string
		want                   want
	}{
		{
			name:          fmt.Sprintf("%s test #1", http.MethodGet),
			requestMethod: http.MethodGet,
			requestPath:   "/api/user/balance",
			requestToken:  "",
			want: want{
				code: http.StatusUnauthorized,
			},
		},
	}
	Routes := *Routes(database, secretKey, "10h")
	ts := httptest.NewServer(&Routes)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := []byte(tt.requestBody)
			reqURL := tt.requestPath + tt.request
			request := httptest.NewRequest(tt.requestMethod, reqURL, bytes.NewBuffer(reqBody))
			request.Header.Set("Content-Type", tt.requestContentType)
			request.Header.Set("Cookie", tt.requestCookie)
			if len(tt.requestToken) > 0 {
				request.Header.Set("Authorization", tt.requestToken)
			}
			// создаём новый Recorder
			w := httptest.NewRecorder()
			Routes.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			// Проверяем StatusCode
			respStatusCode := resp.StatusCode
			wantStatusCode := tt.want.code
			assert.Equal(t, wantStatusCode, respStatusCode, fmt.Errorf("expected StatusCode %d, got %d", wantStatusCode, respStatusCode))
		})
	}
}

func TestWithdraws(t *testing.T) {
	database, dbErr := storage.NewStorage("postgres://junvlkns:BKHdP45va97hTKwWld-6fg85etq62rP8@trumpet.db.elephantsql.com/junvlkns", migrationsDir)
	var secretKey = []byte("G0pher")

	userID := uuid.New()
	cookie, cookieExpires := GenerateCookie(userID, secretKey, "10h")
	hashedPassword, bcrypteErr := bcrypt.GenerateFromPassword([]byte("123"), 4)
	if bcrypteErr != nil {
		log.Println(bcrypteErr)
	}
	user := models.User{
		ID:            userID,
		Username:      "test",
		Password:      string(hashedPassword),
		Cookie:        cookie.String(),
		CookieExpires: cookieExpires,
	}
	userToken := fmt.Sprintf("Bearer %s", cookie.Value)

	newUserErr := database.Repo.CreateUser(&user)
	if newUserErr != nil {
		log.Println("New User Error", newUserErr)
	}

	if dbErr != nil {
		log.Fatal(dbErr)
	}
	type want struct {
		code            int
		location        string
		contentType     string
		contentEncoding string
		responseFormat  bool
		response        string
	}

	tests := []struct {
		name                   string
		request                string
		requestPath            string
		requestMethod          string
		requestBody            string
		requestCompressBody    []byte
		requestContentType     string
		requestAcceptEncoding  string
		requestContentEncoding string
		requestCookie          string
		requestToken           string
		want                   want
	}{
		{
			name:          fmt.Sprintf("%s no content #1", http.MethodGet),
			requestMethod: http.MethodGet,
			requestPath:   "/api/user/withdrawals",
			requestCookie: cookie.String(),
			requestToken:  userToken,
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name:          fmt.Sprintf("%s negative Unauthorized #1", http.MethodGet),
			requestMethod: http.MethodGet,
			requestPath:   "/api/user/withdrawals",
			requestToken:  "",
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:               fmt.Sprintf("%s add order to base", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			requestBody:        "12345678903", // новый номер заказа принят в обработку
			requestPath:        "/api/user/orders",
			requestCookie:      cookie.String(),
			requestToken:       userToken,
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name:               fmt.Sprintf("%s nagetive Unauthorized #2", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestPath:        "/api/user/balance/withdraw",
			requestContentType: "application/json",
			requestBody:        `{"order": "2377225624", "sum": 500}`,
			requestToken:       "",
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:               fmt.Sprintf("%s negative withdraw wrong number", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestPath:        "/api/user/balance/withdraw",
			requestContentType: "application/json",
			requestBody:        `{"order": "123", "sum": 500}`,
			requestCookie:      cookie.String(),
			requestToken:       userToken,
			want: want{
				code: http.StatusUnprocessableEntity,
			},
		},
	}
	Routes := *Routes(database, secretKey, "10h")
	ts := httptest.NewServer(&Routes)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := []byte(tt.requestBody)
			reqURL := tt.requestPath + tt.request
			request := httptest.NewRequest(tt.requestMethod, reqURL, bytes.NewBuffer(reqBody))
			request.Header.Set("Content-Type", tt.requestContentType)
			request.Header.Set("Cookie", tt.requestCookie)
			if len(tt.requestToken) > 0 {
				request.Header.Set("Authorization", tt.requestToken)
			}
			// создаём новый Recorder
			w := httptest.NewRecorder()
			Routes.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			// Проверяем StatusCode
			respStatusCode := resp.StatusCode
			wantStatusCode := tt.want.code
			if tt.requestPath == "/api/user/orders" {
				database.Repo.UpdateOrder(models.ProcessingOrder{
					OrderNum: tt.requestBody,
					Status:   "PROCESSED",
					Accrual:  1000.00,
				})
			}
			assert.Equal(t, wantStatusCode, respStatusCode, fmt.Errorf("expected StatusCode %d, got %d", wantStatusCode, respStatusCode))
		})
	}
}
