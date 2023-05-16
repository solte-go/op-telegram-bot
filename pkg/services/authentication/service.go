package authentication

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/services/hash"
	"telegram-bot/solte.lab/pkg/services/rand"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@` + `[a-z0-9._%+\-]+\.[a-z]{2,16}$`)
	hmac       = hash.NewHMAC(hmacKey)
)

const userPWPepper = "Pinky&Pugovka"
const hmacKey = "Pinky&Pugovka"

type userService struct {
	UserStoreContract
}

type userValidator struct {
	UserStoreContract
	hmac hash.HMAC
}

type UserStoreContract interface {
	CreateUser(user *models.Admin) error
	FindByEmail(email string) (*models.Admin, error)
	UpdateUserData(user *models.Admin) error
	SessionSave(user *models.Admin) error
	FindBySessionToken(hashedToken string) (*models.Admin, error)
	AddNewWordsToDataBase(words []models.Words) error
}

type Service interface {
	UserStoreContract
	Authenticate(email string, pass string) (*models.Admin, error)
	Sanitize(user *models.Admin)
}

func New(store UserStoreContract) Service {
	hmac := hash.NewHMAC(hmacKey)
	uv := &userValidator{
		hmac:              hmac,
		UserStoreContract: store,
	}
	return &userService{
		UserStoreContract: uv,
	}
}

func (u *userService) Authenticate(email string, pswd string) (*models.Admin, error) {

	user, err := u.UserStoreContract.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(pswd+userPWPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}

	return user, nil
}

// Sanitize ...
func (u *userService) Sanitize(user *models.Admin) {
	user.Password = ""
	user.HashedPassword = ""
}

// FindByEmail Searching user by provided email
func (uv *userValidator) FindByEmail(email string) (*models.Admin, error) {
	user := &models.Admin{Email: email}

	if err := runUserValidationFunc(user,
		uv.normalizeEmail,
	); err != nil {
		return nil, err
	}
	return uv.UserStoreContract.FindByEmail(user.Email)
}

// CreateUser Creating new user, function temporary not used
func (uv *userValidator) CreateUser(user *models.Admin) error {
	if err := runUserValidationFunc(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberTokenIfUnset,
		uv.sessionTokenMinBytes,
		uv.sessionToken,
		uv.hashedTokenRequired,
		uv.requireEmail,
		uv.normalizeEmail,
		uv.checkEmailFormat,
		//u.emailIsAvailable,
	); err != nil {
		return err
	}
	err := uv.UserStoreContract.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserData func design for updating user data
func (uv *userValidator) UpdateUserData(user *models.Admin) error {
	if err := runUserValidationFunc(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.sessionTokenMinBytes,
		uv.sessionToken,
		uv.hashedTokenRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.checkEmailFormat,
		//uv.emailIsAvailable,
	); err != nil {
		return err
	}
	return uv.UserStoreContract.UpdateUserData(user)
}

// SessionSave serving session in DB
func (uv *userValidator) SessionSave(user *models.Admin) error {
	if user.Token == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Token = token
	}
	if err := runUserValidationFunc(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.sessionTokenMinBytes,
		uv.sessionToken,
		uv.hashedTokenRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.checkEmailFormat,
	); err != nil {
		return err
	}
	return uv.UserStoreContract.SessionSave(user)
}

// FindBySessionToken Searching user by session token
func (uv *userValidator) FindBySessionToken(token string) (*models.Admin, error) {
	user := models.Admin{
		Token: token,
	}
	if err := runUserValidationFunc(&user,
		uv.sessionToken,
	); err != nil {
		return nil, err
	}
	return uv.UserStoreContract.FindBySessionToken(user.HashedToken)
}
