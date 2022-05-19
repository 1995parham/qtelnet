package handler

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/lucas-clemente/quic-go"
	"github.com/pterm/pterm"
)

// Accepter accept input streams and then print
// their data with their id. please note that this function is a blocking function
// and should run it in a goroutine.
func Accepter(ctx context.Context, conn quic.Connection) {
	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			pterm.Fatal.Printf("cannot accept stream %s\n", err)
		}

		go func() {
			pp := pterm.PrefixPrinter{
				Prefix: pterm.Prefix{
					Text:  fmt.Sprintf("%d", stream.StreamID().StreamNum()),
					Style: &pterm.Style{pterm.BgGray, pterm.FgLightMagenta},
				},
				Scope: pterm.Scope{
					Text:  "",
					Style: nil,
				},
				MessageStyle:     nil,
				Fatal:            false,
				ShowLineNumber:   false,
				LineNumberOffset: 0,
				Writer:           os.Stdout,
				Debugger:         false,
			}

			if _, err := io.Copy(pp.Writer, stream); err != nil {
				pterm.Fatal.Printf("stream %d failed", stream.StreamID().StreamNum())
			}
		}()
	}
}
