package handlers

import (
	"bytes"
	"fmt"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/controllers"
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

var migrationsDir = "file://../../../schema"

func TestUserRegistration(t *testing.T) {
	database, dbErr := storage.NewStorage("postgres://junvlkns:BKHdP45va97hTKwWld-6fg85etq62rP8@trumpet.db.elephantsql.com/junvlkns", migrationsDir)
	var secretKey = []byte("G0pher")

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
		want                   want
	}{
		{
			name:               fmt.Sprintf("%s positive test #1", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": "123"}`,
			requestPath:        "/api/user/register",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "passord": "123"}`,
			requestPath:        "/api/user/register",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #2", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": "123", "asdqs" : 123}`,
			requestPath:        "/api/user/register",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #3", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": 123}`,
			requestPath:        "/api/user/register",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	Routes := *controllers.Routes(database, secretKey, "10h")
	ts := httptest.NewServer(&Routes)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := []byte(tt.requestBody)
			reqURL := tt.requestPath + tt.request
			request := httptest.NewRequest(tt.requestMethod, reqURL, bytes.NewBuffer(reqBody))
			request.Header.Set("Content-Type", tt.requestContentType)
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

func TestUserAuthentication(t *testing.T) {
	database, dbErr := storage.NewStorage("postgres://junvlkns:BKHdP45va97hTKwWld-6fg85etq62rP8@trumpet.db.elephantsql.com/junvlkns", migrationsDir)
	var secretKey = []byte("G0pher")

	userID := uuid.New()
	token, tokenExpires := GenerateToken(userID, secretKey, "10h")
	hashedPassword, bcrypteErr := bcrypt.GenerateFromPassword([]byte("123"), 4)
	if bcrypteErr != nil {
		log.Println(bcrypteErr)
	}
	user := models.User{
		ID:           userID,
		Username:     "test",
		Password:     string(hashedPassword),
		Token:        token,
		TokenExpires: tokenExpires,
	}

	database.Repo.CreateUser(&user)

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
		want                   want
	}{
		{
			name:               fmt.Sprintf("%s positive test #1", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": "123"}`,
			requestPath:        "/api/user/login",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "passord": "123"}`,
			requestPath:        "/api/user/login",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #2", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": "123", "asdqs" : 123}`,
			requestPath:        "/api/user/login",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #3", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": 123}`,
			requestPath:        "/api/user/login",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	Routes := *controllers.Routes(database, secretKey, "10h")
	ts := httptest.NewServer(&Routes)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := []byte(tt.requestBody)
			reqURL := tt.requestPath + tt.request
			request := httptest.NewRequest(tt.requestMethod, reqURL, bytes.NewBuffer(reqBody))
			request.Header.Set("Content-Type", tt.requestContentType)
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
