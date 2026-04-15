package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

type Pkce struct {
	CodeVerifier string
	Challenge string
	State string
}

func GenerateRandomURLEncoded(size int) string {
	// TODO: this may have to be base64 encoding
	randBytes := make([]byte, size) // Create a slice of desired length
	_, err := rand.Read(randBytes) // Fill the slice with random bytes
	if err != nil {
		return "failed" // TODO: test this + check this
	}

	return base64.RawURLEncoding.EncodeToString(randBytes)
}

func GenerateCodeChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))

	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func PrepareOauthFlow() Pkce {
	verifier := GenerateRandomURLEncoded(32)
	challenge := GenerateCodeChallenge(verifier)
	state := GenerateRandomURLEncoded(16)

	return Pkce{CodeVerifier: verifier, Challenge: challenge, State: state}
}

