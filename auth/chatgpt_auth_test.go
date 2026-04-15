package auth_test

import (
	"net/http"
	"testing"

	"github.com/jiasunzhu613/diagnose/auth"
)

func TestRedirectCallbackServer(t *testing.T) {
	expected_code := "testing"
	c := make(chan string)
	go auth.ListenRedirectServer(c)
	go func() {
		http.Get(auth.REDIRECT_URI + "?code=" + expected_code)
	}()

	auth_code := <-c
	if auth_code != expected_code {
		t.Errorf("ListenRedirectServer(c) = %v, want %v", auth_code, expected_code)
	}
}

// TODO: not really a unit test, more like integration test
// func TestStartOauthFlow(t *testing.T) {
// 	result := auth.StartOauthFlow()

// 	log.Println("Received", result)
// }
