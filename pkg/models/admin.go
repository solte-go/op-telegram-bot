package models

type Admin struct {
	Id             uint
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	HashedPassword string `json:"-"`
	Token          string `json:"-"`
	HashedToken    string `json:"-"`
}
