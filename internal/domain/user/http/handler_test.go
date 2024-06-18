package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	v1 "gymSystem/internal/domain/user/http"
	"gymSystem/internal/domain/user/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockUsecase es un mock del caso de uso para realizar pruebas unitarias.
type mockUsecase struct{}

// RegisterUser es una implementación mock del método RegisterUser del caso de uso.
func (m *mockUsecase) RegisterUser(ctx context.Context, request models.RegisterUserRequest) (models.RegisterUserResponse, error) {
	return models.RegisterUserResponse{
		UserID: 1,
	}, nil
}

// TestRegisterUser es una prueba unitaria para el controlador RegisterUser.
func TestRegisterUser(t *testing.T) {
	mockUc := &mockUsecase{}
	handler := v1.NewHandler(mockUc)

	// Configurar el cuerpo de la solicitud
	requestBody := models.RegisterUserRequest{
		Name:               "Billy",
		Lastname1:          "Rivera",
		Lastname2:          "Salinas",
		Email:              "billxd1958@gmail.com",
		Phone:              "7294574940",
		AccountTypeID:      1,
		SubscriptionCostID: 1,
		PaymentTypeID:      1,
		Amount:             29.0,
	}

	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(jsonRequestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.RegisterUser(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var actualResponse models.RegisterUserResponse
	err = json.NewDecoder(rr.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	expectedResponse := models.RegisterUserResponse{UserID: 1}
	if actualResponse != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", actualResponse, expectedResponse)
	}
}
