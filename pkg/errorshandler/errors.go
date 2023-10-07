package errorshandler

const StatusUnauthorized = "Unauthorized"
const StatusBadRequest = "BadRequest"
const StatusNotFound = "NotFound"
const StatusCantReachService = "Internal Service Error"

var (
	ErrAuth                        = &apiErrorHandler{status: StatusUnauthorized, msg: "invalid token"}
	ErrNotFound                    = &apiErrorHandler{status: StatusNotFound, msg: "user not found"}
	ErrBadAccessPrivileges         = &apiErrorHandler{status: StatusUnauthorized, msg: "roles not found"}
	ErrUnauthorized                = &apiErrorHandler{status: StatusUnauthorized, msg: "unauthorized"}
	ErrBadRequest                  = &apiErrorHandler{status: StatusBadRequest, msg: "bad request"}
	ErrDuplicateData               = &apiErrorHandler{status: StatusBadRequest, msg: "duplicate data"}
	ErrDatabaseOperation           = &apiErrorHandler{status: StatusBadRequest, msg: "can't execute query"}
	ErrInvalidUserInfo             = &apiErrorHandler{status: StatusBadRequest, msg: "invalid user data"}
	ErrInvalidUserType             = &apiErrorHandler{status: StatusBadRequest, msg: "invalid user type"}
	ErrInvalidUserInput            = &apiErrorHandler{status: StatusBadRequest, msg: "invalid user input, please check"}
	ErrCantReachAPIEndpoint        = &apiErrorHandler{status: StatusCantReachService, msg: "Internal Service Error"}
	ErrKeyCloakRequestNotSuccessul = &apiErrorHandler{status: StatusCantReachService, msg: "Internal Service Error"}
	ErrInternalServiceError        = &apiErrorHandler{status: StatusCantReachService, msg: "Internal Service Error"}

	//Storage

	ErrStorageInvalidUserType = &apiErrorHandler{status: StatusBadRequest, msg: "invalid user with type freelancer not allowed type"}
)

type APIError interface {
	// APIError returns an HTTP status code and an API-safe error message.
	APIError() (string, string)
}

type apiErrorHandler struct {
	status string
	msg    string
}

func (e apiErrorHandler) Error() string {
	return e.msg
}

func (e apiErrorHandler) APIError() (string, string) {
	return e.status, e.msg
}

type apiWrappedError struct {
	error
	sentinel *apiErrorHandler
}

func (e apiWrappedError) Is(err error) bool {
	return e.sentinel == err
}

func (e apiWrappedError) APIError() (string, string) {
	return e.sentinel.APIError()
}

func WrapError(err error, sentinel *apiErrorHandler) error {
	return apiWrappedError{error: err, sentinel: sentinel}
}
