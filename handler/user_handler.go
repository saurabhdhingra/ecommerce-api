package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"ecommerce-api/domain"
	"ecommerce-api/utils"
)


func (h *APIHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.Service.Signup(req.Username, req.Password)
	if err != nil {
		respondError(w, http.StatusConflict, "Username already taken or invalid input") // Be generic
		return
	}

	token, _ := h.JWTService.GenerateToken(user.ID, user.IsAdmin)
	respondJSON(w, http.StatusCreated, map[string]interface{}{"message": "User created", "token": token})
}

func (h *utils.APIHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, _ := h.JWTService.GenerateToken(user.ID, user.IsAdmin)
	respondJSON(w, http.StatusOK, map[string]interface{}{"message": "Login successful", "token": token, "is_admin": user.IsAdmin})
}