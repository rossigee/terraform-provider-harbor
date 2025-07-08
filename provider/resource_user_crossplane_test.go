package provider

import (
	"regexp"
	"strings"
	"testing"
)

// TestUserIDDetection tests the numeric ID detection logic
func TestUserIDDetection(t *testing.T) {
	tests := []struct {
		name        string
		inputID     string
		isNumericId bool
	}{
		{
			name:        "numeric ID path",
			inputID:     "/users/123",
			isNumericId: true,
		},
		{
			name:        "numeric ID with large number",
			inputID:     "/users/999999",
			isNumericId: true,
		},
		{
			name:        "username without prefix",
			inputID:     "testuser",
			isNumericId: false,
		},
		{
			name:        "username with prefix",
			inputID:     "/users/testuser",
			isNumericId: false,
		},
		{
			name:        "username with special characters",
			inputID:     "/users/test-user-123",
			isNumericId: false,
		},
		{
			name:        "empty ID",
			inputID:     "",
			isNumericId: false,
		},
		{
			name:        "malformed path",
			inputID:     "/users/",
			isNumericId: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the regex pattern used in resourceUserRead
			isNumericId := regexp.MustCompile(`^/users/\d+$`).MatchString(tt.inputID)
			
			if isNumericId != tt.isNumericId {
				t.Errorf("ID detection failed for %s: got %v, want %v", tt.inputID, isNumericId, tt.isNumericId)
			}
		})
	}
}

// TestUsernameExtraction tests username extraction from various ID formats
func TestUsernameExtraction(t *testing.T) {
	tests := []struct {
		name     string
		inputID  string
		expected string
	}{
		{
			name:     "extract from path",
			inputID:  "/users/testuser",
			expected: "testuser",
		},
		{
			name:     "already username",
			inputID:  "testuser",
			expected: "testuser",
		},
		{
			name:     "numeric ID should remain",
			inputID:  "/users/123",
			expected: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the username extraction logic
			searchUsername := tt.inputID
			if strings.HasPrefix(searchUsername, "/users/") {
				searchUsername = strings.TrimPrefix(searchUsername, "/users/")
			}
			
			if searchUsername != tt.expected {
				t.Errorf("Username extraction failed: got %s, want %s", searchUsername, tt.expected)
			}
		})
	}
}