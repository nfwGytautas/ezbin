package ez_client

import (
	"fmt"

	"github.com/nfwGytautas/ezbin/ezbin"
	"github.com/nfwGytautas/ezbin/ezbin/protocol"
	"github.com/nfwGytautas/ezbin/shared"
)

const IDENTITY_FILE = ".ezbin.identity.json"

// UserIdentity is a struct that represents a user's identity
type UserIdentity struct {
	Version        string                           `json:"version"`
	ProtocolInfo   map[string]protocol.ProtocolData `json:"protocolData"`
	KnownProviders map[string]string                `json:"knownProviders"`
}

// Load local user identity from `.ezbin.identity.json`
func LoadUserIdentity() (*UserIdentity, error) {
	// Check if file exists
	homeDir, err := shared.HomeDirectory()
	if err != nil {
		return nil, err
	}

	fullPath := homeDir + "/" + IDENTITY_FILE

	// Check if exists
	exists, err := shared.FileExists(fullPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrIdentityNotFound
	}

	// Load identity
	var identity UserIdentity

	err = shared.ReadJson(fullPath, &identity)
	if err != nil {
		return nil, err
	}

	// Check version
	if identity.Version != VERSION {
		fmt.Println("⚠️ Identity was generated with an older version of ezbin. Please upgrade your identity.")
	}

	return &identity, nil
}

// Generate new user identity
func GenerateUserIdentity() (*UserIdentity, error) {
	fmt.Println("🔨 Constructing new identity...")

	identity := UserIdentity{
		Version:        VERSION,
		ProtocolInfo:   make(map[string]protocol.ProtocolData),
		KnownProviders: make(map[string]string),
	}

	for _, protocol := range ezbin.GetSupportedProtocols() {
		fmt.Printf("	+ %s (%s)\n", protocol.Name(), protocol.Version())
		data, err := protocol.GenerateNew()
		if err != nil {
			return nil, err
		}

		identity.ProtocolInfo[protocol.Name()] = data
	}

	err := identity.Save()
	if err != nil {
		return nil, err
	}

	return &identity, nil
}

// Save identity to file
func (us *UserIdentity) Save() error {
	homeDir, err := shared.HomeDirectory()
	if err != nil {
		return err
	}

	fullPath := homeDir + "/" + IDENTITY_FILE

	err = shared.WriteJson(fullPath, us)
	if err != nil {
		return err
	}

	return nil
}
