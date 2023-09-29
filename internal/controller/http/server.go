package http

import (
	"net/http"

	"github.com/demig00d/auth-service/internal/service"
	"github.com/demig00d/auth-service/pkg/logger"
)

type Server struct {
	Logger  logger.Logger
	routes  []route
	Usecase service.UserService
}

// Server implements http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range s.routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method == route.method {
				w.Header().Set("Content-Type", "application/json")
				route.handler(w, r)
				return
			}
		}
	}
	http.Error(w, "500 internal server error", http.StatusInternalServerError)

}

func NewServer(uc service.UserService, logger logger.Logger) *Server {

	logger.SetPrefix("controller - http ")

	s := Server{
		routes:  []route{},
		Logger:  logger,
		Usecase: uc,
	}

	s.routes = []route{
		newRoute("POST", "/authorize", s.Authorize),
		newRoute("POST", "/refresh", s.Refresh),
	}

	return &s

}
