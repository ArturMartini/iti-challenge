package rest

import (
	"encoding/json"
	"github.com/arturmartini/iti-challenge/entities"
	"github.com/arturmartini/iti-challenge/service"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Endpoint struct {
	Path   string
	Func   func(w http.ResponseWriter, r *http.Request)
	Method string
}

const uri = "/v1/api/password-validate"

var (
	svc       service.Service
	router    = mux.NewRouter()
	endpoints = []Endpoint{
		{
			Path:   uri,
			Func:   validatePassword,
			Method: http.MethodPost,
		},
	}
)

func init() {
	svc = service.New()
	register()
}

func Start() error {
	log.Info("Http server start at")
	return http.ListenAndServe(":8080", router)
}

func register() {
	for _, e := range endpoints {
		router.HandleFunc(e.Path, e.Func).Methods(e.Method)
		log.WithFields(log.Fields{
			"path":   e.Path,
			"method": e.Method,
		}).Info("Register endpoint")
	}
}

func validatePassword(w http.ResponseWriter, r *http.Request) {
	password := entities.Password{}
	err := fromJson(w, r, &password)
	if err != nil {
		return
	}
	valid := svc.ValidateStrongPassword(password)
	handleResponse(entities.Response{Valid: valid}, w)
}

func fromJson(w http.ResponseWriter, r *http.Request, value interface{}) error {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warn("Error when read request body")
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	err = json.Unmarshal(bytes, value)
	if err != nil {
		log.WithError(err).Warn("Error when unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	return err
}

func handleResponse(value interface{}, w http.ResponseWriter) {
	bytes, err := json.Marshal(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unexpected error"))
		return
	}

	w.Write(bytes)
	w.WriteHeader(http.StatusOK)
}
