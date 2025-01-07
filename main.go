package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pquerna/otp/totp"
)

// API response structure
type ApiResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Generate TOTP key
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

	// Respond with key URL
	response := map[string]string{
		"key": key.Secret(),
		"url": key.URL(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Validate TOTP code
func validateTOTP(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OTP string `json:"otp"`
	}
	// Decode the incoming JSON request
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Hardcoded secret key (this should be dynamic based on user in production)
	secret := "JBSWY3DPEHPK3PXP" // You can get this from `/generate` API

	// Validate the OTP
	valid := totp.Validate(input.OTP, secret)
	var response ApiResponse

	if valid {
		response = ApiResponse{Message: "OTP is valid!", Success: true}
	} else {
		response = ApiResponse{Message: "Invalid OTP", Success: false}
	}

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// API routes
	http.HandleFunc("/generate", generateTOTP)
	http.HandleFunc("/validate", validateTOTP)

	// Start the server
	port := ":8000"
	fmt.Println("Starting server on", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
