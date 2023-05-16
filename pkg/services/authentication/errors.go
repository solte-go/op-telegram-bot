package authentication

import "strings"

type userError string

func (ue userError) Error() string {
	return string(ue)
}

func (ue userError) Public() string {
	s := strings.Split(string(ue), " ")
	s[0] = strings.Title(s[0])
	return strings.Join(s, " ")
}

const (
	ErrPasswordRequired               userError = "password required"
	ErrPasswordNotCompileRequirements userError = "password not valid"
	ErrInvalidPassword                userError = "password not valid"
	ErrTokenGenWithError              userError = "token generated with error"
	ErrTokenRequired                  userError = "hashed session token required"
	ErrEmailIsRequired                userError = "email most be not null"
	ErrEmailNotValid                  userError = "email not valid"
	ErrNotFound                       userError = "object not found"
	ErrEmailTaken                     userError = "email address already taken"
)
