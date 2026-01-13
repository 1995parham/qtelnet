package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/1995parham/qtelnet/internal/handler"
	"github.com/quic-go/quic-go"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func main(host, port string, insecure bool, alpn string) {
	addr := fmt.Sprintf("%s:%s", host, port)

	// nolint: exhaustruct, gosec
	tlsConf := &tls.Config{
		InsecureSkipVerify: insecure,
		NextProtos:         []string{alpn},
	}

	ctx, cancel := context.WithCancel(context.Background())

	conn, err := quic.DialAddr(ctx, addr, tlsConf, nil)
	if err != nil {
		pterm.Fatal.Printf("cannot dial to server %s\n", err)
	}
	pterm.Success.Printf("connected\n")
	defer cancel()

	go handler.Accepter(ctx, conn)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		cancel()
	}()

	handler.Prompt(ctx, conn)
}

const (
	// ExitFailure status code.
	ExitFailure = 1

	// NArgs is the exact number of arguments.
	NArgs = 2
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var insecure bool
	var alpn string

	// nolint: exhaustruct
	root := &cobra.Command{
		Use:     "qtelnet address port",
		Short:   "Telnet like application based on quic protocol",
		Example: "qtelnet 127.0.0.1 8080",
		Args:    cobra.ExactArgs(NArgs),
		Run: func(cmd *cobra.Command, args []string) {
			main(args[0], args[1], insecure, alpn)
		},
	}

	root.Flags().BoolVarP(&insecure, "insecure", "K", true, "skip certificate verification")
	root.Flags().StringVarP(&alpn, "alpn", "a", "quic-echo", "ALPN protocol to use")

	if err := root.Execute(); err != nil {
		os.Exit(ExitFailure)
	}
}
