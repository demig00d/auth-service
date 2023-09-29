package service

import (
	"errors"
	"time"

	. "github.com/demig00d/auth-service/internal/model"
	"github.com/demig00d/auth-service/internal/repository"
	"github.com/demig00d/auth-service/pkg/logger"
)

var (
	ErrRefreshTokenCantBeEncrypted = errors.New("refresh token can't be encrypted")
	ErrCantGenerateAccessToken     = errors.New("can't generate access token")
)

type TokenService interface {
	GenerateAccessAndRefreshTokens(GUID) (AccessToken, RefreshToken, error)
	SaveRefreshToken(GUID, RefreshToken) error
	IsValidRefreshToken(GUID, RefreshToken) bool
}

type TokenServiceImpl struct {
	accessLifetime int64
	accessSecret   []byte
	repository     repository.Repository
	logger         logger.Logger
}

func NewTokenService(
	repo repository.Repository,
	accessLifetime int64,
	accessSecret []byte,
	logger logger.Logger,
) TokenServiceImpl {

	logger.SetPrefix("service - token ")

	return TokenServiceImpl{
		accessLifetime: accessLifetime,
		accessSecret:   accessSecret,
		repository:     repo,
		logger:         logger,
	}
}

func (s TokenServiceImpl) GenerateAccessAndRefreshTokens(guid GUID) (AccessToken, RefreshToken, error) {

	now := time.Now().Unix()

	payload := AccessTokenPayload{
		GUID: guid.String(),
		Iat:  now,
		Exp:  now + s.accessLifetime,
	}

	access, err := NewAccessToken(payload, s.accessSecret)
	if err != nil {
		s.logger.Debug(err)
		return AccessToken{}, RefreshToken{}, ErrCantGenerateAccessToken
	}

	refresh := NewRefreshToken()

	return access, refresh, nil
}

func (s TokenServiceImpl) SaveRefreshToken(guid GUID, token RefreshToken) error {

	encryptedRefresh, err := token.Encrypt()
	if err != nil {
		s.logger.Debug(err)
		return ErrRefreshTokenCantBeEncrypted
	}

	return s.repository.UpsertRefreshToken(guid, encryptedRefresh)
}

func (s TokenServiceImpl) IsValidRefreshToken(guid GUID, token RefreshToken) bool {
	storedToken, err := s.repository.FindRefreshToken(guid)
	if err != nil {
		s.logger.Debug(err)
		return false
	}

	return token.IsEqual(storedToken)
}
