package main

import (
	"github.com/1995parham/qtelnet/internal/cmd"
	"github.com/pterm/pterm"
)

func main() {
	header, err := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("q", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("Tel", pterm.NewStyle(pterm.FgLightMagenta)),
		pterm.NewLettersFromStringWithStyle("Net", pterm.NewStyle(pterm.FgLightRed)),
	).Srender()
	if err != nil {
		_ = err
	}

	pterm.DefaultCenter.Println(header)

	cmd.Execute()
}
