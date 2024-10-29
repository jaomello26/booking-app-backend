//go:build test

package models

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestGroupValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name     string
		group    Group
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Valid Group",
			group: Group{
				Name:      "Family Group",
				CreatedBy: 1,
			},
			wantErr: false,
		},
		{
			name: "Missing Name",
			group: Group{
				CreatedBy: 1,
			},
			wantErr:  true,
			errorMsg: "Name",
		},
		{
			name: "Missing CreatedBy",
			group: Group{
				Name: "Family Group",
			},
			wantErr:  true,
			errorMsg: "CreatedBy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.group)
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
