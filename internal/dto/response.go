package dto

type RegisterResponse struct {
	Login string `json:"login"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LogoutResponse struct {
	Response map[string]bool `json:"response"`
}

type UploadDocResponse struct {
	Json map[string]interface{} `json:"json"`
	File string                 `json:"File"`
}

type GetDocResponse struct {
	ID        int32
	Name      string
	Mime      string
	IsFile    bool
	Public    bool
	CreatedAt string
	Grants    []string
}
