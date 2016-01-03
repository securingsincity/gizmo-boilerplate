package service

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/gziphandler"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
	"net/http"
)

type (
	// JSONService will implement server.JSONService and
	// handle all requests to the server.
	// Config is a struct to contain all the needed
	// configuration for our JSONService
	Config struct {
		*config.Server
		MySQL   *config.MySQL
		MongoDB *config.MongoDB
		DB      gorm.DB
		Mongo   *mgo.Database
	}
	JSONService struct {
		Config *Config
	}
)

// NewJSONService will instantiate a JSONService
// with the given configuration.
func NewJSONService(cfg *Config) *JSONService {
	return &JSONService{
		Config: cfg,
	}
}

// Prefix returns the string prefix used for all endpoints within
// this service.
func (s *JSONService) Prefix() string {
	return "/api"
}

// Middleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we're using a GzipHandler middleware to
// compress our responses.
func (s *JSONService) Middleware(h http.Handler) http.Handler {
	return gziphandler.GzipHandler(h)
}

// JSONMiddleware provides a JSONEndpoint hook wrapped around all requests.
// In this implementation, we're using it to provide application logging and to check errors
// and provide generic responses.
func (s *JSONService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {

		status, res, err := j(r)
		if err != nil {
			print(err.Error())
			server.LogWithFields(r).WithFields(logrus.Fields{
				"error": err,
			}).Error("problems with serving request")
			return http.StatusServiceUnavailable, nil, &jsonErr{"sorry, this service is unavailable"}
		}

		server.LogWithFields(r).Info("success!")
		return status, res, nil
	}
}

// JSONEndpoints is a listing of all endpoints available in the JSONService.
func (s *JSONService) JSONEndpoints() map[string]map[string]server.JSONEndpoint {
	return map[string]map[string]server.JSONEndpoint{
		"/cats/{id}": map[string]server.JSONEndpoint{
			"GET":    s.ReadCat,
			"PUT":    s.UpdateCat,
			"DELETE": s.DeleteCat,
		},
		"/cats": map[string]server.JSONEndpoint{
			"POST": s.CreateCat,
			"GET":  s.ReadCats,
		},
		"/dogs": map[string]server.JSONEndpoint{
			"GET": s.ReadDogs,
		},
	}
}

type jsonErr struct {
	Err string `json:"error"`
}

func (e *jsonErr) Error() string {
	return e.Err
}
