package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

func MakeCSRFCookie(secret, sessionID string, ttl time.Duration) (value string, expires time.Time) {
	expires = time.Now().Add(ttl)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(sessionID))
	sum := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(sum), expires
}

func ValidateCSRF(secret, sessionID, token string) bool {
	want, _ := MakeCSRFCookie(secret, sessionID, 0)
	return hmac.Equal([]byte(want), []byte(token))
}
