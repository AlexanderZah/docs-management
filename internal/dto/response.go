package dto

type RegisterResponse struct {
	Login string `json:"login"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type UploadDocResponse struct {
	Json map[string]interface{} `json:"json"`
	File string                 `json:"File"`
}
