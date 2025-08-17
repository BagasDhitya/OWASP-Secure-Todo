package security

import (
	"crypto/rand"
	"encoding/hex"
)

// RandJTI generates a random 16-byte hex string to use as JWT ID.
func RandJTI() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// fallback, meskipun kecil kemungkinannya gagal
		return "fallback-jti"
	}
	return hex.EncodeToString(b)
}
