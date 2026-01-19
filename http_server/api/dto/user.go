package dto

type PostUserRegisterHandlerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PostUserLoginHandlerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PostUserLoginHandlerResponse struct {
	Token string `json:"token"`
}
