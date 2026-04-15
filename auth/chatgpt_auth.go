package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/browser"
)

/*
 OpenAI Codex OAuth Flow:
 		=>Build authorization link
 		=> open in browser
		=> listen for callback on localhost
		=> call token url to get access + refresh tokens
*/
// TODO: Add refresh flow

const (
	OPENAI_AUTH_PATH = "https://auth.openai.com/oauth/authorize"
	OPENAI_TOKEN_PATH = "https://auth.openai.com/oauth/token"
	REDIRECT_URI = "http://localhost:1455/auth/callback"
	CLIENT_ID = "app_EMoamEEZ73f0CkXaXp7hrann"
	SCOPE = "openid profile email offline_access";
)

type ChatGptOAuthServerInfo struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn int `json:"expires_in"`
	IDToken string `json:"id_token"`
	TokenType string `json:"token_type"`
}

func OpenAuthorizationLink(pkce Pkce) error{
	authorizationUrl, err := url.Parse(OPENAI_AUTH_PATH)
	if err != nil {
		return err
	}

	data := authorizationUrl.Query()
	data.Set("response_type", "code")
    data.Set("redirect_uri", REDIRECT_URI)
	data.Set("client_id", CLIENT_ID)
	data.Set("scope", SCOPE)
	data.Set("code_challenge", pkce.Challenge)
	data.Set("code_challenge_method", "S256")
	data.Set("state", pkce.State);
	data.Set("id_token_add_organizations", "true");
	data.Set("codex_cli_simplified_flow", "true");
	data.Set("originator", "codex_cli_rs");

	authorizationUrl.RawQuery = data.Encode()
	log.Println("authorizationUrl", authorizationUrl.String())
	if err := browser.OpenURL(authorizationUrl.String()); err != nil {
		return err
	}

	return nil
}

// TODO: add timeout to this?
// TODO: change this to longstanding server for when user enters "diagnose login"
func ListenRedirectServer(c chan string) {
	log.Println("Starting redirect callback server...")

	var auth_code string
	server := &http.Server{Addr: ":1455"}
	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, req *http.Request) {
		query_params := req.URL.Query()
		log.Println("Extracted query params: ", query_params)

		auth_code = query_params.Get("code")

		go server.Close()
	})

	_ = server.ListenAndServe() // Best effort (should probably do some signal handling to handle graceful shutdown)

	c <- auth_code
}

func exchangeAuthorizationCode(code, verifier string) (*ChatGptOAuthServerInfo, error) {
	data := url.Values{
		"grant_type": {"authorization_code"},
		"code_verifier": {verifier},
		"redirect_uri": {REDIRECT_URI},
		"client_id": {CLIENT_ID},
		"code": {code},
	}

	resp, err := http.Post(
		OPENAI_TOKEN_PATH, 
		"application/x-www-form-urlencoded", 
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		log.Println("Returned because http post was bad")
		return nil, fmt.Errorf("token exchange failed: bad http post request, %s", resp.Status)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errResp.ErrorDescription != "" {
			return nil, fmt.Errorf("token exchange failed: %s - %s", errResp.Error, errResp.ErrorDescription)
		}
		return nil, fmt.Errorf("token exchange failed: %s", resp.Status)
	}

	var oAuthFlowResponse ChatGptOAuthServerInfo
	if err := json.NewDecoder(resp.Body).Decode(&oAuthFlowResponse); err != nil {
		log.Println("Returned because decoding was bad")
		return &oAuthFlowResponse, nil
	}

	return &oAuthFlowResponse, nil
}

func StartOauthFlow() *ChatGptOAuthServerInfo{
	pcke := PrepareOauthFlow()

	// Start redirect server to read from later
	c := make(chan string)
	go ListenRedirectServer(c)

	err := OpenAuthorizationLink(pcke)
	if err != nil {
		return nil
	}

	auth_code := <-c

	oAuthResponse, err := exchangeAuthorizationCode(auth_code, pcke.CodeVerifier)
	if err != nil {
		log.Println(err)
	}
	return oAuthResponse
}
