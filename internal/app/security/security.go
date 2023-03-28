package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/exceptions"

	"github.com/google/uuid"
)

func Encrypt(uuid uuid.UUID, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write(uuid[:])
	dst := h.Sum(nil)
	var fullCookie []byte
	fullCookie = append(fullCookie, uuid[:]...)
	fullCookie = append(fullCookie, dst...)
	return hex.EncodeToString(fullCookie)
}

func Decrypt(hashString string, secret []byte) (uuid.UUID, error) {
	var (
		data []byte
		err  error
		sign []byte
	)

	data, err = hex.DecodeString(hashString)
	if err != nil {
		log.Println(err)
		return uuid.UUID{}, exceptions.ErrNotValidSign
	}
	id, idErr := uuid.FromBytes(data[:16])
	if idErr != nil {
		log.Println(idErr)
	}
	h := hmac.New(sha256.New, secret)
	h.Write(data[:16])
	sign = h.Sum(nil)

	if hmac.Equal(sign, data[16:]) {
		return id, nil
	} else {
		return uuid.UUID{}, exceptions.ErrNotValidSign
	}
}
