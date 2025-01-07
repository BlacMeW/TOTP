package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// GenerateResponse defines the structure of the generate API response
type GenerateResponse struct {
	Key    string `json:"key"`
	URL    string `json:"url"`
	QRCode string `json:"qrCode"`
}

// ValidateRequest defines the structure of the validate API request
type ValidateRequest struct {
	OTP    string `json:"otp"`
	Secret string `json:"secret"`
}

// ValidateResponse defines the structure of the validate API response
type ValidateResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Generate TOTP key and QR Code
func generateTOTP(w http.ResponseWriter, r *http.Request) {
	// Generate a new TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MyApp",
		AccountName: "user@domain.com",
	})
	if err != nil {
		http.Error(w, "Failed to generate TOTP key", http.StatusInternalServerError)
		return
	}

	// Generate QR Code
	qrCode, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Failed to generate QR Code", http.StatusInternalServerError)
		return
	}

	// Encode QR Code to Base64
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

	// Prepare the response
	response := GenerateResponse{
		Key:    key.Secret(),
		URL:    key.URL(),
		QRCode: fmt.Sprintf("data:image/png;base64,%s", qrCodeBase64),
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Validate TOTP code
func validateTOTP(w http.ResponseWriter, r *http.Request) {
	var input ValidateRequest

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate OTP
	valid := totp.Validate(input.OTP, input.Secret)

	// Prepare the response
	response := ValidateResponse{}
	if valid {
		response.Message = "OTP is valid!"
		response.Success = true
	} else {
		response.Message = "Invalid OTP!"
		response.Success = false
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// API routes
	http.HandleFunc("/generate", generateTOTP)
	http.HandleFunc("/validate", validateTOTP)

	// Start the server
	port := ":8000"
	fmt.Println("Server running on", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
