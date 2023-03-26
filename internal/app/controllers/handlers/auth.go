package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/exceptions"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/models"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/storage"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const timeLayout = "2006-01-02 15:04:05"

func UserRegistration(w http.ResponseWriter, r *http.Request, storage *storage.Storage, secretKey []byte, tokenLifetime string) {
	headerContentType := r.Header.Get("Content-Type")
	if !strings.Contains("application/json, application/x-gzip", headerContentType) {
		messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
		return
	}

	var user models.User
	var unmarshalErr *json.UnmarshalTypeError

	b, err := readBodyBytes(r)
	if err != nil {
		messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(b)
	decoder.DisallowUnknownFields()
	errDecode := decoder.Decode(&user)

	if errDecode != nil {
		if errors.As(errDecode, &unmarshalErr) {
			messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
		} else {
			messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
		}
		return
	}

	userID := uuid.New()
	userToken, userTokenExp := GenerateToken(userID, secretKey, tokenLifetime)
	hashedPassword, bcrypteErr := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if bcrypteErr != nil {
		log.Println(bcrypteErr)
	}

	user.ID, user.Password, user.Token, user.TokenExpires = userID, string(hashedPassword), userToken, userTokenExp

	newUserErr := storage.Repo.CreateUser(&user)
	if newUserErr != nil {
		if newUserErr == exceptions.ErrDuplicatePK {
			messageResponse(w, "login is already busy", "application/json", http.StatusConflict)
			return
		}
		messageResponse(w, "Internal Server Error "+newUserErr.Error(), "application/json", http.StatusInternalServerError)
		return
	}

	generatedAt := time.Now().Format(timeLayout)
	expiresAt := userTokenExp.Format(timeLayout)
	tokenDetails := models.Token{
		TokenType:   "Bearer",
		AuthToken:   userToken,
		GeneratedAt: generatedAt,
		ExpiresAt:   expiresAt,
	}
	jsonResp, _ := json.Marshal(tokenDetails)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", tokenDetails.TokenType+" "+tokenDetails.AuthToken)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func UserAuthentication(w http.ResponseWriter, r *http.Request, storage *storage.Storage) {
	headerContentType := r.Header.Get("Content-Type")
	if !strings.Contains("application/json, application/x-gzip", headerContentType) {
		messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
		return
	}

	var user models.User
	var unmarshalErr *json.UnmarshalTypeError

	b, err := readBodyBytes(r)
	if err != nil {
		messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(b)
	decoder.DisallowUnknownFields()
	errDecode := decoder.Decode(&user)

	if errDecode != nil {
		if errors.As(errDecode, &unmarshalErr) {
			messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
		} else {
			messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
		}
		return
	}

	userDB, userDBErr := storage.Repo.ReadUser(user.Username)
	if userDBErr != nil {
		if errors.Is(userDBErr, sql.ErrNoRows) {
			messageResponse(w, "User unauthorized: "+user.Username+", please register at /api/user/register", "application/json", http.StatusUnauthorized)
			return
		}
		messageResponse(w, "Internal Server Error: "+userDBErr.Error(), "application/json", http.StatusInternalServerError)
		return
	}

	cryptErr := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if cryptErr != nil {
		messageResponse(w, "User unauthorized: "+cryptErr.Error(), "application/json", http.StatusUnauthorized)
		return
	}

	if userDB.TokenExpires.Before(time.Now()) {
		log.Println("token expired")
	}

	//// Авторизация по токену
	generatedAt := time.Now().Format(timeLayout)
	expiresAt := userDB.TokenExpires.Format(timeLayout)

	tokenDetails := models.Token{
		TokenType:   "Bearer",
		AuthToken:   userDB.Token,
		GeneratedAt: generatedAt,
		ExpiresAt:   expiresAt,
	}
	jsonResp, _ := json.Marshal(tokenDetails)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", tokenDetails.TokenType+" "+tokenDetails.AuthToken)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
