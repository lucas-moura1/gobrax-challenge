package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDriver_Validate(t *testing.T) {
	tests := []struct {
		name    string
		driver  *Driver
		want    error
		wantErr bool
	}{
		{
			name: "Should return nil",
			driver: &Driver{
				Name:        "John",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Phone:       "1234567890",
				License:     "ABC123",
				LicenseType: LicenseTypeA,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Should return name is invalid",
			driver: &Driver{
				Name:        "L",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Phone:       "1234567890",
				License:     "ABC123",
				LicenseType: LicenseTypeA,
			},
			want:    &ErrorInvalidField{Message: []string{"driver name is invalid"}},
			wantErr: true,
		},
		{
			name: "Should return last name is invalid",
			driver: &Driver{
				Name:        "John",
				LastName:    "D",
				Email:       "john@test.com",
				Phone:       "1234567890",
				License:     "ABC123",
				LicenseType: LicenseTypeA,
			},
			want:    &ErrorInvalidField{Message: []string{"driver last name is invalid"}},
			wantErr: true,
		},
		{
			name: "Should return email is invalid",
			driver: &Driver{
				Name:        "John",
				LastName:    "Doe",
				Email:       "john.doe.com",
				Phone:       "1234567890",
				License:     "ABC123",
				LicenseType: LicenseTypeA,
			},
			want:    &ErrorInvalidField{Message: []string{"driver email is invalid"}},
			wantErr: true,
		},
		{
			name: "Should return phone is invalid",
			driver: &Driver{
				Name:        "John",
				LastName:    "Doe",
				Email:       "john@test.com",
				Phone:       "123456789",
				License:     "ABC123",
				LicenseType: LicenseTypeA,
			},
			want:    &ErrorInvalidField{Message: []string{"driver phone is invalid"}},
			wantErr: true,
		},
		{
			name: "Should return license is invalid",
			driver: &Driver{
				Name:        "John",
				LastName:    "Doe",
				Email:       "john@test.com",
				Phone:       "1234567890",
				License:     "ABC",
				LicenseType: LicenseTypeA,
			},
			want:    &ErrorInvalidField{Message: []string{"driver license is invalid"}},
			wantErr: true,
		},
		{
			name: "Should return license type is invalid",
			driver: &Driver{
				Name:        "John",
				LastName:    "Doe",
				Email:       "john@test.com",
				Phone:       "1234567890",
				License:     "ABC123",
				LicenseType: "X",
			},
			want:    &ErrorInvalidField{Message: []string{"driver license type is invalid"}},
			wantErr: true,
		},
		{
			name: "Should return all fields are invalid",
			driver: &Driver{
				Name:        "L",
				LastName:    "D",
				Email:       "john.doe.com",
				Phone:       "123456789",
				License:     "ABC",
				LicenseType: "X",
			},
			want: &ErrorInvalidField{
				Message: []string{
					"driver name is invalid",
					"driver last name is invalid",
					"driver email is invalid",
					"driver phone is invalid",
					"driver license is invalid",
					"driver license type is invalid",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.driver.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, tt.want, err)
		})
	}
}
