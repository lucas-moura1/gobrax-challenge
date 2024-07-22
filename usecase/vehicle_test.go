package usecase

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lucas-moura1/gobrax-challenge/entity"
	"github.com/lucas-moura1/gobrax-challenge/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_vehicleUsecase_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mockVehicleRepo *repository.MockVehicleRepository)
		want    []*entity.Vehicle
		wantErr bool
	}{
		{
			name: "Should return all vehicles",
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().GetAll().Return([]*entity.Vehicle{
					{
						Brand:        "Toyota",
						VehicleModel: "Camry",
						Year:         2022,
						Plate:        "ABC-1234",
						DriverID:     1,
					},
				}, nil)
			},
			want: []*entity.Vehicle{
				{
					Brand:        "Toyota",
					VehicleModel: "Camry",
					Year:         2022,
					Plate:        "ABC-1234",
					DriverID:     1,
				},
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().GetAll().Return(nil, fmt.Errorf("some error occurred"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockVehicleRepo := repository.NewMockVehicleRepository(ctrl)

			tt.setup(mockVehicleRepo)

			vu := NewVehicleUsecase(mockVehicleRepo)
			got, err := vu.GetAll()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_vehicleUsecase_GetById(t *testing.T) {
	tests := []struct {
		name      string
		vehicleId int
		setup     func(mockVehicleRepo *repository.MockVehicleRepository)
		want      *entity.Vehicle
		wantErr   bool
	}{
		{
			name:      "Should return vehicle by ID",
			vehicleId: 1,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().GetById(1).Return(&entity.Vehicle{
					Brand:        "Toyota",
					VehicleModel: "Camry",
					Year:         2022,
					Plate:        "ABC-1234",
					DriverID:     1,
				}, nil)
			},
			want: &entity.Vehicle{
				Brand:        "Toyota",
				VehicleModel: "Camry",
				Year:         2022,
				Plate:        "ABC-1234",
				DriverID:     1,
			},
			wantErr: false,
		},
		{
			name:      "Should return error for invalid vehicle ID",
			vehicleId: 0,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:      "Should return error when vehicle is not found",
			vehicleId: 2,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().GetById(2).Return(nil, nil)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:      "Should return error",
			vehicleId: 3,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().GetById(3).Return(nil, fmt.Errorf("some error occurred"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockVehicleRepo := repository.NewMockVehicleRepository(ctrl)

			tt.setup(mockVehicleRepo)

			vu := NewVehicleUsecase(mockVehicleRepo)
			got, err := vu.GetById(tt.vehicleId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_vehicleUsecase_Update(t *testing.T) {
	mockUpdatedVehicle := &entity.Vehicle{
		Brand:        "Toyota",
		VehicleModel: "Camry",
		Year:         2023,
		Plate:        "DEF-5678",
	}
	tests := []struct {
		name          string
		vehicleId     int
		updateVehicle *entity.Vehicle
		setup         func(mockVehicleRepo *repository.MockVehicleRepository)
		wantErr       bool
	}{
		{
			name:          "Should update vehicle successfully",
			vehicleId:     1,
			updateVehicle: mockUpdatedVehicle,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().GetById(1).Return(&entity.Vehicle{
					Brand:        "Toyotta",
					VehicleModel: "Canry",
					Year:         2022,
					Plate:        "ABC-1234",
				}, nil)
				mockVehicleRepo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:          "Should return error for invalid vehicle ID",
			vehicleId:     0,
			updateVehicle: mockUpdatedVehicle,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
			},
			wantErr: true,
		},
		{
			name:          "Should return error to get vehicle by ID",
			vehicleId:     1,
			updateVehicle: mockUpdatedVehicle,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().GetById(1).Return(nil, fmt.Errorf("some error occurred"))
			},
			wantErr: true,
		},
		{
			name:          "Should return error when vehicle is not found",
			vehicleId:     2,
			updateVehicle: mockUpdatedVehicle,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().GetById(2).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:      "Should return error for some invalid field",
			vehicleId: 1,
			updateVehicle: &entity.Vehicle{
				Brand:        "T",
				VehicleModel: "Camry",
				Year:         2023,
				Plate:        "DEF-5678",
			},
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().GetById(1).Return(&entity.Vehicle{
					Brand:        "Toyotta",
					VehicleModel: "Canry",
					Year:         2022,
					Plate:        "ABC-1234",
				}, nil)
			},
			wantErr: true,
		},
		{
			name:          "Should return error",
			vehicleId:     3,
			updateVehicle: mockUpdatedVehicle,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().GetById(3).Return(new(entity.Vehicle), nil)
				mockVehicleRepo.EXPECT().Update(gomock.Any()).Return(errors.New("some error occurred"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockVehicleRepo := repository.NewMockVehicleRepository(ctrl)

			tt.setup(mockVehicleRepo)

			vu := NewVehicleUsecase(mockVehicleRepo)
			err := vu.Update(tt.vehicleId, tt.updateVehicle)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_vehicleUsecase_Delete(t *testing.T) {
	tests := []struct {
		name      string
		vehicleId int
		setup     func(mockVehicleRepo *repository.MockVehicleRepository)
		wantErr   bool
	}{
		{
			name:      "Should delete vehicle",
			vehicleId: 1,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().Delete(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Should return error for invalid vehicle ID",
			vehicleId: 0,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
			},
			wantErr: true,
		},
		{
			name:      "Should return error",
			vehicleId: 3,
			setup: func(mockVehicleRepo *repository.MockVehicleRepository) {
				mockVehicleRepo.EXPECT().Delete(3).Return(errors.New("some error occurred"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockVehicleRepo := repository.NewMockVehicleRepository(ctrl)

			tt.setup(mockVehicleRepo)

			vu := NewVehicleUsecase(mockVehicleRepo)
			err := vu.Delete(tt.vehicleId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
