package types

import "github.com/demig00d/auth-service/internal/model"

type ErrorRespBody struct {
	Error string `json:"error"`
}

func NewErrorRespBody(err error) ErrorRespBody {
	return ErrorRespBody{
		Error: err.Error(),
	}
}

type RespBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewRespBody(access model.AccessToken, refresh model.RefreshToken) RespBody {
	return RespBody{
		AccessToken:  access.String(),
		RefreshToken: refresh.String(),
	}
}
