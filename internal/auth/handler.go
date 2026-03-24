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

// @Summary Получить тестовый JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body DummyLoginRequest true "Роль"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]interface{}
// @Router /dummyLogin [post]
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

// @Summary Авторизация
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Email и пароль"
// @Success 200 {object} TokenResponse
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
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

// @Summary Регистрация
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Email и пароль"
// @Success 201 {object} TokenResponse
// @Failure 400 {object} map[string]interface{}
// @Router /register [post]
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