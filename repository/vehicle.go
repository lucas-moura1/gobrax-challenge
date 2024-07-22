package repository

import (
	"errors"

	"github.com/lucas-moura1/gobrax-challenge/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type VehicleRepository interface {
	GetAll() ([]*entity.Vehicle, error)
	GetById(vehicleId int) (*entity.Vehicle, error)
	Update(vehicle *entity.Vehicle) error
	Delete(vehicleId int) error
}

type vehicleRepository struct {
	log *zap.SugaredLogger
	db  *gorm.DB
}

func NewVehicleRepository(log *zap.SugaredLogger, db *gorm.DB) *vehicleRepository {
	return &vehicleRepository{log: log, db: db}
}

func (vr vehicleRepository) GetAll() ([]*entity.Vehicle, error) {
	var vehicles []*entity.Vehicle
	err := vr.db.Find(&vehicles).Error
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (vr vehicleRepository) GetById(vehicleId int) (*entity.Vehicle, error) {
	vehicle := new(entity.Vehicle)
	err := vr.db.First(vehicle, vehicleId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		vr.log.Errorw("error getting vehicle by id", "vehicleId", vehicleId, "error", err)
		return nil, err
	}
	return vehicle, nil
}

func (vr vehicleRepository) Update(vehicle *entity.Vehicle) error {
	err := vr.db.Save(vehicle).Error
	if err != nil {
		vr.log.Errorw("error updating vehicle", "vehicle", vehicle, "error", err)
		return err
	}
	return nil
}

func (vr vehicleRepository) Delete(vehicleId int) error {
	err := vr.db.Delete(&entity.Vehicle{}, vehicleId).Error
	if err != nil {
		vr.log.Errorw("error deleting vehicle", "vehicleId", vehicleId, "error", err)
		return err
	}
	return nil
}
