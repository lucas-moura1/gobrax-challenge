package usecase

import (
	"errors"

	"github.com/lucas-moura1/gobrax-challenge/entity"
	"github.com/lucas-moura1/gobrax-challenge/repository"
	"go.uber.org/zap"
)

var ErrDriverNotFound = errors.New("driver not found")

type DriverUsecase interface {
	GetAll() ([]*entity.Driver, error)
	GetById(driverId int, includeVehicle bool) (*entity.Driver, error)
	Create(driver *entity.Driver) error
	AddVehicle(driverId int, vehicle *entity.Vehicle) error
	Update(driverId int, driver *entity.Driver) error
	Delete(driverId int) error
}

type driverUsecase struct {
	log   *zap.SugaredLogger
	dRepo repository.DriverRepository
}

func NewDriverUsecase(log *zap.SugaredLogger, dRepo repository.DriverRepository) *driverUsecase {
	return &driverUsecase{log: log, dRepo: dRepo}
}

func (du driverUsecase) GetAll() ([]*entity.Driver, error) {
	drivers, err := du.dRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return drivers, nil
}

func (du driverUsecase) GetById(driverId int, includeVehicle bool) (*entity.Driver, error) {
	if driverId <= 0 {
		return nil, &entity.ErrorInvalidField{
			Message: []string{"driver id is invalid"},
		}
	}
	driver, err := du.dRepo.GetById(driverId, includeVehicle)
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func (du driverUsecase) Create(driver *entity.Driver) error {
	if driver == nil {
		return &entity.ErrorInvalidField{
			Message: []string{"driver is invalid"},
		}
	}
	err := driver.Validate()
	if err != nil {
		return err
	}

	err = du.dRepo.Create(driver)
	if err != nil {
		return err
	}
	return nil
}

func (du driverUsecase) AddVehicle(driverId int, vehicle *entity.Vehicle) error {
	if driverId <= 0 {
		return &entity.ErrorInvalidField{
			Message: []string{"driver id is invalid"},
		}
	}
	if vehicle == nil {
		return &entity.ErrorInvalidField{
			Message: []string{"vehicle is invalid"},
		}
	}

	err := vehicle.Validate()
	if err != nil {
		return err
	}

	driver, err := du.dRepo.GetById(driverId, false)
	if err != nil {
		return err
	}
	if driver == nil {
		return ErrDriverNotFound
	}
	err = du.dRepo.AddVehicle(driver, vehicle)
	if err != nil {
		return err
	}
	return nil
}

func (du driverUsecase) Update(driverId int, updateDriver *entity.Driver) error {
	if driverId <= 0 {
		return &entity.ErrorInvalidField{
			Message: []string{"driver id is invalid"},
		}
	}
	driver, err := du.dRepo.GetById(driverId, false)
	if err != nil {
		return err
	}
	if driver == nil {
		return ErrDriverNotFound
	}
	if updateDriver.Name != "" {
		driver.Name = updateDriver.Name
	}
	if updateDriver.LastName != "" {
		driver.LastName = updateDriver.LastName
	}
	if updateDriver.Email != "" {
		driver.Email = updateDriver.Email
	}
	if updateDriver.Phone != "" {
		driver.Phone = updateDriver.Phone
	}
	if updateDriver.License != "" {
		driver.License = updateDriver.License
	}
	if updateDriver.LicenseType != "" {
		driver.LicenseType = updateDriver.LicenseType
	}

	err = driver.Validate()
	if err != nil {
		return err
	}

	err = du.dRepo.Update(driver)
	if err != nil {
		return err
	}
	return nil
}

func (du driverUsecase) Delete(driverId int) error {
	if driverId <= 0 {
		return &entity.ErrorInvalidField{
			Message: []string{"driver id is invalid"},
		}
	}
	err := du.dRepo.Delete(driverId)
	if err != nil {
		return err
	}
	return nil
}
