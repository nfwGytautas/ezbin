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

// Request for downloading a package
type PackageDownloadRequest struct {
	Package string `json:"package"`
	Version string `json:"version"`
}

// Response for downloading a package
type PackageDownloadResponse struct {
	Okay        bool   `json:"okay"`
	PacketCount uint32 `json:"packetCount"`
	FullSize    uint64 `json:"fullSize"`
}

// Request for uploading a package
type PackageUploadRequest struct {
	Package     string `json:"package"`
	Version     string `json:"version"`
	FullSize    uint64 `json:"fullSize"`
	PacketCount uint32 `json:"packetCount"`
}

// Response for uploading a package
type PackageUploadResponse struct {
	Okay bool `json:"okay"`
}
