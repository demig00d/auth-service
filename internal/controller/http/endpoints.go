package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/demig00d/auth-service/internal/controller/http/types"
	"github.com/demig00d/auth-service/internal/service"
)

/*
== Endpoints ==

Handlers
  * Authorize
  * Refresh

Errors:
	* if there is an error in the input data, the handler returns HTTP 400,
	* in case of validation error in business logic the handler returns HTTP 404
	* in case of other errors HTTP 503 is returned
*/

func (s *Server) Authorize(w http.ResponseWriter, r *http.Request) {
	var authReqBody types.AuthorizeReqBody
	err := json.NewDecoder(r.Body).Decode(&authReqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	guid, err := authReqBody.ToGUID()
	if err != nil {
		s.Logger.Info(err)
		errBody := types.NewErrorRespBody(err)

		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(errBody)
		return
	}

	access, refresh, err := s.Usecase.Authorize(guid)
	if err != nil {
		s.Logger.Info(err)
		errBody := types.NewErrorRespBody(err)

		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(errBody)
		return
	}

	body := types.NewRespBody(access, refresh)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(body)
}

func (s *Server) Refresh(w http.ResponseWriter, r *http.Request) {

	var refreshReqBody types.RefreshReqBody
	err := json.NewDecoder(r.Body).Decode(&refreshReqBody)
	if err != nil {
		s.Logger.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	guid, refresh, err := refreshReqBody.ToGUIDAndRefreshToken()
	if err != nil {
		s.Logger.Info(err)
		errBody := types.NewErrorRespBody(err)

		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(errBody)
		return
	}

	access, refresh, err := s.Usecase.Refresh(guid, refresh)
	if err != nil {
		s.Logger.Info(err)
		errBody := types.NewErrorRespBody(err)

		if errors.Is(err, service.ErrInvalidRefreshToken) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadGateway)
		}
		_ = json.NewEncoder(w).Encode(errBody)
		return
	}

	w.WriteHeader(http.StatusOK)
	body := types.NewRespBody(access, refresh)
	_ = json.NewEncoder(w).Encode(body)
}
