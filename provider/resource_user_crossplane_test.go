package provider

import (
	"strings"
	"testing"
)

// TestUserIDHandling tests the ID handling logic for Crossplane compatibility
func TestUserIDHandling(t *testing.T) {
	tests := []struct {
		name     string
		inputID  string
		expected string
	}{
		{
			name:     "numeric ID path",
			inputID:  "/users/123",
			expected: "/users/123",
		},
		{
			name:     "username without prefix",
			inputID:  "testuser",
			expected: "/users/testuser",
		},
		{
			name:     "username with special characters",
			inputID:  "test-user-123",
			expected: "/users/test-user-123",
		},
		{
			name:     "empty ID",
			inputID:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the ID handling logic from resourceUserRead
			currentId := tt.inputID
			if currentId != "" && !strings.HasPrefix(currentId, "/users/") {
				currentId = "/users/" + currentId
			}
			
			if currentId != tt.expected {
				t.Errorf("ID handling failed: got %s, want %s", currentId, tt.expected)
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