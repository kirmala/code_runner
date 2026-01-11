package service

import (
	"code_processor/http_server/models"
	"code_processor/http_server/repository"
	"code_processor/http_server/usecases/mocks"
	"errors"
	"testing"
)

func TestGetUserId(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		expectedUID string
		expectedErr error
	}{
		{
			name:        "session key is empty",
			key:         "",
			expectedErr: repository.NotFound,
		},
		{
			name:        "session key not found",
			key:         "nonexistent_key",
			expectedErr: repository.NotFound,
		},
		{
			name:        "session key found",
			key:         "valid_key",
			expectedUID: "user123",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			task := NewTask(
				nil,
				&mocks.SessionRepo{
					Sessions: map[string]models.Session{
						"valid_key": {
							UserId:    "user123",
							SessionId: "valid_key",
						},
					},
				},
				nil,
			)

			userId, err := task.GetUserId(tt.key)

			if tt.expectedErr != nil {
				if err == nil || !errors.Is(err, tt.expectedErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if *userId != tt.expectedUID {
				t.Fatalf("expected userId %s, got %s", tt.expectedUID, *userId)
			}
		})
	}
}
