package types

import (
	"github.com/demig00d/auth-service/internal/model"
)

type AuthorizeReqBody struct {
	GUID string `json:"GUID"`
}

type RefreshReqBody struct {
	GUID         string
	RefreshToken string `json:"refresh_token"`
}

func (b AuthorizeReqBody) ToGUID() (model.GUID, error) {
	return model.NewGUID(b.GUID)
}

func (b RefreshReqBody) ToGUIDAndRefreshToken() (model.GUID, model.RefreshToken, error) {
	guid, err := model.NewGUID(b.GUID)
	if err != nil {
		return model.GUID{}, model.RefreshToken{}, err
	}

	token, err := model.RefreshTokenFromString(b.RefreshToken)
	if err != nil {
		return model.GUID{}, model.RefreshToken{}, err
	}

	return guid, token, nil
}
