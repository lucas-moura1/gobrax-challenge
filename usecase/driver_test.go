package usecase

import (
	"fmt"
	"testing"

	"github.com/lucas-moura1/gobrax-challenge/entity"
	"github.com/lucas-moura1/gobrax-challenge/repository"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func Test_driveUsecase_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mockDriveRepo *repository.MockDriverRepository)
		want    []*entity.Driver
		wantErr bool
	}{
		{
			name: "Should return all drives",
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetAll().Return([]*entity.Driver{
					{
						Name:        "Lucas",
						LastName:    "Moura",
						Email:       "lucas@test.com",
						Phone:       "123456789",
						License:     "123456",
						LicenseType: "A",
					},
				}, nil)
			},
			want: []*entity.Driver{
				{
					Name:        "Lucas",
					LastName:    "Moura",
					Email:       "lucas@test.com",
					Phone:       "123456789",
					License:     "123456",
					LicenseType: "A",
				},
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetAll().Return(nil, fmt.Errorf("some error occurred"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriveRepo := repository.NewMockDriverRepository(ctrl)

			tt.setup(mockDriveRepo)

			vu := NewDriverUsecase(zap.NewNop().Sugar(), mockDriveRepo)
			got, err := vu.GetAll()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_driveUsecase_GetById(t *testing.T) {
	tests := []struct {
		name           string
		driverId       int
		includeVehicle bool
		setup          func(mockDriveRepo *repository.MockDriverRepository)
		want           *entity.Driver
		wantErr        bool
	}{
		{
			name:           "Should return driver without vehicle",
			driverId:       1,
			includeVehicle: false,
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(1, false).Return(&entity.Driver{
					Name:        "Lucas",
					LastName:    "Moura",
					Email:       "lucas@test.com",
					Phone:       "123456789",
					License:     "123456",
					LicenseType: "A",
				}, nil)
			},
			want: &entity.Driver{
				Name:        "Lucas",
				LastName:    "Moura",
				Email:       "lucas@test.com",
				Phone:       "123456789",
				License:     "123456",
				LicenseType: "A",
			},
			wantErr: false,
		},
		{
			name:           "Should return driver with vehicle",
			driverId:       2,
			includeVehicle: true,
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(2, true).Return(&entity.Driver{
					Name:        "John",
					LastName:    "Doe",
					Email:       "john@test.com",
					Phone:       "987654321",
					License:     "654321",
					LicenseType: "B",
					Vehicles: []entity.Vehicle{
						{
							Brand:        "Toyota",
							VehicleModel: "Corolla",
							Year:         2021,
							Plate:        "XYZ-9876",
						},
					},
				}, nil)
			},
			want: &entity.Driver{
				Name:        "John",
				LastName:    "Doe",
				Email:       "john@test.com",
				Phone:       "987654321",
				License:     "654321",
				LicenseType: "B",
				Vehicles: []entity.Vehicle{
					{
						Brand:        "Toyota",
						VehicleModel: "Corolla",
						Year:         2021,
						Plate:        "XYZ-9876",
					},
				},
			},
			wantErr: false,
		},
		{
			name:           "Should return error for invalid driver ID",
			driverId:       -1,
			includeVehicle: false,
			setup:          func(mockDriveRepo *repository.MockDriverRepository) {},
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "Should return nil for driver not found",
			driverId:       3,
			includeVehicle: false,
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(3, false).Return(nil, ErrDriverNotFound)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriveRepo := repository.NewMockDriverRepository(ctrl)

			tt.setup(mockDriveRepo)

			vu := NewDriverUsecase(zap.NewNop().Sugar(), mockDriveRepo)
			got, err := vu.GetById(tt.driverId, tt.includeVehicle)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_driveUsecase_Create(t *testing.T) {
	tests := []struct {
		name    string
		driver  *entity.Driver
		setup   func(mockDriveRepo *repository.MockDriverRepository)
		wantErr bool
	}{
		{
			name: "Should create driver successfully",
			driver: &entity.Driver{
				Name:        "John",
				LastName:    "Doe",
				Email:       "john@test.com",
				Phone:       "21987654321",
				License:     "654321",
				LicenseType: "B",
			},
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().Create(&entity.Driver{
					Name:        "John",
					LastName:    "Doe",
					Email:       "john@test.com",
					Phone:       "21987654321",
					License:     "654321",
					LicenseType: "B",
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "Should return error for invalid driver",
			driver:  nil,
			setup:   func(mockDriveRepo *repository.MockDriverRepository) {},
			wantErr: true,
		},
		{
			name: "Should return error for invalid driver fields",
			driver: &entity.Driver{
				Name:        "",
				LastName:    "",
				Email:       "",
				Phone:       "",
				License:     "",
				LicenseType: "",
			},
			setup:   func(mockDriveRepo *repository.MockDriverRepository) {},
			wantErr: true,
		},
		{
			name: "Should return error for repository create failure",
			driver: &entity.Driver{
				Name:        "John",
				LastName:    "Doe",
				Email:       "john@test.com",
				Phone:       "21987654321",
				License:     "654321",
				LicenseType: "B",
			},
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().Create(&entity.Driver{
					Name:        "John",
					LastName:    "Doe",
					Email:       "john@test.com",
					Phone:       "21987654321",
					License:     "654321",
					LicenseType: "B",
				}).Return(fmt.Errorf("some error occurred"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriveRepo := repository.NewMockDriverRepository(ctrl)

			tt.setup(mockDriveRepo)

			vu := NewDriverUsecase(zap.NewNop().Sugar(), mockDriveRepo)
			err := vu.Create(tt.driver)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func Test_driveUsecase_AddVehicle(t *testing.T) {
	mockVehicle := &entity.Vehicle{
		Brand:        "Toyota",
		VehicleModel: "Corolla",
		Year:         2021,
		Plate:        "XYZ-9876",
	}
	tests := []struct {
		name     string
		driverId int
		vehicle  *entity.Vehicle
		setup    func(mockDriveRepo *repository.MockDriverRepository)
		wantErr  bool
	}{
		{
			name:     "Should add vehicle successfully",
			driverId: 1,
			vehicle:  mockVehicle,
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(1, false).Return(&entity.Driver{
					Name:        "Lucas",
					LastName:    "Moura",
					Email:       "lucas@test.com",
					Phone:       "123456789",
					License:     "123456",
					LicenseType: "A",
				}, nil)
				mockDriveRepo.EXPECT().AddVehicle(gomock.Any(), mockVehicle).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Should return error for invalid driver ID",
			driverId: -1,
			vehicle:  nil,
			setup:    func(mockDriveRepo *repository.MockDriverRepository) {},
			wantErr:  true,
		},
		{
			name:     "Should return error for invalid vehicle",
			driverId: 2,
			vehicle:  nil,
			setup:    func(mockDriveRepo *repository.MockDriverRepository) {},
			wantErr:  true,
		},
		{
			name:     "Should return error when vehicle fields are invalid",
			driverId: 4,
			vehicle: &entity.Vehicle{
				Brand:        "T",
				VehicleModel: "C",
				Year:         1885,
				Plate:        "XYZ-987",
			},
			setup:   func(mockDriveRepo *repository.MockDriverRepository) {},
			wantErr: true,
		},
		{
			name:     "Should return error for driver not found",
			driverId: 3,
			vehicle:  mockVehicle,
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(3, false).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:     "Should return error to get driver by id",
			driverId: 3,
			vehicle:  mockVehicle,
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(3, false).Return(nil, fmt.Errorf("some error occurred"))
			},
			wantErr: true,
		},
		{
			name:     "Should return error for repository add vehicle failure",
			driverId: 4,
			vehicle:  mockVehicle,
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(4, false).Return(&entity.Driver{
					Name:        "John",
					LastName:    "Doe",
					Email:       "john@test.com",
					Phone:       "987654321",
					License:     "654321",
					LicenseType: "B",
				}, nil)
				mockDriveRepo.EXPECT().AddVehicle(gomock.Any(), mockVehicle).Return(fmt.Errorf("some error occurred"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriveRepo := repository.NewMockDriverRepository(ctrl)

			tt.setup(mockDriveRepo)

			vu := NewDriverUsecase(zap.NewNop().Sugar(), mockDriveRepo)
			err := vu.AddVehicle(tt.driverId, tt.vehicle)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func Test_driveUsecase_Update(t *testing.T) {
	tests := []struct {
		name         string
		driverId     int
		updateDriver *entity.Driver
		setup        func(mockDriveRepo *repository.MockDriverRepository)
		wantErr      bool
	}{
		{
			name:     "Should update driver successfully",
			driverId: 1,
			updateDriver: &entity.Driver{
				Name:        "Lucas",
				LastName:    "Moura",
				Email:       "lucas@test.com",
				Phone:       "21987654321",
				License:     "21323232",
				LicenseType: "B",
			},
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(1, false).Return(&entity.Driver{
					Name:        "L",
					LastName:    "M",
					Email:       "lucas.test@test.com",
					Phone:       "123456789",
					License:     "123456",
					LicenseType: "A",
				}, nil)
				mockDriveRepo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:         "Should return error for invalid driver ID",
			driverId:     -1,
			updateDriver: nil,
			setup:        func(mockDriveRepo *repository.MockDriverRepository) {},
			wantErr:      true,
		},
		{
			name:         "Should return error to get driver by id",
			driverId:     2,
			updateDriver: new(entity.Driver),
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(2, false).Return(nil, fmt.Errorf("some error occurred"))
			},
			wantErr: true,
		},
		{
			name:         "Should return error for driver not found",
			driverId:     2,
			updateDriver: &entity.Driver{},
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(2, false).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:     "Should return error for invalid driver fields",
			driverId: 3,
			updateDriver: &entity.Driver{
				Name: "J",
			},
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(3, false).Return(&entity.Driver{
					Name:        "John",
					LastName:    "Doe",
					Email:       "john@test.com",
					Phone:       "987654321",
					License:     "654321",
					LicenseType: "B",
				}, nil)
			},
			wantErr: true,
		},
		{
			name:     "Should return error for repository update failure",
			driverId: 4,
			updateDriver: &entity.Driver{
				Name: "John",
			},
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().GetById(4, false).Return(&entity.Driver{
					Name:        "Johnn",
					LastName:    "Doe",
					Email:       "john@test.com",
					Phone:       "21987654321",
					License:     "654321",
					LicenseType: "B",
				}, nil)
				mockDriveRepo.EXPECT().Update(gomock.Any()).Return(fmt.Errorf("some error occurred"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriveRepo := repository.NewMockDriverRepository(ctrl)

			tt.setup(mockDriveRepo)

			vu := NewDriverUsecase(zap.NewNop().Sugar(), mockDriveRepo)
			err := vu.Update(tt.driverId, tt.updateDriver)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func Test_driveUsecase_Delete(t *testing.T) {
	tests := []struct {
		name     string
		driverId int
		setup    func(mockDriveRepo *repository.MockDriverRepository)
		wantErr  bool
	}{
		{
			name:     "Should delete driver successfully",
			driverId: 1,
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().Delete(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Should return error for invalid driver ID",
			driverId: -1,
			setup:    func(mockDriveRepo *repository.MockDriverRepository) {},
			wantErr:  true,
		},
		{
			name:     "Should return error for repository delete failure",
			driverId: 2,
			setup: func(mockDriveRepo *repository.MockDriverRepository) {
				mockDriveRepo.EXPECT().Delete(2).Return(fmt.Errorf("some error occurred"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriveRepo := repository.NewMockDriverRepository(ctrl)

			tt.setup(mockDriveRepo)

			vu := NewDriverUsecase(zap.NewNop().Sugar(), mockDriveRepo)
			err := vu.Delete(tt.driverId)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}
