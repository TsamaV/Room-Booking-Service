package auth

import (
	"Avito/internal/user"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	AdminUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	UserUUID  = uuid.MustParse("00000000-0000-0000-0000-000000000002")
)

type UserRepo interface {
	Create(u *user.User) (*user.User, error)
	GetByID(id uuid.UUID) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}

type AuthService struct {
	userRepo  UserRepo
	jwtSecret string
}

func NewAuthService(userRepo UserRepo, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (service *AuthService) Login(email, password string) (string, error) {
	existedUser, _ := service.userRepo.FindByEmail(email)
	if existedUser == nil {
		return "", errors.New(ErrWrongCredentials)
	}
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}
	return service.generateToken(existedUser.ID, existedUser.Role)
}

func (service *AuthService) Register(email, password string) (string, error) {
	existedUser, _ := service.userRepo.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	user := &user.User{
		Email:    email,
		PasswordHash: string(hashedPassword),
	}
	_, err = service.userRepo.Create(user)
	if err != nil {
		return "", err
	}
	return service.generateToken(user.ID, user.Role)
}

func (service *AuthService) DummyLogin(role string) (string, error) {
	if role != "admin" && role != "user" {
		return "", errors.New("Invalid role, must be admin or user")
	}

	var userID uuid.UUID
	if role == "admin" {
		userID = AdminUUID
	} else {
		userID = UserUUID
	}

	_, err := service.userRepo.GetByID(userID)
	if err != nil {
		user := &user.User{
			ID:    userID,
			Email: role + "@dummy.com",
			Role:  role,
		}
		if _, err := service.userRepo.Create(user); err != nil {
			return "", err
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(service.jwtSecret))
}

func (service *AuthService) generateToken(userID uuid.UUID, role string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID.String(),
        "role":    role,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    })
    return token.SignedString([]byte(service.jwtSecret))
}