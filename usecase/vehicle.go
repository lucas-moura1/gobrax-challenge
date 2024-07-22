package usecase

import (
	"errors"

	"github.com/lucas-moura1/gobrax-challenge/entity"
	"github.com/lucas-moura1/gobrax-challenge/repository"
)

var ErrVehicleNotFound = errors.New("vehicle not found")

type VehicleUsecase interface {
	GetAll() ([]*entity.Vehicle, error)
	GetById(vehicleId int) (*entity.Vehicle, error)
	Update(vehicleId int, updateVehicle *entity.Vehicle) error
	Delete(vehicleId int) error
}

type vehicleUsecase struct {
	vRepo repository.VehicleRepository
}

func NewVehicleUsecase(vRepo repository.VehicleRepository) *vehicleUsecase {
	return &vehicleUsecase{vRepo: vRepo}
}

func (vu vehicleUsecase) GetAll() ([]*entity.Vehicle, error) {
	vehicles, err := vu.vRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (vu vehicleUsecase) GetById(vehicleId int) (*entity.Vehicle, error) {
	if vehicleId <= 0 {
		return nil, &entity.ErrorInvalidField{
			Message: []string{"vehicle id is invalid"},
		}
	}

	vehicle, err := vu.vRepo.GetById(vehicleId)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (vu vehicleUsecase) Update(vehicleId int, updateVehicle *entity.Vehicle) error {
	if vehicleId <= 0 {
		return &entity.ErrorInvalidField{
			Message: []string{"vehicle id is invalid"},
		}
	}

	vehicle, err := vu.vRepo.GetById(vehicleId)
	if err != nil {
		return err
	}
	if vehicle == nil {
		return ErrVehicleNotFound
	}

	if updateVehicle.Brand != "" {
		vehicle.Brand = updateVehicle.Brand
	}
	if updateVehicle.Plate != "" {
		vehicle.Plate = updateVehicle.Plate
	}
	if updateVehicle.VehicleModel != "" {
		vehicle.VehicleModel = updateVehicle.VehicleModel
	}
	if updateVehicle.Year != 0 {
		vehicle.Year = updateVehicle.Year
	}

	err = vehicle.Validate()
	if err != nil {
		return err
	}

	err = vu.vRepo.Update(vehicle)
	if err != nil {
		return err
	}
	return nil
}

func (vu vehicleUsecase) Delete(vehicleId int) error {
	if vehicleId <= 0 {
		return &entity.ErrorInvalidField{
			Message: []string{"vehicle id is invalid"},
		}
	}

	err := vu.vRepo.Delete(vehicleId)
	if err != nil {
		return err
	}
	return nil
}
