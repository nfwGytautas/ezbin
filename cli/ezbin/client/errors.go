package ez_client

import "errors"

var (
	ErrIdentityNotFound = errIdentityNotFoundProxy()
)

var (
	errIdentityNotFound = errors.New("identity not found")
)

func errIdentityNotFoundProxy() error {
	return errIdentityNotFound
}
