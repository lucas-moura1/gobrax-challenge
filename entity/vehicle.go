package entity

import (
	"regexp"

	"gorm.io/gorm"
)

const (
	regexPlate string = `^[A-Z]{3}-\w{4}$`
)

type Vehicle struct {
	gorm.Model
	Brand        string
	VehicleModel string
	Year         int
	Plate        string
	DriverID     uint
}

func (v Vehicle) Validate() error {
	err := new(ErrorInvalidField)
	v.validateBrand(err)
	v.validateVehicleModel(err)
	v.validateYear(err)
	v.validatePlate(err)
	if len(err.Message) > 0 {
		return err
	}
	return nil
}

func (v Vehicle) validateBrand(err *ErrorInvalidField) {
	if v.Brand == "" || len(v.Brand) < 3 {
		err.Message = append(err.Message, "vehicle brand is invalid")
	}
}

func (v Vehicle) validateVehicleModel(err *ErrorInvalidField) {
	if v.VehicleModel == "" || len(v.VehicleModel) < 3 {
		err.Message = append(err.Message, "vehicle model is invalid")
	}
}

func (v Vehicle) validateYear(err *ErrorInvalidField) {
	// The first car was made in 1886
	if v.Year <= 1886 {
		err.Message = append(err.Message, "vehicle year is invalid")
	}
}

func (v Vehicle) validatePlate(err *ErrorInvalidField) {
	if !regexp.MustCompile(regexPlate).MatchString(v.Plate) {
		err.Message = append(err.Message, "vehicle plate is invalid")
	}
}
