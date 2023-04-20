package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EncryptID(id string) string {
	hash := sha256.Sum256([]byte(id))
	return hex.EncodeToString(hash[:])
}

// ValidateID memvalidasi sebuah id yang telah dienkripsi dengan SHA256
func ValidateID(id string, encryptedID string) bool {
	return encryptedID == EncryptID(id)
}

func ReadRequest(ctx *gin.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(ctx.Request.Context(), request)
}

func SuccessResponse(data interface{}) Result {
	return Result{
		Data: data,
		Response: Response{
			Code:    http.StatusOK,
			Message: "Ok",
			Status:  "success",
		},
	}
}

func ErrorResponse(code int, message string) Result {
	return Result{
		Data: nil,
		Response: Response{
			Code:    code,
			Message: message,
			Status:  "error",
		},
	}
}

func GenerateRandomID() string {
	// Generate a random byte array
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the byte array to a string
	id := hex.EncodeToString(bytes)

	return id
}
