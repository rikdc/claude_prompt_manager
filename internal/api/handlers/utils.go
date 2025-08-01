package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/claude-code-template/prompt-manager/internal/database"
)

// GetOrCreateConversation finds an existing conversation by session ID or creates a new one.
// It uses a direct database query to efficiently find conversations by session ID.
// If no match is found, it creates a new conversation with optional context data.
func GetOrCreateConversation(db *database.DB, sessionID string, data map[string]interface{}) (int, error) {
	// Try to find existing conversation for this session using efficient lookup
	conv, err := db.GetConversationBySessionID(sessionID)
	if err == nil {
		// Found existing conversation
		return conv.ID, nil
	}

	// Check if error is "not found" - if so, create new conversation
	// For other errors, return them
	if err.Error() != "conversation not found" {
		return 0, fmt.Errorf("failed to lookup conversation by session ID: %w", err)
	}

	// Create new conversation
	workingDir := ExtractStringFromData(data, "cwd")
	transcriptPath := ExtractStringFromData(data, "transcript_path")

	newConv, err := db.CreateConversation(sessionID, nil, workingDir, transcriptPath)
	if err != nil {
		return 0, fmt.Errorf("failed to create conversation: %w", err)
	}

	return newConv.ID, nil
}

// ExtractStringFromData safely extracts a string value from map data.
// Returns a pointer to the string if the key exists and the value is a non-empty string,
// otherwise returns nil.
func ExtractStringFromData(data map[string]interface{}, key string) *string {
	if value, exists := data[key]; exists {
		if str, ok := value.(string); ok && str != "" {
			return &str
		}
	}
	return nil
}

// ErrorResponse sends a standardized error response in JSON format.
// It sets the appropriate content type, status code, and response structure
// consistent across all handlers.
func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: false,
		Error:   &message,
	}
	json.NewEncoder(w).Encode(response)
}
