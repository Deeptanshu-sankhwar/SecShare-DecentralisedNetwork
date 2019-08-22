package encryptionproperties

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"net"
)

// GenerateKeyPair - generates a new key pair
func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 2048)

	return privkey, &privkey.PublicKey
}

// PerformHandshake performs handshake with encryption done
func PerformHandshake(conn net.Conn, pub *rsa.PublicKey) *rsa.PublicKey {
	// sending public key
	encoder := json.NewEncoder(conn)
	encoder.Encode(pub)

	// receiving public key
	serverkeys := &rsa.PublicKey{}
	decoder := json.NewDecoder(conn)
	decoder.Decode(&serverkeys)
	return serverkeys
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, _ := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)

	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, _ := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	return plaintext
}
