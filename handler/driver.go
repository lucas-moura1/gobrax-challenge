package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/lucas-moura1/gobrax-challenge/entity"
	"github.com/lucas-moura1/gobrax-challenge/usecase"
)

type driverRequest struct {
	Name        string `json:"name"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	License     string `json:"license"`
	LicenseType string `json:"licenseType"`
}

type DriverHandler struct {
	DriverUsecase usecase.DriverUsecase
}

func (dh DriverHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	drivers, err := dh.DriverUsecase.GetAll()
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	json.NewEncoder(w).Encode(drivers)
}

func (dh DriverHandler) GetById(w http.ResponseWriter, r *http.Request) {
	driverId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Errorf("driverId must be a number"))
		return
	}

	includeVehicle := r.URL.Query().Get("includeVehicle")
	if includeVehicle == "" {
		includeVehicle = "false"
	}

	includeVehicleBool, err := strconv.ParseBool(includeVehicle)
	if err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Errorf("includeVehicle must be a boolean"))
		return
	}

	driver, err := dh.DriverUsecase.GetById(driverId, includeVehicleBool)
	if err != nil {
		if reflect.TypeOf(err).String() == "*entity.ErrorInvalidField" {
			errorHandler(w, http.StatusBadRequest, err)
			return
		}
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	if driver == nil {
		errorHandler(w, http.StatusNotFound, fmt.Errorf("driver not found"))
		return
	}
	json.NewEncoder(w).Encode(driver)
}

func (dh DriverHandler) Create(w http.ResponseWriter, r *http.Request) {
	driverReq := new(driverRequest)
	err := json.NewDecoder(r.Body).Decode(driverReq)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}

	driver := &entity.Driver{
		Name:        driverReq.Name,
		LastName:    driverReq.LastName,
		Email:       driverReq.Email,
		Phone:       driverReq.Phone,
		License:     driverReq.License,
		LicenseType: driverReq.LicenseType,
	}

	err = dh.DriverUsecase.Create(driver)
	if err != nil {
		if reflect.TypeOf(err).String() == "*entity.ErrorInvalidField" {
			errorHandler(w, http.StatusBadRequest, err)
			return
		}
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (dh DriverHandler) AddVehicle(w http.ResponseWriter, r *http.Request) {
	driverId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Errorf("driverId must be a number"))
		return
	}

	vehicleReq := new(vehicleRequest)
	err = json.NewDecoder(r.Body).Decode(vehicleReq)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}

	vehicle := &entity.Vehicle{
		Plate:        vehicleReq.Plate,
		Brand:        vehicleReq.Brand,
		VehicleModel: vehicleReq.VehicleModel,
		Year:         vehicleReq.Year,
	}
	err = dh.DriverUsecase.AddVehicle(driverId, vehicle)
	if err != nil {
		if reflect.TypeOf(err).String() == "*entity.ErrorInvalidField" {
			errorHandler(w, http.StatusBadRequest, err)
			return
		}
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (dh DriverHandler) Update(w http.ResponseWriter, r *http.Request) {
	driverId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Errorf("driverId must be a number"))
		return
	}

	driverReq := new(driverRequest)
	err = json.NewDecoder(r.Body).Decode(driverReq)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}

	driver := &entity.Driver{
		Name:        driverReq.Name,
		LastName:    driverReq.LastName,
		Email:       driverReq.Email,
		Phone:       driverReq.Phone,
		License:     driverReq.License,
		LicenseType: driverReq.LicenseType,
	}

	err = dh.DriverUsecase.Update(driverId, driver)
	if err != nil {
		if reflect.TypeOf(err).String() == "*entity.ErrorInvalidField" {
			errorHandler(w, http.StatusBadRequest, err)
			return
		}
		if errors.Is(err, usecase.ErrDriverNotFound) {
			errorHandler(w, http.StatusNotFound, err)
			return
		}
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (dh DriverHandler) Delete(w http.ResponseWriter, r *http.Request) {
	driverId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Errorf("driverId must be a number"))
		return
	}

	err = dh.DriverUsecase.Delete(driverId)
	if err != nil {
		if reflect.TypeOf(err).String() == "*entity.ErrorInvalidField" {
			errorHandler(w, http.StatusBadRequest, err)
			return
		}
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
