package handler

import (
	"bufio"
	"context"
	"os"

	"github.com/quic-go/quic-go"
	"github.com/pterm/pterm"
)

// Prompt reads user input from stdin and sends it to the server
// via QUIC streams. This function is blocking and should be called
// from the main goroutine.
func Prompt(ctx context.Context, conn *quic.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		stream, err := conn.OpenStreamSync(ctx)
		if err != nil {
			pterm.Error.Printf("cannot open stream: %s\n", err)

			continue
		}

		if _, err := stream.Write([]byte(line + "\n")); err != nil {
			pterm.Error.Printf("cannot write to stream: %s\n", err)
		}

		if err := stream.Close(); err != nil {
			pterm.Error.Printf("cannot close stream: %s\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		pterm.Fatal.Printf("error reading input: %s\n", err)
	}
}
