package ez_client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nfwGytautas/ezbin/ezbin/connection"
	"github.com/nfwGytautas/ezbin/shared"
)

func GetPackage(i *UserIdentity, pck string, peer string) error {
	fmt.Printf("📦 Getting package: %v\n", pck)

	if !i.KnowsPeer(peer) {
		return ErrPeerNotFound
	}

	// Open a connection to peer
	connData := i.Peers[peer]

	conn, err := connection.ConnectC2P(connection.C2PConnectionParameters{
		Peer:           connData,
		UserIdentifier: i.Identifier,
	})
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
		return errors.New("package not found")
	}

	packageDir := i.PackageDir + "/"

	err = conn.DownloadPackage(packageInfo[0], packageInfo[1], packageDir, pckInfo)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Package %v downloaded into: %s\n", pck, packageDir)

	return nil
}

func RemovePackage(i *UserIdentity, pck string) error {
	fmt.Printf("📦 Removing package: %v\n", pck)

	outDir := i.PackageDir + "/"

	// Remove package
	err := shared.DeleteDirectory(outDir + pck)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Package %v removed\n", pck)

	return nil
}

func ListPackages(i *UserIdentity) error {
	outDir := i.PackageDir + "/"

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
