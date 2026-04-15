package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

var (
	stdinBuf bytes.Buffer
	stdoutBuf bytes.Buffer
)

func test() error {
        // Create arbitrary command.
        c := exec.Command("bash")

        // Start the command with a pty.
        ptmx, err := pty.Start(c)
        if err != nil {
                return err
        }

		stdinMw := io.MultiWriter(ptmx, &stdinBuf)
		stdoutMw := io.MultiReader(ptmx, &stdoutBuf)

		// Set up reader for ptmx???
		// ptmxReader := bufio.NewScanner(ptmx)
		
        // Make sure to close the pty at the end.
        defer func() { _ = ptmx.Close() }() // Best effort.

        // Handle pty size.
        ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGWINCH)
        go func() {
                for range ch {
                        if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
                                log.Printf("error resizing pty: %s", err)
                        }
                }
        }()
        ch <- syscall.SIGWINCH // Initial resize.
        defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

        // Set stdin in raw mode.
		// NOTE: essentially disallow original shell instance to be able to be interacted with (? i think)
        oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
        if err != nil {
                panic(err)
        }
        defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

        // Copy stdin to the pty and the pty to stdout.
        // NOTE: The goroutine will keep reading until the next keystroke before returning.
        go func() { _, _ = io.Copy(stdinMw, os.Stdin) }()
		// go func() {
		// 	for {
		// 		time.Sleep(10 * time.Second)
		// 		fmt.Println("This is stdin: ", stdinBuf.String())
		// 		fmt.Println("This is stdout: ", stdoutBuf.String())
		// 	}
		// }()
        _, _ = io.Copy(os.Stdout, stdoutMw)

        return nil
}

func main() {
        if err := test(); err != nil {
                log.Fatal(err)
        }
}