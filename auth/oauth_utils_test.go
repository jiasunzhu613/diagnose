package auth_test

import (
	"encoding/base64"
	"testing"

	"github.com/jiasunzhu613/diagnose/auth"
)

// func TestGenerateCodeChallenge(t *testing.T) {
// 	verifier := "dBjftJeZ4CVP-mB92K27uhbUbP1E6as7mt35ghasabc"
//     expected := base64.RawURLEncoding.EncodeToString([]byte("E9Melhoa2OqFrTRua8_Y679I4YmSggmbe65Uq56M-7Y"))

// 	if challenge := auth.GenerateCodeChallenge(verifier); challenge != expected {
// 		t.Errorf("GenerateCodeChallenge(%v) = %v, want %v", verifier, challenge, expected)
// 	}
// }

func TestGenerateRandomURLEncoded(t *testing.T) {
	size := 32
	urlEncoded := auth.GenerateRandomURLEncoded(size)
	decoded, err := base64.RawURLEncoding.DecodeString(urlEncoded)
	if err != nil {
		t.Error("problem decoding base64 url encoded string")
	}

	if len(decoded) != size {
		t.Errorf("GenerateRandomURLEncoded(%v) = %v, decoded = %v, got length: %v, want %v", 
			size, urlEncoded, decoded, len(decoded), size)
	}
}
