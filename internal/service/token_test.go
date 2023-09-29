package service

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/demig00d/auth-service/internal/model"
	"github.com/demig00d/auth-service/internal/repository/mock"
	"github.com/demig00d/auth-service/pkg/logger"
)

var (
	guid, _           = model.NewGUID("B29AFC40-CA47-1067-B31D-00DD010662DA")
	token, _          = model.RefreshTokenFromString("G65cEalG2Yv9JGLTBvUfwG65cEalG2Yv9JGLTBvUfwd4ZDsqXTky4d4ZDsqXTky4")
	anotherToken, _   = model.RefreshTokenFromString("JjtUbYN0alR0CujQCK3w3gssWXYkJ8ddV6ohdRf9HQj01g8H7Veq9TZkyd4ZDsqX")
	encryptedToken, _ = token.Encrypt()
	lifetime          = int64(120)
	secret            = []byte("secret")
)

func TestGenerateAccessAndRefreshTokens(t *testing.T) {
	repo := mock.NewRepositoryMock(guid, encryptedToken)
	tokenService := NewTokenService(repo, lifetime, secret, logger.EmptyLogger)

	access, _, err := tokenService.GenerateAccessAndRefreshTokens(guid)
	if err != nil {
		t.Errorf("no errors expected, but got %s", err.Error())
	}

	CheckAccessToken(t, access, lifetime, secret)
}

func TestIsValidRefreshToken_Valid(t *testing.T) {
	repo := mock.NewRepositoryMock(guid, encryptedToken)
	tokenService := NewTokenService(repo, lifetime, secret, logger.EmptyLogger)

	isValid := tokenService.IsValidRefreshToken(guid, token)
	if !isValid {
		t.Error("expected valid")
	}

}

func TestIsValidRefreshToken_Invalid(t *testing.T) {
	repo := mock.NewRepositoryMock(guid, encryptedToken)
	tokenService := NewTokenService(repo, lifetime, secret, logger.EmptyLogger)

	isValid := tokenService.IsValidRefreshToken(guid, anotherToken)
	if isValid {
		t.Error("expected invalid")
	}
}

func CheckAccessToken(t *testing.T, access model.AccessToken, lifetime int64, secret []byte) {
	jwtParts := strings.Split(access.String(), ".")
	data, _ := base64.RawURLEncoding.DecodeString(jwtParts[1])

	var payload model.AccessTokenPayload
	json.Unmarshal(data, &payload)

	if payload.GUID != guid.String() {
		t.Errorf("GUID mismatch in payload, got != expected: %s != %s", guid.String(), payload.GUID)
	}

	gotLifetime := payload.Exp - payload.Iat
	if gotLifetime != lifetime {
		t.Errorf("lifetime mismatch in payload, got != expected: %d != %d", gotLifetime, lifetime)
	}

	if !access.HasValidSignature(secret) {
		t.Error("expected valid access token")
	}
}
