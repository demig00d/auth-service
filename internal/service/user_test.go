package service

import (
	"testing"

	"github.com/demig00d/auth-service/internal/repository/mock"
	"github.com/demig00d/auth-service/pkg/logger"
)

func TestUserService(t *testing.T) {
	repo := mock.NewRepositoryMock(guid, encryptedToken)

	tokenService := NewTokenService(repo, lifetime, secret, logger.EmptyLogger)
	userService := NewUserService(tokenService, logger.EmptyLogger)

	_, refresh, err := userService.Authorize(guid)
	if err != nil {
		t.Error("Unexpected error", err)
	}

	_, _, err = userService.Refresh(guid, refresh)
	if err != nil {
		t.Error("Unexpected error", err)
	}

}
