package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Nikitapopov/Habbit/internal/apperror"
	"github.com/Nikitapopov/Habbit/internal/handlers"
	"github.com/Nikitapopov/Habbit/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

var _ handlers.Handler = &handler{}

const (
	usersURL = "/users"
	userURL  = "/user/:uuid"
)

type handler struct {
	service Service
	logger  logging.Logger
}

type Id struct {
	Id string
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handler {
	return &handler{
		service: *service,
		logger:  *logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.getUsers))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.createUser))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.getUserByUUID))
	router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodPatch, userURL, apperror.Middleware(h.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.deleteUser))
}

func (h *handler) getUsers(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(201)
	w.Write([]byte("these are users"))
	return nil
}

func (h *handler) createUser(w http.ResponseWriter, r *http.Request) error {
	var dto CreateUserDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	id, err := h.service.Create(context.TODO(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(&Id{Id: id})
	if err != nil {
		h.logger.Error("ID converting to result has failed")
	}
	w.Write(jsonResp)
	return nil
}

func (h *handler) getUserByUUID(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(200)
	w.Write([]byte("this is user by uuid"))
	return nil
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	w.Write([]byte("this is updating user"))
	return nil
}

func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	w.Write([]byte("this is partially updating user"))
	return nil
}

func (h *handler) deleteUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	w.Write([]byte("this is delete of user"))
	return nil
}
