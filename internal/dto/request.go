package dto

type Meta struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}

type RegisterRequest struct {
	AdminToken string `json:"token"`
	Login      string `json:"login"`
	Password   string `json:"pswd"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"pswd"`
}
