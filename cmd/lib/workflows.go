package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/jiasunzhu613/diagnose/api"
	"github.com/jiasunzhu613/diagnose/envconfig"
)

var (
	StdoutBuf bytes.Buffer
	StderrBuf bytes.Buffer
)

var client api.CompletionClient[*api.GenerateRequest, api.GenerateResponse]

func ExecWorkflow(arg []string) error {
	fmt.Println("execCmd workflow")
	cmd := exec.Command(arg[0], arg[1:]...)

	// Multiplex dup the file streams of cmd to current shell and byte buffer
	stdoutMw := io.MultiWriter(os.Stdout, &StdoutBuf)
	stderrMw := io.MultiWriter(os.Stderr, &StderrBuf)

	cmd.Stdout = stdoutMw
	cmd.Stderr = stderrMw

	// We want to catch error here to pass it into LLM
	if err := cmd.Run(); err != nil {
		fmt.Println("======== ERROR OCCURED ========")
		fmt.Println("Error details:", err.Error())
		fmt.Println("Remaining buffer in stdout: ", StdoutBuf.String())
		fmt.Println("Remaining buffer in stderr: ", StderrBuf.String())
		fmt.Println("======== FIX ========")

		client, clientErr := api.NewOllamaClient(envconfig.BASE_OLLAMA_LOCAL)
		if clientErr != nil {
			return clientErr
		}

		// TODO: move this into helper function
		request := &api.GenerateRequest{
			Model: envconfig.OLLAMA_QWEN,
			Prompt: "What is wrong here?\n" + err.Error() + StderrBuf.String() + StdoutBuf.String(),
		}
		
		// response := make([]api.GenerateResponse, 0)
		client.GenerateCompletion(context.Background(), request, func(resp api.GenerateResponse) error {
			fmt.Print(resp.Response)

			return nil
		})
	}

	return nil
}