package http

import (
	"code_processor/http_server/api/http/types"
	"code_processor/http_server/models"
	"code_processor/http_server/usecases"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

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
// @Param name body types.PostUserRegisterHandlerRequest true "user login and password"
// @Success 201 "Created"
// @Failure 400 {string} string "Bad request"
// @Failure 208 {string} string "Key already exists"
// @Router /user/register [post]
func (s *User) postRegisterHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatePostUserRegisterHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	newUser := models.User{Id: uuid.New().String(), Login: req.Username, Password: req.Password}

	err = s.service.PostRegister(newUser)
	types.ProcessError(w, err, nil, 201)
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
	req, err := types.CreatePostUserRegisterHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	SessionId, err := s.service.PostLogin(req.Username, req.Password)
	types.ProcessError(w, err, &types.PostUserLoginHandlerResponse{Token: SessionId}, 200)
}

// WithUserHandlers registers user-related HTTP handlers.
func (s *User) WithUserHandlers(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", s.postRegisterHandler)
		r.Post("/login", s.postLoginHandler)
	})
}
