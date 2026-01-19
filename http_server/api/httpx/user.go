package httpx

import (
	"code_processor/http_server/api"
	"code_processor/http_server/api/dto"
	"code_processor/http_server/models"
	"code_processor/http_server/usecases"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreatePostUserRegisterHandlerRequest(r *http.Request) (*dto.PostUserRegisterHandlerRequest, error) {
	var req dto.PostUserRegisterHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, api.ErrBadRequest{}
	}
	return &req, nil
}

func CreatePostUserLoginHandlerRequest(r *http.Request) (*dto.PostUserLoginHandlerRequest, error) {
	var req dto.PostUserLoginHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, api.ErrBadRequest{}
	}
	return &req, nil
}

type User struct {
	service usecases.User
}

func NewUserHandler(service usecases.User) *User {
	return &User{service: service}
}

// @Summary Register a user
// @Description Registers new user and adds them to data base
// @Tags user
// @Accept  json
// @Param user body dto.PostUserRegisterHandlerRequest true "user login and password"
// @Success 201 "Created"
// @Failure 400 {object} HTTPError "Bad request"
// @Failure 208 {object} HTTPError "Key already exists"
// @Router /user/register [post]
func (s *User) postRegisterHandler(w http.ResponseWriter, r *http.Request) {
	req, err := CreatePostUserRegisterHandlerRequest(r)
	if err != nil {
		WriteResponse(w, err, nil)
		return
	}

	newUser := models.User{Id: uuid.New(), Login: req.Username, Password: req.Password}

	err = s.service.PostRegister(newUser)
	WriteResponse(w, err, nil)
}

// @Summary Login a user
// @Description Logins new user and creates new session for him
// @Tags user
// @Accept  json
// @Param name body types.PostUserLoginHandlerRequest true "user login and password"
// @Success 200 {user} types.PostUserLoginHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Router /user/login [post]
func (s *User) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	req, err := CreatePostUserRegisterHandlerRequest(r)
	if err != nil {
		WriteResponse(w, err, nil)
		return
	}

	SessionId, err := s.service.PostLogin(req.Username, req.Password)
	WriteResponse(w, err, dto.PostUserLoginHandlerResponse{Token: SessionId.String()})
}

// WithUserHandlers registers user-related HTTP handlers.
func (s *User) WithUserHandlers(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", s.postRegisterHandler)
		r.Post("/login", s.postLoginHandler)
	})
}
