package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/lucas-moura1/gobrax-challenge/entity"
	"github.com/lucas-moura1/gobrax-challenge/usecase"
)

type vehicleRequest struct {
	Plate        string `json:"plate"`
	Brand        string `json:"brand"`
	VehicleModel string `json:"vehicleModel"`
	Year         int    `json:"year"`
}

type VehicleHandler struct {
	VehicleUsecase usecase.VehicleUsecase
}

func (vh VehicleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	vehicles, err := vh.VehicleUsecase.GetAll()
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	json.NewEncoder(w).Encode(vehicles)
}

func (vh VehicleHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vehicleId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Errorf("vehicleId must be a number"))
		return
	}

	vehicle, err := vh.VehicleUsecase.GetById(vehicleId)
	if err != nil {
		if reflect.TypeOf(err).String() == "*entity.ErrorInvalidField" {
			errorHandler(w, http.StatusBadRequest, err)
			return
		}
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	if vehicle == nil {
		errorHandler(w, http.StatusNotFound, fmt.Errorf("vehicle not found"))
		return
	}
	json.NewEncoder(w).Encode(vehicle)
}

func (vh VehicleHandler) Update(w http.ResponseWriter, r *http.Request) {
	vehicleId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Errorf("vehicleId must be a number"))
		return
	}

	var vehicle vehicleRequest
	err = json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	updateVehicle := &entity.Vehicle{
		Plate:        vehicle.Plate,
		Brand:        vehicle.Brand,
		VehicleModel: vehicle.VehicleModel,
		Year:         vehicle.Year,
	}

	err = vh.VehicleUsecase.Update(vehicleId, updateVehicle)
	if err != nil {
		if reflect.TypeOf(err).String() == "*entity.ErrorInvalidField" {
			errorHandler(w, http.StatusBadRequest, err)
			return
		}
		if err == usecase.ErrVehicleNotFound {
			errorHandler(w, http.StatusNotFound, err)
			return
		}
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (vh VehicleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vehicleId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Errorf("vehicleId must be a number"))
		return
	}

	err = vh.VehicleUsecase.Delete(vehicleId)
	if err != nil {
		if reflect.TypeOf(err).String() == "*entity.ErrorInvalidField" {
			errorHandler(w, http.StatusBadRequest, err)
			return
		}
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
