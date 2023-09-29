package service

import (
	"errors"

	. "github.com/demig00d/auth-service/internal/model"
	"github.com/demig00d/auth-service/pkg/logger"
)

var (
	ErrCantSaveRefreshToken = errors.New("can't save refresh token")
	ErrInvalidRefreshToken  = errors.New("refresh token is invalid")
)

type UserService interface {
	Authorize(GUID) (AccessToken, RefreshToken, error)
	Refresh(GUID, RefreshToken) (AccessToken, RefreshToken, error)
}

type UserServiceImpl struct {
	tokenService TokenService
	logger       logger.Logger
}

func NewUserService(tokenService TokenService, logger logger.Logger) UserServiceImpl {
	logger.SetPrefix("service - user ")

	return UserServiceImpl{
		tokenService: tokenService,
		logger:       logger,
	}
}

func (s UserServiceImpl) Authorize(guid GUID) (AccessToken, RefreshToken, error) {
	access, refresh, err := s.tokenService.GenerateAccessAndRefreshTokens(guid)
	if err != nil {
		return AccessToken{}, RefreshToken{}, err
	}

	err = s.tokenService.SaveRefreshToken(guid, refresh)
	if err != nil {
		return AccessToken{}, RefreshToken{}, err
	}

	return access, refresh, nil
}

func (s UserServiceImpl) Refresh(guid GUID, token RefreshToken) (AccessToken, RefreshToken, error) {

	isValid := s.tokenService.IsValidRefreshToken(guid, token)
	if !isValid {
		return AccessToken{}, RefreshToken{}, ErrInvalidRefreshToken
	}
	newAccess, newRefresh, err := s.tokenService.GenerateAccessAndRefreshTokens(guid)
	if err != nil {
		return AccessToken{}, RefreshToken{}, err
	}

	err = s.tokenService.SaveRefreshToken(guid, newRefresh)
	if err != nil {
		return AccessToken{}, RefreshToken{}, ErrCantSaveRefreshToken
	}

	return newAccess, newRefresh, nil
}
