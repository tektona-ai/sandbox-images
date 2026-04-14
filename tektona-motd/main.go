package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	version   = "dev"
	buildDate = "unknown"
)

var (
	styleLabel  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8A8A8A"))
	styleValue  = lipgloss.NewStyle().Foreground(lipgloss.Color("#E4E4E4"))
	styleAccent = lipgloss.NewStyle().Foreground(lipgloss.Color("#CDF12E"))
	styleMuted  = lipgloss.NewStyle().Foreground(lipgloss.Color("#5F5F5F"))
	styleRule   = lipgloss.NewStyle().Foreground(lipgloss.Color("#3A3A3A"))
)

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("tektona-motd %s (built %s)\n", version, buildDate)
		return
	}

	info := gather()

	var b strings.Builder
	b.WriteString(renderBanner())
	b.WriteString("\n")
	b.WriteString(renderHeader(info))
	b.WriteString("\n")
	b.WriteString(renderSystem(info))
	b.WriteString("\n")
	b.WriteString(renderResources(info))
	b.WriteString("\n")
	b.WriteString(renderFooter())
	b.WriteString("\n")

	fmt.Fprint(os.Stdout, b.String())
}

func section(title string) string {
	const width = 60
	pad := width - len(title) - 3
	if pad < 0 {
		pad = 0
	}
	return styleAccent.Render("▎ "+title) + " " + styleRule.Render(strings.Repeat("─", pad)) + "\n"
}

func kv(label, value string) string {
	return "  " + styleLabel.Render(fmt.Sprintf("%-12s", label)) + styleValue.Render(value) + "\n"
}
