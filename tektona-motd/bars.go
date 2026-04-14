package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const barWidth = 32

func renderBar(used, total uint64, detail string) string {
	if total == 0 {
		return styleMuted.Render("  (unavailable)")
	}
	ratio := float64(used) / float64(total)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	filled := int(ratio * float64(barWidth))
	if filled == 0 && used > 0 {
		filled = 1
	}
	return paintBar(ratio, filled) + "  " +
		styleValue.Render(fmt.Sprintf("%4.1f%%", ratio*100)) + "  " +
		styleMuted.Render(detail)
}

func renderLoadBar(load float64, ncpu int) string {
	if ncpu <= 0 {
		return styleMuted.Render("  (unavailable)")
	}
	ratio := load / float64(ncpu)
	if ratio > 1 {
		ratio = 1
	}
	filled := int(ratio * float64(barWidth))
	return paintBar(ratio, filled) + "  " +
		styleValue.Render(fmt.Sprintf("%.2f", load)) + "  " +
		styleMuted.Render(fmt.Sprintf("load1 (%d CPUs)", ncpu))
}

// paintBar renders the filled/empty segments. Filled uses a medium-density
// block (▒) for a dotted / translucent look; empty uses a light block (░).
func paintBar(ratio float64, filled int) string {
	color := lipgloss.Color("#CDF12E")
	switch {
	case ratio >= 0.85:
		color = lipgloss.Color("#FF5F5F")
	case ratio >= 0.60:
		color = lipgloss.Color("#E6B800")
	}
	return lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("⣿", filled)) +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#3A3A3A")).Render(strings.Repeat("⣀", barWidth-filled))
}
