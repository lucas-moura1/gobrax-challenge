package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucas-moura1/gobrax-challenge/entity"
	"github.com/lucas-moura1/gobrax-challenge/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVehicleHandler_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mockVehicleUsecase *usecase.MockVehicleUsecase)
		wantErr bool
	}{
		{
			name: "Should return all vehicles",
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().GetAll().Return(make([]*entity.Vehicle, 0), nil)
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().GetAll().Return(nil, fmt.Errorf("some error occurred"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockVehicleUsecase := usecase.NewMockVehicleUsecase(ctrl)
			tt.setup(mockVehicleUsecase)

			vh := VehicleHandler{
				VehicleUsecase: mockVehicleUsecase,
			}

			req := httptest.NewRequest(http.MethodGet, "/vehicles", nil)
			respWritter := httptest.NewRecorder()

			vh.GetAll(respWritter, req)
			if tt.wantErr {
				assert.Equal(t, http.StatusInternalServerError, respWritter.Code)
				return
			}
			assert.Equal(t, http.StatusOK, respWritter.Code)
		})
	}
}

func TestVehicleHandler_GetById(t *testing.T) {
	tests := []struct {
		name         string
		pathValue    string
		setup        func(mockVehicleUsecase *usecase.MockVehicleUsecase)
		wantCode     int
		wantError    bool
		wantErrorMsg string
	}{
		{
			name:      "Should return vehicle by ID",
			pathValue: "1",
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().GetById(1).Return(new(entity.Vehicle), nil)
			},
			wantCode:  http.StatusOK,
			wantError: false,
		},
		{
			name:         "Should return bad request error when vehicleId is not a number",
			pathValue:    "abc",
			setup:        func(mockVehicleUsecase *usecase.MockVehicleUsecase) {},
			wantCode:     http.StatusBadRequest,
			wantError:    true,
			wantErrorMsg: "vehicleId must be a number",
		},
		{
			name:      "Should return bad request error when vehicleId is invalid",
			pathValue: "0",
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().GetById(0).Return(nil, &entity.ErrorInvalidField{
					Message: []string{"vehicle id is invalid"},
				})
			},
			wantCode:     http.StatusBadRequest,
			wantError:    true,
			wantErrorMsg: "vehicle id is invalid",
		},
		{
			name:      "Should return bad request error when vehicleId is not found",
			pathValue: "999",
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().GetById(999).Return(nil, nil)
			},
			wantCode:     http.StatusNotFound,
			wantError:    true,
			wantErrorMsg: "vehicle not found",
		},
		{
			name:      "Should return internal server error",
			pathValue: "2",
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().GetById(2).Return(nil, fmt.Errorf("some error occurred"))
			},
			wantCode:     http.StatusInternalServerError,
			wantError:    true,
			wantErrorMsg: "some error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockVehicleUsecase := usecase.NewMockVehicleUsecase(ctrl)
			tt.setup(mockVehicleUsecase)

			vh := VehicleHandler{
				VehicleUsecase: mockVehicleUsecase,
			}

			req := httptest.NewRequest(http.MethodGet, "/vehicles/{id}", nil)
			req.SetPathValue("id", tt.pathValue)
			respWriter := httptest.NewRecorder()

			vh.GetById(respWriter, req)
			if tt.wantError {
				assert.Contains(t, respWriter.Body.String(), tt.wantErrorMsg)
			}
			assert.Equal(t, tt.wantCode, respWriter.Code)
		})
	}
}

func TestVehicleHandler_Update(t *testing.T) {
	mockBody := `{"plate": "ABC123", "brand": "Toyota", "vehicleModel": "Corolla", "year": 2022}`
	tests := []struct {
		name         string
		pathValue    string
		requestBody  string
		setup        func(mockVehicleUsecase *usecase.MockVehicleUsecase)
		wantCode     int
		wantError    bool
		wantErrorMsg string
	}{
		{
			name:        "Should update vehicle successfully",
			pathValue:   "1",
			requestBody: mockBody,
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().Update(1, &entity.Vehicle{
					Plate:        "ABC123",
					Brand:        "Toyota",
					VehicleModel: "Corolla",
					Year:         2022,
				}).Return(nil)
			},
			wantCode:  http.StatusOK,
			wantError: false,
		},
		{
			name:         "Should return bad request error when vehicleId is not a number",
			pathValue:    "abc",
			requestBody:  mockBody,
			setup:        func(mockVehicleUsecase *usecase.MockVehicleUsecase) {},
			wantCode:     http.StatusBadRequest,
			wantError:    true,
			wantErrorMsg: "vehicleId must be a number",
		},
		{
			name:         "Should return bad request error when request body is invalid",
			pathValue:    "1",
			requestBody:  `{"plate"}`,
			setup:        func(mockVehicleUsecase *usecase.MockVehicleUsecase) {},
			wantCode:     http.StatusBadRequest,
			wantError:    true,
			wantErrorMsg: "invalid request body",
		},
		{
			name:        "Should return error when vehicle is not found",
			pathValue:   "2",
			requestBody: mockBody,
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().Update(2, gomock.Any()).Return(usecase.ErrVehicleNotFound)
			},
			wantCode:     http.StatusNotFound,
			wantError:    true,
			wantErrorMsg: "vehicle not found",
		},
		{
			name:        "Should return bad request error when vehicle brand is invalid",
			pathValue:   "3",
			requestBody: `{"brand": "T"}`,
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().Update(3, gomock.Any()).Return(&entity.ErrorInvalidField{
					Message: []string{"vehicle brand is invalid"},
				})
			},
			wantCode:     http.StatusBadRequest,
			wantError:    true,
			wantErrorMsg: "vehicle brand is invalid",
		},
		{
			name:        "Should return internal server error",
			pathValue:   "3",
			requestBody: mockBody,
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().Update(3, gomock.Any()).Return(fmt.Errorf("some error occurred"))
			},
			wantCode:     http.StatusInternalServerError,
			wantError:    true,
			wantErrorMsg: "some error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockVehicleUsecase := usecase.NewMockVehicleUsecase(ctrl)
			tt.setup(mockVehicleUsecase)

			vh := VehicleHandler{
				VehicleUsecase: mockVehicleUsecase,
			}

			req := httptest.NewRequest(http.MethodPut, "/vehicles/{id}", strings.NewReader(tt.requestBody))
			req.SetPathValue("id", tt.pathValue)
			respWriter := httptest.NewRecorder()

			vh.Update(respWriter, req)
			if tt.wantError {
				assert.Contains(t, respWriter.Body.String(), tt.wantErrorMsg)
			}
			assert.Equal(t, tt.wantCode, respWriter.Code)
		})
	}
}

func TestVehicleHandler_Delete(t *testing.T) {
	tests := []struct {
		name         string
		pathValue    string
		setup        func(mockVehicleUsecase *usecase.MockVehicleUsecase)
		wantCode     int
		wantError    bool
		wantErrorMsg string
	}{
		{
			name:      "Should delete vehicle successfully",
			pathValue: "1",
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().Delete(1).Return(nil)
			},
			wantCode:  http.StatusOK,
			wantError: false,
		},
		{
			name:         "Should return bad request error when vehicleId is not a number",
			pathValue:    "abc",
			setup:        func(mockVehicleUsecase *usecase.MockVehicleUsecase) {},
			wantCode:     http.StatusBadRequest,
			wantError:    true,
			wantErrorMsg: "vehicleId must be a number",
		},
		{
			name:      "Should return bad request error when vehicleId is invalid",
			pathValue: "0",
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().Delete(0).Return(&entity.ErrorInvalidField{
					Message: []string{"vehicle id is invalid"},
				})
			},
			wantCode:     http.StatusBadRequest,
			wantError:    true,
			wantErrorMsg: "vehicle id is invalid",
		},
		{
			name:      "Should return internal server error",
			pathValue: "3",
			setup: func(mockVehicleUsecase *usecase.MockVehicleUsecase) {
				mockVehicleUsecase.EXPECT().Delete(3).Return(fmt.Errorf("some error occurred"))
			},
			wantCode:     http.StatusInternalServerError,
			wantError:    true,
			wantErrorMsg: "some error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockVehicleUsecase := usecase.NewMockVehicleUsecase(ctrl)
			tt.setup(mockVehicleUsecase)

			vh := VehicleHandler{
				VehicleUsecase: mockVehicleUsecase,
			}

			req := httptest.NewRequest(http.MethodDelete, "/vehicles/{id}", nil)
			req.SetPathValue("id", tt.pathValue)
			respWriter := httptest.NewRecorder()

			vh.Delete(respWriter, req)
			if tt.wantError {
				assert.Contains(t, respWriter.Body.String(), tt.wantErrorMsg)
			}
			assert.Equal(t, tt.wantCode, respWriter.Code)
		})
	}
}
