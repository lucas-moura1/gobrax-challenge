package handler

import (
	"errors"
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

func TestDriverHandler_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mockDriverUsecase *usecase.MockDriverUsecase)
		wantErr bool
	}{
		{
			name: "Should return all drivers",
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().GetAll().Return(make([]*entity.Driver, 0), nil)
			},
			wantErr: false,
		},
		{
			name: "Should return error",
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().GetAll().Return(nil, fmt.Errorf("some error occurred"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriverUsecase := usecase.NewMockDriverUsecase(ctrl)
			tt.setup(mockDriverUsecase)

			dh := DriverHandler{
				DriverUsecase: mockDriverUsecase,
			}

			req := httptest.NewRequest(http.MethodGet, "/drivers", nil)
			respWritter := httptest.NewRecorder()

			dh.GetAll(respWritter, req)
			if tt.wantErr {
				assert.Equal(t, http.StatusInternalServerError, respWritter.Code)
				return
			}
			assert.Equal(t, http.StatusOK, respWritter.Code)
		})
	}
}

func TestDriverHandler_GetById(t *testing.T) {
	tests := []struct {
		name                string
		pathValue           string
		queryIncludeVehicle string
		setup               func(mockDriverUsecase *usecase.MockDriverUsecase)
		wantStatus          int
		wantError           bool
		wantErrMsg          string
	}{
		{
			name:      "Should return driver by ID successfully",
			pathValue: "1",
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().GetById(1, false).Return(new(entity.Driver), nil)
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name:                "Should return error when driver ID is not a number",
			pathValue:           "abc",
			queryIncludeVehicle: "true",
			setup:               func(mockDriverUsecase *usecase.MockDriverUsecase) {},
			wantStatus:          http.StatusBadRequest,
			wantError:           true,
			wantErrMsg:          "driverId must be a number",
		},
		{
			name:                "Should return error when includeVehicle is not a boolean",
			pathValue:           "1",
			queryIncludeVehicle: "abc",
			setup:               func(mockDriverUsecase *usecase.MockDriverUsecase) {},
			wantStatus:          http.StatusBadRequest,
			wantError:           true,
			wantErrMsg:          "includeVehicle must be a boolean",
		},
		{
			name:                "Should return error when driver is not found",
			pathValue:           "1",
			queryIncludeVehicle: "true",
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().GetById(1, true).Return(nil, nil)
			},
			wantStatus: http.StatusNotFound,
			wantError:  true,
			wantErrMsg: "driver not found",
		},
		{
			name:                "Should return bad request error when driverId is invalid",
			pathValue:           "0",
			queryIncludeVehicle: "true",
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().GetById(0, true).Return(nil, &entity.ErrorInvalidField{
					Message: []string{"driverId is invalid"},
				})
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			wantErrMsg: "driverId is invalid",
		},
		{
			name:                "Should return internal server error",
			pathValue:           "2",
			queryIncludeVehicle: "false",
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().GetById(2, false).Return(nil, fmt.Errorf("some error occurred"))
			},
			wantStatus: http.StatusInternalServerError,
			wantError:  true,
			wantErrMsg: "some error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriverUsecase := usecase.NewMockDriverUsecase(ctrl)
			tt.setup(mockDriverUsecase)

			dh := DriverHandler{
				DriverUsecase: mockDriverUsecase,
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/drivers/{id}?includeVehicle=%s", tt.queryIncludeVehicle), nil)
			req.SetPathValue("id", tt.pathValue)
			respWriter := httptest.NewRecorder()

			dh.GetById(respWriter, req)
			if tt.wantError {
				assert.Contains(t, respWriter.Body.String(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.wantStatus, respWriter.Code)
		})
	}
}

func TestDriverHandler_Create(t *testing.T) {
	mockBody := `{"name": "John", "lastName": "Doe", "email": "john.doe@example.com", "phone": "1234567890", "license": "ABC123", "licenseType": "car"}`
	tests := []struct {
		name        string
		requestBody string
		setup       func(mockDriverUsecase *usecase.MockDriverUsecase)
		wantStatus  int
		wantError   bool
		wantErrMsg  string
	}{
		{
			name:        "Should create driver successfully",
			requestBody: mockBody,
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name:        "Should return bad request error when request body is not a valid JSON",
			requestBody: `{"name"}`,
			setup:       func(mockDriverUsecase *usecase.MockDriverUsecase) {},
			wantStatus:  http.StatusInternalServerError,
			wantError:   true,
			wantErrMsg:  "{\"error\":\"invalid character '}' after object key\"}",
		},
		{
			name:        "Should return bad request error when returns invalid field error",
			requestBody: `{"name": "John", "lastName": "Doe", "email": "john.doe@example.com", "phone": "1234567890", "license": "21232123", "licenseType": "Y"}`,
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().Create(gomock.Any()).Return(&entity.ErrorInvalidField{
					Message: []string{"license is invalid", "licenseType is invalid"},
				})
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			wantErrMsg: "license is invalid,licenseType is invalid",
		},
		{
			name:        "Should return internal server error",
			requestBody: mockBody,
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().Create(gomock.Any()).Return(errors.New("some error occurred"))
			},
			wantStatus: http.StatusInternalServerError,
			wantError:  true,
			wantErrMsg: "some error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriverUsecase := usecase.NewMockDriverUsecase(ctrl)
			tt.setup(mockDriverUsecase)

			dh := DriverHandler{
				DriverUsecase: mockDriverUsecase,
			}

			reqBody := strings.NewReader(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/drivers", reqBody)
			respWriter := httptest.NewRecorder()

			dh.Create(respWriter, req)
			if tt.wantError {
				assert.Contains(t, respWriter.Body.String(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.wantStatus, respWriter.Code)
		})
	}
}
func TestDriverHandler_AddVehicle(t *testing.T) {
	mockBody := `{"plate": "ABC123", "brand": "Toyota", "vehicleModel": "Camry", "year": 2022}`
	tests := []struct {
		name        string
		pathValue   string
		requestBody string
		setup       func(mockDriverUsecase *usecase.MockDriverUsecase)
		wantStatus  int
		wantError   bool
		wantErrMsg  string
	}{
		{
			name:        "Should add vehicle successfully",
			pathValue:   "1",
			requestBody: mockBody,
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().AddVehicle(1, gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name:        "Should return bad request error when driverId is not a number",
			pathValue:   "abc",
			requestBody: mockBody,
			setup:       func(mockDriverUsecase *usecase.MockDriverUsecase) {},
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
			wantErrMsg:  "driverId must be a number",
		},
		{
			name:        "Should return bad request error when request body is not a valid JSON",
			pathValue:   "1",
			requestBody: `{"plate"}`,
			setup:       func(mockDriverUsecase *usecase.MockDriverUsecase) {},
			wantStatus:  http.StatusInternalServerError,
			wantError:   true,
			wantErrMsg:  "{\"error\":\"invalid character '}' after object key\"}",
		},
		{
			name:        "Should return bad request error when returns invalid field error",
			pathValue:   "1",
			requestBody: `{"plate": "ABC123", "brand": "T", "vehicleModel": "Camry", "year": 2022}`,
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().AddVehicle(1, gomock.Any()).Return(&entity.ErrorInvalidField{
					Message: []string{"vehicle brand is invalid"},
				})
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			wantErrMsg: "vehicle brand is invalid",
		},
		{
			name:        "Should return internal server error",
			pathValue:   "1",
			requestBody: `{"plate": "ABC123", "brand": "Toyota", "vehicleModel": "Camry", "year": 2022}`,
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().AddVehicle(1, gomock.Any()).Return(errors.New("some error occurred"))
			},
			wantStatus: http.StatusInternalServerError,
			wantError:  true,
			wantErrMsg: "some error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriverUsecase := usecase.NewMockDriverUsecase(ctrl)
			tt.setup(mockDriverUsecase)

			dh := DriverHandler{
				DriverUsecase: mockDriverUsecase,
			}

			reqBody := strings.NewReader(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/drivers/%s/vehicle", tt.pathValue), reqBody)
			req.SetPathValue("id", tt.pathValue)
			respWriter := httptest.NewRecorder()

			dh.AddVehicle(respWriter, req)
			if tt.wantError {
				assert.Contains(t, respWriter.Body.String(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.wantStatus, respWriter.Code)
		})
	}
}

func TestDriverHandler_Update(t *testing.T) {
	mockBody := `{"name": "John", "lastName": "Doe", "email": "john.doe@example.com", "phone": "1234567890", "license": "ABC123", "licenseType": "B"}`
	tests := []struct {
		name        string
		pathValue   string
		requestBody string
		setup       func(mockDriverUsecase *usecase.MockDriverUsecase)
		wantStatus  int
		wantError   bool
		wantErrMsg  string
	}{
		{
			name:        "Should update driver successfully",
			pathValue:   "1",
			requestBody: mockBody,
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().Update(1, gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name:        "Should return bad request error when driverId is not a number",
			pathValue:   "abc",
			requestBody: mockBody,
			setup:       func(mockDriverUsecase *usecase.MockDriverUsecase) {},
			wantStatus:  http.StatusBadRequest,
			wantError:   true,
			wantErrMsg:  "driverId must be a number",
		},
		{
			name:        "Should return bad request error when request body is not a valid JSON",
			pathValue:   "1",
			requestBody: `{"name"}`,
			setup:       func(mockDriverUsecase *usecase.MockDriverUsecase) {},
			wantStatus:  http.StatusInternalServerError,
			wantError:   true,
			wantErrMsg:  "{\"error\":\"invalid character '}' after object key\"}",
		},
		{
			name:        "Should return bad request error when returns invalid field error",
			pathValue:   "1",
			requestBody: `{"name": "John", "lastName": "Doe", "email": "john.doe@example.com", "phone": "1234567890", "license": "21232123", "licenseType": "Y"}`,
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().Update(1, gomock.Any()).Return(&entity.ErrorInvalidField{
					Message: []string{"license is invalid", "licenseType is invalid"},
				})
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			wantErrMsg: "license is invalid,licenseType is invalid",
		},
		{
			name:        "Should return not found error when driver is not found",
			pathValue:   "1",
			requestBody: mockBody,
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().Update(1, gomock.Any()).Return(usecase.ErrDriverNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantError:  true,
			wantErrMsg: "driver not found",
		},
		{
			name:        "Should return internal server error",
			pathValue:   "1",
			requestBody: mockBody,
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().Update(1, gomock.Any()).Return(errors.New("some error occurred"))
			},
			wantStatus: http.StatusInternalServerError,
			wantError:  true,
			wantErrMsg: "some error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriverUsecase := usecase.NewMockDriverUsecase(ctrl)
			tt.setup(mockDriverUsecase)

			dh := DriverHandler{
				DriverUsecase: mockDriverUsecase,
			}

			reqBody := strings.NewReader(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/drivers/{id}", reqBody)
			req.SetPathValue("id", tt.pathValue)
			respWriter := httptest.NewRecorder()

			dh.Update(respWriter, req)
			if tt.wantError {
				assert.Contains(t, respWriter.Body.String(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.wantStatus, respWriter.Code)
		})
	}
}

func TestDriverHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		pathValue  string
		setup      func(mockDriverUsecase *usecase.MockDriverUsecase)
		wantStatus int
		wantError  bool
		wantErrMsg string
	}{
		{
			name:      "Should delete driver successfully",
			pathValue: "1",
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().Delete(1).Return(nil)
			},
			wantStatus: http.StatusNoContent,
			wantError:  false,
		},
		{
			name:       "Should return bad request error when driverId is not a number",
			pathValue:  "abc",
			setup:      func(mockDriverUsecase *usecase.MockDriverUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			wantErrMsg: "driverId must be a number",
		},
		{
			name:      "Should return bad request error when driverId is invalid",
			pathValue: "0",
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().Delete(0).Return(&entity.ErrorInvalidField{
					Message: []string{"driverId is invalid"},
				})
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			wantErrMsg: "driverId is invalid",
		},
		{
			name:      "Should return internal server error",
			pathValue: "1",
			setup: func(mockDriverUsecase *usecase.MockDriverUsecase) {
				mockDriverUsecase.EXPECT().Delete(1).Return(errors.New("some error occurred"))
			},
			wantStatus: http.StatusInternalServerError,
			wantError:  true,
			wantErrMsg: "some error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDriverUsecase := usecase.NewMockDriverUsecase(ctrl)
			tt.setup(mockDriverUsecase)

			dh := DriverHandler{
				DriverUsecase: mockDriverUsecase,
			}

			req := httptest.NewRequest(http.MethodDelete, "/drivers/{id}", nil)
			req.SetPathValue("id", tt.pathValue)
			respWriter := httptest.NewRecorder()

			dh.Delete(respWriter, req)
			if tt.wantError {
				assert.Contains(t, respWriter.Body.String(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.wantStatus, respWriter.Code)
		})
	}
}
