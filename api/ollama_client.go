package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
)

// Client encapsulates client state for interacting with the ollama
// service. Use [ClientFromEnvironment] to create new Clients.
type OllamaClient struct {
	base *url.URL // API base path
	http *http.Client
}

func NewOllamaClient(base string) (*OllamaClient, error) {
	url, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	
	return &OllamaClient{base:url, http:&http.Client{}}, nil
}

// Performs an HTTP request
// Allows any req and response type to be performed 
// ctx is used to create http request with context
func (c *OllamaClient) do(ctx context.Context, method, endpoint string, reqData, respData any) error {
	var reqBody io.Reader
	var err error

	switch reqData := reqData.(type) {
	case io.Reader:
		reqBody = reqData
	case nil:
		// noop
	default:
		// A request body struct came through
		body, err := json.Marshal(reqData)
		if err != nil {
			return err
		}

		reqBody = bytes.NewReader(body)
	}
	
	reqUrl := c.base.JoinPath(endpoint)

	request, err := http.NewRequestWithContext(ctx, method, reqUrl.String(), reqBody)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", fmt.Sprintf("diagnose-go (%s, %s) Go/%s", runtime.GOARCH, runtime.GOOS, runtime.Version()))
	
	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if len(respBody) > 0 && respData != nil {
		if err := json.Unmarshal(respBody, respData); err != nil {
			return err
		}
	}

	return nil
}

const MEGABYTE = 1_000_000
const MAX_BUFFER_SIZE = 8 * MEGABYTE
// Takes application/x-ndjson
// This is streaming JSON which maintains long lived connection between client and server and allows the server
// to send JSON objects as it processes them and allows the client to deliver objects as they come 
func (c *OllamaClient) stream(ctx context.Context, method, endpoint string, reqData any, fn func([]byte) error) error {
	var reqBody io.Reader
	var err error

	switch reqData := reqData.(type) {
	case io.Reader:
		reqBody = reqData
	case nil:
		// noop
	default:
		// A request body struct came through
		body, err := json.Marshal(reqData)
		if err != nil {
			return err
		}

		reqBody = bytes.NewReader(body)
	}
	
	reqUrl := c.base.JoinPath(endpoint)

	request, err := http.NewRequestWithContext(ctx, method, reqUrl.String(), reqBody)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/x-ndjson") // Newline delimited streaming JSON
	request.Header.Set("User-Agent", fmt.Sprintf("diagnose-go (%s, %s) Go/%s", runtime.GOARCH, runtime.GOOS, runtime.Version()))
	
	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	buffArr := make([]byte, 0, MAX_BUFFER_SIZE) // Indicate capacity of MAX_BUFFER_SIZE

	scanner.Buffer(buffArr, MAX_BUFFER_SIZE)
	for scanner.Scan() {
		var errorResponse struct {
			Error     string `json:"error,omitempty"`
		}

		bytes := scanner.Bytes()
		if err := json.Unmarshal(bytes, &errorResponse); err != nil {
			return err
		}

		if errorResponse.Error != "" {
			return errors.New(errorResponse.Error)
		}

		if err := fn(bytes); err != nil {
			return err
		}
	}

	return nil
}

type GenerateResponseFunction func(GenerateResponse) error
func (c *OllamaClient) GenerateCompletion(ctx context.Context, req *GenerateRequest, fn func(GenerateResponse) error) error {
	return c.stream(ctx, http.MethodPost, "api/generate", *req, func(bytes []byte) error {
		var response GenerateResponse

		if err := json.Unmarshal(bytes, &response); err != nil {
			return err
		}

		return fn(response)
	})
}