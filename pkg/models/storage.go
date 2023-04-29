package models

type Page struct {
	UserID   int
	UserName string
	URLId    int
	URL      string
}

type Words struct {
	Topic   string
	Letter  string
	Suomi   string
	Russian string
	English string
}
