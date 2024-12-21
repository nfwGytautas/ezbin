package requests

import (
	"strings"
)

// Convert a header to a request
func HeaderToRequest(header string) any {
	switch strings.TrimRight(header, "\x00") {
	case HeaderHandshake:
		return &HandshakeRequest{}
	case HeaderPackageInfo:
		return &PackageInfoRequest{}
	}

	return nil
}

// Convert a header to a response
func HeaderToResponse(header string) any {
	switch strings.TrimRight(header, "\x00") {
	case HeaderHandshake:
		return &HandshakeRequest{}
	case HeaderPackageInfo:
		return &PackageInfoResponse{}
	}

	return nil
}
