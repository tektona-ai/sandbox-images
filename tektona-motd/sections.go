package main

import (
	"fmt"
	"strings"
)

func renderHeader(i Info) string {
	var b strings.Builder
	greeting := fmt.Sprintf("Welcome, %s", i.User)
	b.WriteString("  " + styleAccent.Render(greeting))
	b.WriteString("\n")

	if i.SandboxID != "" {
		b.WriteString("  " + styleMuted.Render("sandbox ") + styleValue.Render(i.SandboxID) + "\n")
	}
	return b.String()
}

func renderSystem(i Info) string {
	var b strings.Builder
	b.WriteString(section("System"))
	b.WriteString(kv("Host", i.Hostname))
	b.WriteString(kv("OS", i.OS))
	b.WriteString(kv("Kernel", i.Kernel))
	b.WriteString(kv("Uptime", formatUptime(i.Uptime)))
	return b.String()
}

func renderResources(i Info) string {
	var b strings.Builder
	b.WriteString(section("Resources"))

	memUsedKB := uint64(0)
	if i.MemTotalKB > i.MemAvailKB {
		memUsedKB = i.MemTotalKB - i.MemAvailKB
	}
	memDetail := fmt.Sprintf("%s / %s",
		formatBytes(memUsedKB*1024), formatBytes(i.MemTotalKB*1024))
	b.WriteString("  " + styleLabel.Render(fmt.Sprintf("%-12s", "Memory")) +
		renderBar(memUsedKB, i.MemTotalKB, memDetail) + "\n")

	diskUsed := uint64(0)
	if i.DiskTotalB > i.DiskFreeB {
		diskUsed = i.DiskTotalB - i.DiskFreeB
	}
	diskDetail := fmt.Sprintf("%s / %s on /",
		formatBytes(diskUsed), formatBytes(i.DiskTotalB))
	b.WriteString("  " + styleLabel.Render(fmt.Sprintf("%-12s", "Disk")) +
		renderBar(diskUsed, i.DiskTotalB, diskDetail) + "\n")

	b.WriteString("  " + styleLabel.Render(fmt.Sprintf("%-12s", "Load")) +
		renderLoadBar(i.Load1, i.NCPU) + "\n")

	return b.String()
}

func renderFooter() string {
	var b strings.Builder
	b.WriteString(section("Quick reference"))
	b.WriteString("  " + styleMuted.Render("docs:    ") +
		styleValue.Render("https://docs.tektona.ai") + "\n")
	b.WriteString("  " + styleMuted.Render("support: ") +
		styleValue.Render("https://github.com/tektona-ai/sandbox-images/issues") + "\n")
	b.WriteString("  " + styleMuted.Render("CLI:     ") +
		styleValue.Render("tektonactl --help") + "\n")
	return b.String()
}
