package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVehicle_Validate(t *testing.T) {
	tests := []struct {
		name    string
		vehicle *Vehicle
		want    error
		wantErr bool
	}{
		{
			name: "Should return nil",
			vehicle: &Vehicle{
				Brand:        "Toyota",
				VehicleModel: "Camry",
				Year:         2022,
				Plate:        "ABC-1234",
				DriverID:     1,
			},
			wantErr: false,
		},
		{
			name: "Should return brand is invalid",
			vehicle: &Vehicle{
				Brand:        "T",
				VehicleModel: "Camry",
				Year:         2022,
				Plate:        "ABC-1234",
				DriverID:     1,
			},
			want:    &ErrorInvalidField{Message: []string{"vehicle brand is invalid"}},
			wantErr: true,
		},
		{
			name: "Should return vehicle model is invalid",
			vehicle: &Vehicle{
				Brand:        "Toyota",
				VehicleModel: "C",
				Year:         2022,
				Plate:        "ABC-1234",
				DriverID:     1,
			},
			want:    &ErrorInvalidField{Message: []string{"vehicle model is invalid"}},
			wantErr: true,
		},
		{
			name: "Should return year is invalid",
			vehicle: &Vehicle{
				Brand:        "Toyota",
				VehicleModel: "Camry",
				Year:         1886,
				Plate:        "ABC-1234",
				DriverID:     1,
			},
			want:    &ErrorInvalidField{Message: []string{"vehicle year is invalid"}},
			wantErr: true,
		},
		{
			name: "Should return plate is invalid",
			vehicle: &Vehicle{
				Brand:        "Toyota",
				VehicleModel: "Camry",
				Year:         2022,
				Plate:        "ABC1234",
				DriverID:     1,
			},
			want:    &ErrorInvalidField{Message: []string{"vehicle plate is invalid"}},
			wantErr: true,
		},
		{
			name: "Should return all fields are invalid",
			vehicle: &Vehicle{
				Brand:        "T",
				VehicleModel: "C",
				Year:         1886,
				Plate:        "ABC1234",
				DriverID:     1,
			},
			want: &ErrorInvalidField{
				Message: []string{
					"vehicle brand is invalid",
					"vehicle model is invalid",
					"vehicle year is invalid",
					"vehicle plate is invalid",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.vehicle.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, tt.want, err)
		})
	}
}
