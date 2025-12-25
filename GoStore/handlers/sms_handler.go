package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/ramG-reddy/sms-store/models"
	"github.com/ramG-reddy/sms-store/services"
)

// SMSHandler handles HTTP requests for SMS operations
type SMSHandler struct {
	smsService *services.SMSService
}

// NewSMSHandler creates a new SMS handler instance
func NewSMSHandler(smsService *services.SMSService) *SMSHandler {
	return &SMSHandler{
		smsService: smsService,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// GetUserMessages handles GET /v0/user/{user_id}/messages
func (h *SMSHandler) GetUserMessages(w http.ResponseWriter, r *http.Request) {
	// Extract user_id from URL path
	// Expected format: /v0/user/{user_id}/messages
	re := regexp.MustCompile(`^/v0/user/([^/]+)/messages$`)
	matches := re.FindStringSubmatch(r.URL.Path)

	if len(matches) != 2 {
		log.Printf("Invalid URL format: %s", r.URL.Path)
		respondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	userID := matches[1]

	// Validate user_id (phone number format)
	if !isValidPhoneNumber(userID) {
		log.Printf("Invalid user_id format: %s", userID)
		respondWithError(w, http.StatusBadRequest, "Invalid user_id format. Expected phone number.")
		return
	}

	log.Printf("Received request to get messages for user: %s", userID)

	// Retrieve messages from service
	messages, err := h.smsService.GetMessagesByUserID(r.Context(), userID)
	if err != nil {
		log.Printf("Error retrieving messages for user %s: %v", userID, err)
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve messages")
		return
	}

	// Return empty array if no messages found
	if messages == nil {
		messages = make([]*models.SMSRecord, 0)
	}

	log.Printf("Successfully retrieved %d messages for user: %s", len(messages), userID)
	respondWithJSON(w, http.StatusOK, messages)
}

// HealthCheck handles GET /health
func (h *SMSHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":  "UP",
		"service": "sms-store",
	}
	respondWithJSON(w, http.StatusOK, health)
}

// isValidPhoneNumber validates phone number format
// Accepts: +1234567890 or 1234567890 (10-15 digits)
func isValidPhoneNumber(phoneNumber string) bool {
	// Pattern: optional +, digit 1-9, followed by 9-14 more digits
	pattern := `^\+?[1-9]\d{9,14}$`
	matched, err := regexp.MatchString(pattern, phoneNumber)
	if err != nil {
		log.Printf("Error validating phone number: %v", err)
		return false
	}
	return matched
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	respondWithJSON(w, statusCode, errorResponse)
}
