package entity

import (
	"net/mail"
	"regexp"

	"gorm.io/gorm"
)

const (
	LicenseTypeACC string = "ACC"
	LicenseTypeA   string = "A"
	LicenseTypeA1  string = "A1"
	LicenseTypeAB  string = "AB"
	LicenseTypeB   string = "B"
	LicenseTypeB1  string = "B1"
	LicenseTypeC   string = "C"
	LicenseTypeC1  string = "C1"
	LicenseTypeD   string = "D"
	LicenseTypeD1  string = "D1"
	LicenseTypeBE  string = "BE"
	LicenseTypeCE  string = "CE"
	LicenseTypeC1E string = "C1E"
	LicenseTypeDE  string = "DE"
	LicenseTypeD1E string = "D1E"

	regexPhone   string = `((\+|\(|0)?\d{1,3})?((\s|\)|\-))?(\d{10})$`
	regexLicense string = `^[a-zA-Z0-9]{6,11}$`
)

type Driver struct {
	gorm.Model
	Name        string
	LastName    string
	Email       string
	Phone       string
	License     string
	LicenseType string
	Vehicles    []Vehicle
}

func (d Driver) Validate() error {
	err := new(ErrorInvalidField)

	d.validateName(err)
	d.validateLastName(err)
	d.validateEmail(err)
	d.validatePhone(err)
	d.validateLicense(err)
	d.validateLicenseType(err)

	if len(err.Message) > 0 {
		return err
	}
	return nil
}

func (d Driver) validateName(err *ErrorInvalidField) {
	if d.Name == "" || len(d.Name) < 3 {
		err.Message = append(err.Message, "driver name is invalid")
	}
}

func (d Driver) validateLastName(err *ErrorInvalidField) {
	if d.LastName == "" || len(d.LastName) < 3 {
		err.Message = append(err.Message, "driver last name is invalid")
	}
}

func (d Driver) validateEmail(err *ErrorInvalidField) {
	_, errMail := mail.ParseAddress(d.Email)
	if errMail != nil {
		err.Message = append(err.Message, "driver email is invalid")
	}
}

func (d Driver) validatePhone(err *ErrorInvalidField) {
	if !regexp.MustCompile(regexPhone).MatchString(d.Phone) {
		err.Message = append(err.Message, "driver phone is invalid")
	}
}

func (d Driver) validateLicense(err *ErrorInvalidField) {
	if !regexp.MustCompile(regexLicense).MatchString(d.License) {
		err.Message = append(err.Message, "driver license is invalid")
	}
}

func (d Driver) validateLicenseType(err *ErrorInvalidField) {
	switch d.LicenseType {
	case LicenseTypeA, LicenseTypeB, LicenseTypeC, LicenseTypeD,
		LicenseTypeACC, LicenseTypeA1, LicenseTypeAB, LicenseTypeB1,
		LicenseTypeC1, LicenseTypeD1, LicenseTypeBE, LicenseTypeCE,
		LicenseTypeC1E, LicenseTypeDE, LicenseTypeD1E:
		return
	}
	err.Message = append(err.Message, "driver license type is invalid")
}
