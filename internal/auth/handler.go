package auth

import (
	"Avito/pkg/req"
	"Avito/pkg/res"
	"net/http"
)

type AuthHandler struct {
	AuthService *AuthService
}

func NewAuthHandler(router *http.ServeMux, service *AuthService) {
	handler := &AuthHandler{
		AuthService: service,
	}

	router.HandleFunc("POST /dummyLogin", handler.DummyLogin)
	router.HandleFunc("POST /login", handler.Login)
	router.HandleFunc("POST /register", handler.Register)
}

func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	body, err := req.Decode[DummyLoginRequest](r)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	token, err := h.AuthService.DummyLogin(body.Role)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	res.JSON(w, http.StatusOK, TokenResponse{Token: token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := req.Decode[LoginRequest](r)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	token, err := h.AuthService.Login(body.Email, body.Password)
	if err != nil {
		res.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	res.JSON(w, http.StatusOK, TokenResponse{Token: token})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := req.Decode[RegisterRequest](r)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	token, err := h.AuthService.Register(body.Email, body.Password)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	res.JSON(w, http.StatusCreated, TokenResponse{Token: token})
}