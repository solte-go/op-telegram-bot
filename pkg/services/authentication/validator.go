package authentication

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/services/rand"
)

type userValidationFunc func(*models.Admin) error

func runUserValidationFunc(user *models.Admin, fns ...userValidationFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (uv *userValidator) bcryptPassword(user *models.Admin) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPWPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) passwordRequired(user *models.Admin) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *models.Admin) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordNotCompileRequirements
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *models.Admin) error {
	if user.HashedPassword == "" {
		return ErrInvalidPassword
	}
	return nil
}

func (uv *userValidator) setRememberTokenIfUnset(user *models.Admin) error {
	if user.Token != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Token = token
	return nil
}

func (uv *userValidator) sessionTokenMinBytes(user *models.Admin) error {
	if user.Token == "" {
		return nil
	}
	n, err := rand.NBytes(user.Token)
	if err != nil {
		return err
	}

	if n < 32 {
		return ErrTokenGenWithError
	}
	return nil
}

func (uv *userValidator) sessionToken(user *models.Admin) error {
	if user.Token == "" {
		return nil
	}
	user.HashedToken = hmac.Hash(user.Token)
	//fmt.Println(user.HashedToken)
	return nil
}

func (uv *userValidator) hashedTokenRequired(user *models.Admin) error {
	if user.HashedToken == "" {
		return ErrTokenRequired
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *models.Admin) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *models.Admin) error {
	if user.Email == "" {
		return ErrEmailIsRequired
	}
	return nil
}

func (uv *userValidator) checkEmailFormat(user *models.Admin) error {
	if user.Email == "" {
		return ErrEmailNotValid
	}
	if !emailRegex.MatchString(user.Email) {
		return ErrEmailNotValid
	}
	return nil
}
