//go:build test

package models

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestUserValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name     string
		user     User
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Valid User",
			user: User{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Invalid Email",
			user: User{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr:  true,
			errorMsg: "Email",
		},
		{
			name: "Short Name",
			user: User{
				Name:     "Jo",
				Email:    "john@example.com",
				Password: "password123",
			},
			wantErr:  true,
			errorMsg: "Name",
		},
		{
			name: "Short Password",
			user: User{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "pass",
			},
			wantErr:  true,
			errorMsg: "Password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.user)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}
			if tt.wantErr && err != nil && tt.errorMsg != "" {
				if !containsValidationError(err, tt.errorMsg) {
					t.Errorf("Expected error containing '%s' but got: %v", tt.errorMsg, err)
				}
			}
		})
	}
}

func containsValidationError(err error, field string) bool {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			if e.Field() == field {
				return true
			}
		}
	}
	return false
}
