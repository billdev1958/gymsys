package v1

import (
	"encoding/json"
	"gymSystem/internal/domain/user"
	"gymSystem/internal/domain/user/models"
	"net/http"
)

type handler struct {
	uc user.Usecase
}

func NewHandler(uc user.Usecase) *handler {
	return &handler{uc: uc}
}

func (h handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	request := models.RegisterUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	response, err := h.uc.RegisterUser(r.Context(), request)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// request := models.RegisterUserRequest{}
// err := json.NewDecoder(r.Body).Decode(&request)
// if err != nil {
// 	http.Error(w, "invalid request", http.StatusBadRequest)
// 	return
// }

// response, err := h.uc.RegisterUser(r.Context(), request)
// if err != nil {
// 	http.Error(w, "error", http.StatusInternalServerError)
// 	return
// }

// w.Header().Set("Content-Type", "application/json")
// json.NewEncoder(w).Encode(response)
