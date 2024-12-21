package ez_client

import "errors"

var (
	ErrIdentityNotFound = errors.New("identity not found")
	ErrPeerNotFound     = errors.New("peer not found")
	ErrPeerExists       = errors.New("peer already exists")
)
