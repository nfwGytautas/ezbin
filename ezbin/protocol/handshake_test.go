package protocol_test

import (
	"testing"

	"github.com/nfwGytautas/ezbin/ezbin/protocol"
)

func TestHandshakeCreate(t *testing.T) {
	hs, err := protocol.NewHandshakeProtocol()
	if err != nil {
		t.Fatalf("Failed to create handshake protocol: %v", err)
	}

	if hs.PrivateKey == "" {
		t.Fatalf("Private key is empty")
	}

	if hs.PublicKey == "" {
		t.Fatalf("Public key is empty")
	}

	if hs.Name() != protocol.PROTOCOL_HANDSHAKE {
		t.Fatalf("Protocol name is not correct")
	}

	if hs.GetShareableKey() != hs.PublicKey {
		t.Fatalf("Shareable key is not the public key")
	}
}

func TestHandshakeEncrypt(t *testing.T) {
	message := "Handshake encrypt test"
	hs, err := protocol.NewHandshakeProtocol()
	if err != nil {
		t.Fatalf("Failed to create handshake protocol: %v", err)
	}

	encrypted, err := hs.Encrypt([]byte(message))
	if err != nil {
		t.Fatalf("Failed to encrypt message: %v", err)
	}

	if string(encrypted) == message {
		t.Fatalf("Message was not encrypted")
	}
}

func TestHandshakeDecrypt(t *testing.T) {
	message := "Handshake decrypt test"
	hs, err := protocol.NewHandshakeProtocol()
	if err != nil {
		t.Fatalf("Failed to create handshake protocol: %v", err)
	}

	encrypted, err := hs.Encrypt([]byte(message))
	if err != nil {
		t.Fatalf("Failed to encrypt message: %v", err)
	}

	decrypted, err := hs.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt message: %v", err)
	}

	if string(decrypted) != message {
		t.Fatalf("Message was not decrypted")
	}
}

func TestHandshakeEncryptDecryptReconstructed(t *testing.T) {
	message := "Handshake decrypt test"
	publicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA59hDoFQIW6KS9AUXdkRiEklzL2IkENNi1Ao3oZ4pvCzH5+amuY5oqG1s4OaYwZ5Qlc89bZ2j/nrllTxMhjImlgv5s+u4l38i5UhcF1wWChvtjEKwWfn7T1M/IHVATeAQ28MTU08TFupwfWVmMPH9JpGqPnt6DejEOtfJfiXT7ki9vRKcGn+qOJ4hJarnJShKINvEI6VKT5zdTTbWLHNngjkJBrSM/qjMhYFAJKpAfc1JaB8VFNC/d/v1LNJhjTDx7VPNAu+cf2T3Bk4qt3uwHCQHAa2TaCgh19uIfZRU9Se3NMEe5GJoPjXspUvuER6RcKMAuNzjMNp3MVydyFYHlQIDAQAB"
	privateKey := "MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDn2EOgVAhbopL0BRd2RGISSXMvYiQQ02LUCjehnim8LMfn5qa5jmiobWzg5pjBnlCVzz1tnaP+euWVPEyGMiaWC/mz67iXfyLlSFwXXBYKG+2MQrBZ+ftPUz8gdUBN4BDbwxNTTxMW6nB9ZWYw8f0mkao+e3oN6MQ618l+JdPuSL29Epwaf6o4niElquclKEog28QjpUpPnN1NNtYsc2eCOQkGtIz+qMyFgUAkqkB9zUloHxUU0L93+/Us0mGNMPHtU80C75x/ZPcGTiq3e7AcJAcBrZNoKCHX24h9lFT1J7c0wR7kYmg+NeylS+4RHpFwowC43OMw2ncxXJ3IVgeVAgMBAAECggEAeAukyIl6YmhFixBv248A8NMTTz+TyRqLG5vGvmp01biiMdNeFMpGKp+uNq1v/yEIPOm3tuBfH89mvOUiAoJJNHwy6RRu2hK8cNgMxxOpXcakM3H8ejpUA/jowNe1Wh1g3Ume4g4Zpk3xvRwZ09IY8DWQXxX0VutlX8qHzEet+ry3HpFfici1erxvdsZDDh55QOa41GTfiX+lS60CspPpaojIGRPtpx9RF2LIpqMSnrRS6rcNadrD2ZDAGljIMh6OtQ8Pfshjeccxj4zSj+uq9yYAoGGUxjsyjKzxvpOgG8h9on2kA9oAVjZ9PwHZZw7eh3/S5ZKTuta5fpuJ6Rqi6QKBgQD3pSjpnBZvskPzqaF/GVKATwkdJkYy/6vXDcHxEt3m6v8CftUSQEsqS98RuXuBOHxSmW7gmiH/FHMzFMH7ozCKqL8MZoR9+mZbirU0h8tpeba0qRxKto28AU0RQM2xLEG7d6Wbi2t9MiWl7ImR/a6XlY/eedlYBlF3FpcE1lnhbwKBgQDvqqQcLWLSXFUfJEVxfQ9S9c7lu5bq57gPj/boB7cO7bmuZvNUdQe8Rul47Zs4oEpbvS+y1df5dvK1m5xSN5/NlVG97S+cGLt4YdqzIQ0m0pS7x19/4zKtHYUMMNxv2x7LY5SANoiIN+dFVmBf11MGaAKYHnkGYU00JgURuIKdOwKBgQDufE43EssEjA8mc0CETtWFnRdwy/AkotVQx/3ydDHgdIRaWdxFtEbul5xdzFsk6UnIndwKTkTZCk+abK4W8GQJ1FIP1hZX37F9DMpOqUt56u3Jc2Y8iStbV4FpURgFPFKc/68raQt9yLI65NzjDAN8FVs0a/Gj9Im1frq2vNpX3wKBgQCUKXbI4Jn+GAybcu3nSfvmOoXMahrTX6rfHA30xYg6l2Y51fVJ2guNLn15P9K8wAMYEa3iLecVlp5W/Ts3bKHDEzN0aaQMKRIESuJL6Pvba0V9jLSSOB+E/AHbVn2APQMdk5MjbBMduwmjSNHNji0KgdRQvE3vTsnOmk559QnyLwKBgQCo81VBRb+c3h668cxEItSU+97crVUJuhRRyVCGzXtcHHGoCnBiB8nij9Bl0/n9sLq17Zqe5lM51BlRU4ZnQNhHXi6hVY20DPP6Ro6mbFu18dbas96rpHzX3d6qV1dcjgZ4IKM4dzcLYfIjUqUdMSk7i3+8SdYa9c0t0Qug5ytVVg=="

	hs := protocol.HandshakeProtocol{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	encrypted, err := hs.Encrypt([]byte(message))
	if err != nil {
		t.Fatalf("Failed to encrypt message: %v", err)
	}

	decrypted, err := hs.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt message: %v", err)
	}

	if string(decrypted) != message {
		t.Fatalf("Message was not decrypted")
	}

}
