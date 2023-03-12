package controllers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/exceptions"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/security"

	"github.com/google/uuid"
)

var gzipContentTypes = "application/x-gzip, application/javascript, application/json, text/css, text/html, text/plain, text/xml"

func messageResponse(w http.ResponseWriter, message, ContentType string, httpStatusCode int) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, jsonRespErr := json.Marshal(resp)
	if jsonRespErr != nil {
		log.Println(jsonRespErr)
	}
	w.Write(jsonResp)
}

func readBodyBytes(r *http.Request) (io.ReadCloser, error) {
	// GZIP decode
	if len(r.Header["Content-Encoding"]) > 0 && r.Header["Content-Encoding"][0] == "gzip" {
		// Read body
		bodyBytes, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			return nil, readErr
		}
		defer r.Body.Close()

		newR, gzErr := gzip.NewReader(io.NopCloser(bytes.NewBuffer(bodyBytes)))
		if gzErr != nil {
			log.Println(gzErr)
			return nil, gzErr
		}
		defer newR.Close()

		return newR, nil
	} else {
		return r.Body, nil
	}
}

func GenerateCookie(userID uuid.UUID) (http.Cookie, time.Time) {
	session := security.Encrypt(userID, security.SecretKey)
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "session", Value: session, Expires: expiration, Path: "/"}
	return cookie, expiration
}

func GetToken(r *http.Request) (uuid.UUID, error) {
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return uuid.UUID{}, exceptions.ErrNoAuth
	}
	tokenValue := strings.Split(auth, "Bearer ")
	if len(tokenValue) < 2 {
		return uuid.UUID{}, exceptions.ErrNoAuth
	}
	authToken := tokenValue[1]
	userID, tokenDecryptErr := security.Decrypt(authToken, security.SecretKey)
	if tokenDecryptErr != nil {
		return uuid.UUID{}, tokenDecryptErr
	}
	return userID, nil
}

func ParseCookie(cookieStr string) (string, error) {
	cookieInfo := strings.Split(cookieStr, "; ")
	for _, pairs := range cookieInfo {
		elements := strings.Split(pairs, "=")
		if elements[0] == "session" {
			return elements[1], nil
		}
	}
	return "", exceptions.ErrNoCookie
}

func valid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur *= 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number /= 10
	}
	return luhn % 10
}
