package docs

type Document struct {
	Name   string
	Public bool
	Token  string
	Mime   string
	Grants []string
	Json   map[string]interface{}
	File   []byte
}
