package repository

import (
	"errors"

	"github.com/lucas-moura1/gobrax-challenge/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DriverRepository interface {
	GetAll() ([]*entity.Driver, error)
	GetById(driverId int, includeVehicle bool) (*entity.Driver, error)
	Create(driver *entity.Driver) error
	AddVehicle(driver *entity.Driver, vehicle *entity.Vehicle) error
	Update(driver *entity.Driver) error
	Delete(driverId int) error
}

type driverRepository struct {
	log *zap.SugaredLogger
	db  *gorm.DB
}

func NewDriverRepository(log *zap.SugaredLogger, db *gorm.DB) *driverRepository {
	return &driverRepository{log: log, db: db}
}

func (dr driverRepository) GetAll() ([]*entity.Driver, error) {
	var drivers []*entity.Driver
	err := dr.db.Find(&drivers).Error
	if err != nil {
		return nil, err
	}
	return drivers, nil
}

func (dr driverRepository) GetById(driverId int, includeVehicle bool) (*entity.Driver, error) {
	driver := new(entity.Driver)
	query := dr.db
	if includeVehicle {
		query = query.Preload("Vehicles")
	}
	err := query.First(driver, driverId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		dr.log.Errorw("error getting driver by id", "driverId", driverId, "error", err)
		return nil, err
	}
	return driver, nil
}

func (dr driverRepository) Create(driver *entity.Driver) error {
	return dr.db.Create(driver).Error
}

func (dr driverRepository) AddVehicle(driver *entity.Driver, vehicle *entity.Vehicle) error {
	err := dr.db.Model(driver).Association("Vehicles").Append(vehicle)
	if err != nil {
		dr.log.Errorw("error adding vehicle to driver",
			"driverId", driver.ID, "vehicle", vehicle, "error", err)
		return err
	}
	return nil
}

func (dr driverRepository) Update(driver *entity.Driver) error {
	err := dr.db.Save(driver).Error
	if err != nil {
		dr.log.Errorw("error updating driver", "driver", driver, "error", err)
		return err
	}
	return nil
}

func (dr driverRepository) Delete(driverId int) error {
	err := dr.db.Delete(&entity.Driver{}, driverId).Error
	if err != nil {
		dr.log.Errorw("error deleting driver", "driverId", driverId, "error", err)
		return err
	}
	return nil
}
