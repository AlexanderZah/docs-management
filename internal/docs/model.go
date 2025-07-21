package docs

import (
	"time"
)

type Document struct {
	ID        int32
	Name      string
	IsFile    bool
	Public    bool
	Token     string
	Mime      string
	Grants    []string
	Json      map[string]interface{}
	Content   []byte
	CreatedAt time.Time
}
