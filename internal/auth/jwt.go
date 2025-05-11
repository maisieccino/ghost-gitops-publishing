// internal/auth/jwt.go

package auth

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// FromKey turns Ghost’s “<id>:<secret>” Admin API key
// into a signed HS256 JWT valid for 10 minutes.
func FromKey(adminKey, apiURL string) (string, error) {
	parts := strings.SplitN(adminKey, ":", 2)
	if len(parts) != 2 {
		return "", nil // not a key, probably an already-signed JWT
	}
	id, secretHex := parts[0], parts[1]

	secret, err := hex.DecodeString(secretHex)
	if err != nil {
		return "", err
	}

	iat := time.Now().Unix()
	exp := iat + 600 // 10 minutes

	// Ghost checks the aud against the major version path
	aud := "/v5/admin/"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": iat,
		"exp": exp,
		"aud": aud,
	})
	token.Header["kid"] = id
	return token.SignedString(secret)
}
