package httpErrors

const (
	ErrBadRequest         = "bad request"
	ErrEmailAlreadyExists = "user with given email already exists"
	ErrNoSuchUser         = "user not found"
	ErrWrongCredentials   = "wrong Credentials"
	ErrNotFound           = "not Found"
	ErrUnauthorized       = "unauthorized"
	ErrForbidden          = "forbidden"
	ErrBadQueryParams     = "invalid query params"
)
