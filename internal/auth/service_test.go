package auth_test

import (
	"Avito/internal/auth"
	"Avito/internal/user"
	"testing"

	"github.com/google/uuid"
)

type MockUserRepo struct {
	user *user.User
}

func (m *MockUserRepo) Create(u *user.User) (*user.User, error) {
	return u, nil
}

func (m *MockUserRepo) GetByID(id uuid.UUID) (*user.User, error) {
	if m.user != nil {
		return m.user, nil
	}
	return nil, nil
}

func (m *MockUserRepo) FindByEmail(email string) (*user.User, error) {
	if m.user != nil && m.user.Email == email {
		return m.user, nil
	}
	return nil, nil
}

func TestDummyLoginAdmin(t *testing.T) {
	service := auth.NewAuthService(&MockUserRepo{}, "secret")

	token, err := service.DummyLogin("admin")
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("expected token, got empty string")
	}
}

func TestDummyLoginUser(t *testing.T) {
	service := auth.NewAuthService(&MockUserRepo{}, "secret")

	token, err := service.DummyLogin("user")
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("expected token, got empty string")
	}
}

func TestDummyLoginInvalidRole(t *testing.T) {
	service := auth.NewAuthService(&MockUserRepo{}, "secret")

	_, err := service.DummyLogin("superadmin")
	if err == nil {
		t.Fatal("expected error for invalid role, got nil")
	}
}

func TestRegisterSuccess(t *testing.T) {
	service := auth.NewAuthService(&MockUserRepo{}, "secret")

	token, err := service.Register("test@test.com", "password")
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("expected token, got empty string")
	}
}

func TestRegisterAlreadyExists(t *testing.T) {
	service := auth.NewAuthService(&MockUserRepo{
		user: &user.User{
			Email: "test@test.com",
		},
	}, "secret")

	_, err := service.Register("test@test.com", "password")
	if err == nil {
		t.Fatal("expected error for existing user, got nil")
	}
}