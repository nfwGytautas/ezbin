package protocol_test

import (
	"testing"

	"github.com/nfwGytautas/ezbin/ezbin/protocol"
)

func TestNoOpCreate(t *testing.T) {
	no, err := protocol.NewNoOpProtocol()
	if err != nil {
		t.Fatalf("Failed to create noop protocol: %v", err)
	}

	if no.Name() != protocol.PROTOCOL_NOOP {
		t.Fatalf("Protocol name is not correct")
	}

	if no.GetShareableKey() != "" {
		t.Fatalf("Shareable key is not empty")
	}
}

func TestNoOpEncrypt(t *testing.T) {
	message := "Noop encrypt test"
	no, err := protocol.NewNoOpProtocol()
	if err != nil {
		t.Fatalf("Failed to create noop protocol: %v", err)
	}

	encrypted, err := no.Encrypt([]byte(message))
	if err != nil {
		t.Fatalf("Failed to encrypt message: %v", err)
	}

	if string(encrypted) != message {
		t.Fatalf("Message was encrypted")
	}
}

func TestNoOpDecrypt(t *testing.T) {
	message := "Noop decrypt test"
	no, err := protocol.NewNoOpProtocol()
	if err != nil {
		t.Fatalf("Failed to create noop protocol: %v", err)
	}

	encrypted, err := no.Encrypt([]byte(message))
	if err != nil {
		t.Fatalf("Failed to encrypt message: %v", err)
	}

	decrypted, err := no.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt message: %v", err)
	}

	if string(decrypted) != message {
		t.Fatalf("Message was not decrypted")
	}
}
