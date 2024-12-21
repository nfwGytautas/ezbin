package requests

// Request for getting package info
type PackageInfoRequest struct {
	Package string `json:"package"`
}

// Response for package info
type PackageInfoResponse struct {
	Exists bool  `json:"exists"`
	Size   int64 `json:"size"`
}
