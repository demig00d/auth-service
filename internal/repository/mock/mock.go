package mock

import (
	. "github.com/demig00d/auth-service/internal/model"
	"github.com/demig00d/auth-service/internal/repository"
)

type RepositoryMock struct {
	storage map[string]string
}

func NewRepositoryMock(guid GUID, encryptedToken Encrypted[RefreshToken]) RepositoryMock {

	storage := make(map[string]string)
	storage[guid.String()] = encryptedToken.String()

	return RepositoryMock{
		storage: storage,
	}
}

func (r RepositoryMock) UpsertRefreshToken(guid GUID, refreshToken Encrypted[RefreshToken]) error {
	r.storage[guid.String()] = refreshToken.String()
	return nil
}

func (r RepositoryMock) FindRefreshToken(guid GUID) (Encrypted[RefreshToken], error) {
	token, ok := r.storage[guid.String()]
	if !ok {
		return Encrypted[RefreshToken]{}, repository.ErrCantFindRefreshToken
	}

	encryptedToken, _ := EncryptedFromString[RefreshToken](token)
	return encryptedToken, nil
}
