package dto

type PostUserRegisterHandlerRequest struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type PostUserLoginHandlerRequest struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type PostUserLoginHandlerResponse struct {
	Token string `json:"token"`
}
