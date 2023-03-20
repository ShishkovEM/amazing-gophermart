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

func TestPostOrder(t *testing.T) {
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

	subUserID := uuid.New()
	subCookie, subCookieExpires := GenerateCookie(subUserID, secretKey, "10h")
	subHashedPassword, subBcrypteErr := bcrypt.GenerateFromPassword([]byte("123"), 4)
	if subBcrypteErr != nil {
		log.Println(subBcrypteErr)
	}
	subUser := models.User{
		ID:            subUserID,
		Username:      "test2",
		Password:      string(subHashedPassword),
		Cookie:        subCookie.String(),
		CookieExpires: subCookieExpires,
	}

	database.Repo.CreateUser(&subUser)

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
			name:               fmt.Sprintf("%s negative #1", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			requestBody:        "12345-678903", // неверный формат номера заказа
			requestPath:        "/api/user/orders",
			requestCookie:      cookie.String(),
			requestToken:       userToken,
			want: want{
				code: http.StatusUnprocessableEntity,
			},
		},
		{
			name:               fmt.Sprintf("%s negative #2", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json", // неверный Content-Type неверный формат запроса
			requestBody:        "12345-678903",
			requestPath:        "/api/user/orders",
			requestCookie:      cookie.String(),
			requestToken:       userToken,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative #3", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			requestBody:        "1234567890", // не проходит проверку по Луну
			requestPath:        "/api/user/orders",
			requestCookie:      cookie.String(),
			requestToken:       userToken,
			want: want{
				code: http.StatusUnprocessableEntity,
			},
		},
		{
			name:               fmt.Sprintf("%s negative #4", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			requestBody:        "1234567890", // без cookie
			requestPath:        "/api/user/orders",
			requestCookie:      "",
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
			assert.Equal(t, respStatusCode, wantStatusCode, fmt.Errorf("expected StatusCode %d, got %d", wantStatusCode, respStatusCode))
		})
	}
}

func TestGetOrders(t *testing.T) {
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

	subUserID := uuid.New()
	subCookie, subCookieExpires := GenerateCookie(subUserID, secretKey, "10h")
	subHashedPassword, subBcrypteErr := bcrypt.GenerateFromPassword([]byte("123"), 4)
	if subBcrypteErr != nil {
		log.Println(subBcrypteErr)
	}
	subUser := models.User{
		ID:            subUserID,
		Username:      "test2",
		Password:      string(subHashedPassword),
		Cookie:        subCookie.String(),
		CookieExpires: subCookieExpires,
	}
	subUserToken := fmt.Sprintf("Bearer %s", subCookie.Value)
	database.Repo.CreateUser(&subUser)

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
			name:               fmt.Sprintf("%s post test order #1", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			requestBody:        "1230", // новый номер заказа принят в обработку
			requestPath:        "/api/user/orders",
			requestCookie:      cookie.String(),
			requestToken:       userToken,
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name:          fmt.Sprintf("%s positive #1", http.MethodGet),
			requestMethod: http.MethodGet,
			requestPath:   "/api/user/orders", // нет заказов в базе
			requestCookie: subCookie.String(),
			requestToken:  subUserToken,
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name:          fmt.Sprintf("%s negative #1", http.MethodGet),
			requestMethod: http.MethodGet,
			requestPath:   "/api/user/orders",
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
			headerContentType := resp.Header.Get("Content-Length")
			log.Println("LENGTH: ", headerContentType)
			assert.Equal(t, respStatusCode, wantStatusCode, fmt.Errorf("expected StatusCode %d, got %d", wantStatusCode, respStatusCode))
		})
	}
}
