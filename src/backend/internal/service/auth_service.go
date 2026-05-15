package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"unihub-workshop/internal/middleware"
	"unihub-workshop/internal/model"
	"unihub-workshop/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepo
	secret   string
}

func NewAuthService(userRepo *repository.UserRepo, secret string) *AuthService {
	return &AuthService{userRepo: userRepo, secret: secret}
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByStudentID(ctx, req.StudentID)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Thử so sánh bằng bcrypt trước
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err == nil {
		// Khớp mã hóa bcrypt
	} else if user.PasswordHash == req.Password {
		// Khớp văn bản thô (Plaintext)
	} else {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := middleware.GenerateJWT(s.secret, user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.LoginResponse{Token: token, User: *user}, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
