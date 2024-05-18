package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateSecret(originalSecret string, options ...int64) string {
	var window int64 = 1800 // Default window size (30 minutes)
	if len(options) > 0 {
		window = options[0] // Override default window size if provided
	}

	// Get the current time 30 minutes window
	currentTime := time.Now().Unix() / window

	// Convert the current time to a string
	timeString := fmt.Sprintf("%d", currentTime)

	// Create a new HMAC using SHA256
	h := hmac.New(sha256.New, []byte(originalSecret))

	// Write the time string to the HMAC
	h.Write([]byte(timeString))

	// Get the resulting HMAC
	secret := h.Sum(nil)

	// Convert the HMAC to a hexadecimal string
	return hex.EncodeToString(secret)
}

func VerifySecret(originalSecret, secret string) bool {
	// Get the current time 30 minutes window
	currentTime := time.Now().Unix() / (30 * 60)

	// Convert the current time to a string
	timeString := fmt.Sprintf("%d", currentTime)

	// Create a new HMAC using SHA256
	h := hmac.New(sha256.New, []byte(originalSecret))

	// Write the time string to the HMAC
	h.Write([]byte(timeString))

	// Get the resulting HMAC
	expectedSecret := h.Sum(nil)

	// Convert the HMAC to a hexadecimal string
	expectedSecretString := hex.EncodeToString(expectedSecret)

	// Compare the expected secret with the provided secret
	return secret == expectedSecretString
}

func GenerateJWT(secret string,
	claims jwt.MapClaims,
	exp time.Duration,
) string {

	claims["exp"] = time.Now().Add(exp).Unix()

	log.Println("Time expiring in hours and minutes local time: ", time.Now().Add(exp).Local())
	log.Println("Time duration in seconds: ", exp.Seconds())

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Errorf("Error signing token: ", err)
		return ""
	}

	return signedToken
}

func VerifyJWT(secret, token string) (jwt.MapClaims, error) {
	// Parse the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Get the claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func EncryptPassword(password string) string {
	// Create a new HMAC using SHA256
	h := hmac.New(sha256.New, []byte(password))

	// Get the resulting HMAC
	encryptedPassword := h.Sum(nil)

	// Convert the HMAC to a hexadecimal string
	return hex.EncodeToString(encryptedPassword)
}

func GenerateRandomString(options ...int) string {
	var length int = 36
	if len(options) > 0 {
		length = options[0]
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		// Generate a random index within the charset length
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[idx.Int64()]
	}
	return string(b)
}
