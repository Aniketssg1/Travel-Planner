package service

import (
	"context"
	"os"
	"regexp"
	"time"

	apperrors "backend/errors"
	"backend/middleware"
	"backend/models"
	"backend/repository"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

type RegisterResult struct {
	Message string
}

type LoginResult struct {
	Token string
	User  map[string]interface{}
}

func (s *AuthService) Register(ctx context.Context, input models.RegisterInput) (*RegisterResult, error) {
	if !emailRegex.MatchString(input.Email) {
		return nil, apperrors.BadRequest("invalid email format")
	}

	if len(input.Password) < 6 {
		return nil, apperrors.BadRequest("password must be at least 6 characters")
	}

	_, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err == nil {
		return nil, apperrors.Conflict("email already registered")
	}
	if err != mongo.ErrNoDocuments {
		return nil, apperrors.Internal("database error")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, apperrors.Internal("failed to process password")
	}

	user := models.User{
		ID:        primitive.NewObjectID(),
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Role:      "user",
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.Internal("failed to create user")
	}

	return &RegisterResult{Message: "User registered successfully"}, nil
}

func (s *AuthService) Login(ctx context.Context, input models.LoginInput) (*LoginResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err == mongo.ErrNoDocuments {
		return nil, apperrors.Unauthorized("invalid email or password")
	}
	if err != nil {
		return nil, apperrors.Internal("database error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, apperrors.Unauthorized("invalid email or password")
	}

	now := time.Now()
	claims := &middleware.Claims{
		UserID: user.ID.Hex(),
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "travel-planner-api",
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, apperrors.Internal("failed to generate token")
	}

	return &LoginResult{
		Token: tokenStr,
		User: map[string]interface{}{
			"id":    user.ID.Hex(),
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	}, nil
}
