package ezbin_client

import (
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
	"github.com/nfwGytautas/ezbin/ezbin"
	"github.com/nfwGytautas/ezbin/ezbin/connection"
	"github.com/nfwGytautas/ezbin/ezbin/protocol"
	"github.com/nfwGytautas/ezbin/shared"
)

const IDENTITY_FILE = ".ezbin.identity.json"
const PACKAGE_DIR = ".ezbin"

// Peer info
type PeerInfo struct {
	Address  string `json:"address"`
	Key      string `json:"key"`
	Protocol string `json:"protocol"`
}

// UserIdentity is a struct that represents a user's identity
type UserIdentity struct {
	Version    string              `json:"version"`
	Identifier string              `json:"identifier"`
	PackageDir string              `json:"packageDir"`
	Protocols  protocol.Protocols  `json:"protocolData"`
	Peers      map[string]PeerInfo `json:"Peers"`
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
		return nil, ezbin.ErrIdentityNotFound
	}

	// Load identity
	var identity UserIdentity

	err = shared.ReadJson(fullPath, &identity)
	if err != nil {
		return nil, err
	}

	// Check version
	if identity.Version != ezbin.VERSION {
		fmt.Println("⚠️ Identity was generated with an older version of ezbin. Please upgrade your identity.")
	}

	return &identity, nil
}

// Generate new user identity
func GenerateUserIdentity() (*UserIdentity, error) {
	fmt.Println("🔨 Constructing new identity...")

	identity := UserIdentity{
		Version: ezbin.VERSION,
		Peers:   make(map[string]PeerInfo),
	}

	// Generate identifier
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	identity.Identifier = uuid.String()

	// Generate protocols
	protocols, err := protocol.GenerateProtocols()
	if err != nil {
		return nil, err
	}
	identity.Protocols = protocols

	// Package directory
	homeDir, err := shared.HomeDirectory()
	if err != nil {
		return nil, err
	}

	identity.PackageDir = homeDir + "/" + PACKAGE_DIR + "/" + identity.Identifier

	// Save
	err = identity.Save()
	if err != nil {
		return nil, err
	}

	return &identity, nil
}

// Save identity to file
func (ui *UserIdentity) Save() error {
	homeDir, err := shared.HomeDirectory()
	if err != nil {
		return err
	}

	fullPath := homeDir + "/" + IDENTITY_FILE

	err = shared.WriteJson(fullPath, ui)
	if err != nil {
		return err
	}

	return nil
}

// List peers from identity
func (ui *UserIdentity) ListPeers() {
	if len(ui.Peers) == 0 {
		fmt.Println("⚠️ No peers in identity, try adding one with `ezbin peer add`")
		return
	}

	fmt.Println("Peers:")
	for name, c := range ui.Peers {
		fmt.Printf("	- %s %s %s\n", name, c.Address, c.Protocol)
	}
}

// Set the protocol to use for the peer
func (ui *UserIdentity) SetProtocol(peer string, protocol string) error {
	peerInfo, err := ui.GetPeer(peer)
	if err != nil {
		return err
	}

	peerInfo.Protocol = protocol

	ui.Peers[peer] = peerInfo

	// Save
	err = ui.Save()
	if err != nil {
		return err
	}

	return nil
}

// Check if peer exists
func (ui *UserIdentity) KnowsPeer(name string) bool {
	_, ok := ui.Peers[name]
	return ok
}

// Get peer or return error
func (ui *UserIdentity) GetPeer(name string) (PeerInfo, error) {
	peer, ok := ui.Peers[name]

	if !ok {
		return PeerInfo{}, ezbin.ErrPeerNotFound
	}

	return peer, nil
}

// Add peer to identity
func (ui *UserIdentity) AddPeer(name string, addr string, key string, verify bool) error {
	// Check if peer already exists
	if ui.KnowsPeer(name) {
		return ezbin.ErrPeerExists
	}

	ui.Peers[name] = PeerInfo{
		Address: addr,
		Key:     key,
	}

	// Try and connect to the peer
	if verify {
		fmt.Println("Verifying connection to peer...")

		connection, err := ui.connectToPeer(
			name,
		)
		if err != nil {
			return err
		}
		defer connection.Close()
	}

	// Save
	err := ui.Save()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

// Remove peer from identity
func (ui *UserIdentity) RemovePeer(addr string) error {
	if _, ok := ui.Peers[addr]; !ok {
		return ezbin.ErrPeerNotFound
	}

	delete(ui.Peers, addr)
	err := ui.Save()
	if err != nil {
		return err
	}

	return nil
}

// Check connection to all known peers
func (ui *UserIdentity) CheckPeers() {
	fmt.Println("Checking connections to peers...")

	for peer, _ := range ui.Peers {
		connection, err := ui.connectToPeer(
			peer,
		)
		if err != nil {
			fmt.Printf("	- %s: ❌ %s\n", peer, err)
			continue
		}

		defer connection.Close()

		fmt.Printf("	- %s: ✅\n", peer)
	}
}

func (ui *UserIdentity) GetPackage(pck string, peer string) error {
	fmt.Printf("📦 Getting package: %v\n", pck)

	conn, err := ui.connectToPeer(peer)
	if err != nil {
		return err
	}
	defer conn.Close()

	// TODO: Spanner
	packageInfo := strings.Split(pck, "@")

	// Get package info
	pckInfo, err := conn.GetPackageInfo(packageInfo[0])
	if err != nil {
		return err
	}

	if !pckInfo.Exists {
		return ezbin.ErrPackageNotFound
	}

	packageDir := ui.PackageDir + "/"

	err = conn.DownloadPackage(packageInfo[0], packageInfo[1], packageDir, pckInfo)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Package %v downloaded into: %s\n", pck, packageDir)

	return nil
}

func (ui *UserIdentity) RemovePackage(pck string) error {
	fmt.Printf("📦 Removing package: %v\n", pck)

	outDir := ui.PackageDir + "/"

	// Remove package
	err := shared.DeleteDirectory(outDir + pck)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Package %v removed\n", pck)

	return nil
}

func (ui *UserIdentity) ListPackages() error {
	outDir := ui.PackageDir + "/"

	// List all packages
	packages, err := shared.GetSubdirectories(outDir)
	if err != nil {
		return err
	}

	if len(packages) == 0 {
		fmt.Println("⚠️ No packages found")
	}

	fmt.Println("📦 Packages:")
	for _, pck := range packages {
		if strings.Contains(pck, ".ezbin") {
			continue
		}

		fmt.Println("  + " + pck)

		versions, err := shared.GetSubdirectories(outDir + pck)
		if err != nil {
			return err
		}

		for _, version := range versions {
			fmt.Println("  +--- " + version)
		}
	}

	return nil
}

func (ui *UserIdentity) PublishPackage(pck string, version string, peer string) error {
	fmt.Printf("📦 Publishing package: %v\n", pck)

	conn, err := ui.connectToPeer(peer)
	if err != nil {
		return err
	}
	defer conn.Close()

	currentDir, err := shared.CurrentDirectory()
	if err != nil {
		return err
	}

	pck = strings.ReplaceAll(pck, "/", "")
	packageDir := currentDir + "/" + pck

	fmt.Printf("Creating package from: %s\n", packageDir)

	tmpPath := ui.PackageDir + "/.ezbin/" + pck + "@" + version + ".tar.gz"

	err = shared.TarCompressDirectory(packageDir, tmpPath)
	if err != nil {
		return err
	}

	// Publish package
	err = conn.UploadPackage(pck, version, tmpPath)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Package %v published\n", pck)

	return nil
}

func (ui *UserIdentity) connectToPeer(peer string) (*client2P, error) {
	if !ui.KnowsPeer(peer) {
		return nil, ezbin.ErrPeerNotFound
	}

	c := client2P{}

	// Open a connection to peer
	connData := ui.Peers[peer]

	protocol, err := ui.Protocols.Get(connData.Protocol)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", connData.Address)
	if err != nil {
		return nil, err
	}

	c.conn = conn
	c.frame = connection.NewFrame(conn, make([]byte, connection.HANDSHAKE_BUFFER_SIZE))

	// Handshake with the peer
	err = c.handshake(ui.Identifier, connData.Key, protocol)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &c, nil
}
