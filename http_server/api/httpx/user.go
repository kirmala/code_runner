package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/api"
	"github.com/kirmala/code_runner/http_server/api/dto"
	"github.com/kirmala/code_runner/http_server/domain"
	"github.com/kirmala/code_runner/http_server/service"
	"github.com/labstack/echo/v5"
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
	service service.User
}

func NewUserHandler(service service.User) *User {
	return &User{service: service}
}

// @Summary Register a user
// @Description Registers new user and adds them to data base
// @Tags user
// @Accept  json
// @Param user body dto.PostUserRegisterHandlerRequest true "user login and password"
// @Success 201 "Created"
// @Failure 400 {object} dto.Error "Bad request"
// @Failure 409 {object} dto.Error "Key already exists"
// @Router /user/register [post]
func (s *User) postRegisterHandler(c *echo.Context) error {
	req, err := CreatePostUserRegisterHandlerRequest(c.Request())
	if err != nil {
		return err
	}

	newUser := domain.User{Id: uuid.New(), Login: req.Login, Password: req.Password}

	err = s.service.Register(newUser)
	if err != nil {
		return err
	}

	
	return c.NoContent(http.StatusCreated)
}

// @Summary Login a user
// @Description Logins new user and creates new session for him
// @Tags user
// @Accept  json
// @Param user body dto.PostUserLoginHandlerRequest true "user login and password"
// @Success 200 {object} dto.PostUserLoginHandlerResponse
// @Failure 400 {object} dto.Error "Bad request"
// @Failure 401 {object} dto.Error "Unauthorized"
// @Router /user/login [post]
func (s *User) postLoginHandler(c *echo.Context) error {
	req, err := CreatePostUserLoginHandlerRequest(c.Request())
	if err != nil {
		return err
	}

	SessionId, err := s.service.Login(c.Request().Context(), req.Login, req.Password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &dto.PostUserLoginHandlerResponse{Token: SessionId.String()})
}

// WithUserHandlers registers user-related HTTP handlers.
func (s *User) WithUserHandlers(g *echo.Group) {
	g.POST("/user/register", echo.HandlerFunc(s.postRegisterHandler))
	g.POST("/user/login", echo.HandlerFunc(s.postLoginHandler))
}
