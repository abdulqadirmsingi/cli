// Package cmd contains every CLI command DevPulse exposes.
//
// 🧠 Go Lesson #21: The `cmd` package is the top-level boundary between
// user interaction and business logic. Commands live here; actual work
// lives in internal/. This separation makes code easy to test and reason about.
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// rootCmd is the base command — what runs when you type just `pulse`.
//
// 🧠 Go Lesson #22: Package-level variables are initialized once when the
// program starts. &cobra.Command{...} creates a pointer to a new struct,
// filling in fields by name. Fields not listed get their zero value.
var rootCmd = &cobra.Command{
	Use:   "pulse",
	Short: "track ur grind, no cap 🔥",
	Long:  buildBanner(),
}

// Execute is the public entry point called from main().
//
// 🧠 Go Lesson #23: Exported identifiers start with a CAPITAL letter.
// Execute() is visible to main.go. execute() would be package-private.
// This is Go's entire visibility system — no public/private keywords.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// init() is a special Go function — it runs automatically once the package
// is loaded, before main(). Use it only for wiring up, not complex logic.
//
// 🧠 Go Lesson #24: Every package can have multiple init() functions.
// They all run in the order they appear. Here we use them in each cmd/*.go
// file to register subcommands onto rootCmd — a clean Cobra pattern.
func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// buildBanner renders the ASCII art header with lipgloss colors.
// It's a plain function (not a method) because it needs no receiver.
func buildBanner() string {
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D4FF")).Bold(true)
	purple := lipgloss.NewStyle().Foreground(lipgloss.Color("#9D4EDD"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))

	art := cyan.Render(`
  ██████╗ ███████╗██╗   ██╗██████╗ ██╗   ██╗██╗     ███████╗███████╗
  ██╔══██╗██╔════╝██║   ██║██╔══██╗██║   ██║██║     ██╔════╝██╔════╝
  ██║  ██║█████╗  ██║   ██║██████╔╝██║   ██║██║     ███████╗█████╗
  ██║  ██║██╔══╝  ╚██╗ ██╔╝██╔═══╝ ██║   ██║██║     ╚════██║██╔══╝
  ██████╔╝███████╗ ╚████╔╝ ██║     ╚██████╔╝███████╗███████║███████╗
  ╚═════╝ ╚══════╝  ╚═══╝  ╚═╝      ╚═════╝ ╚══════╝╚══════╝╚══════╝`)

	return art + "\n\n" +
		purple.Render("  track ur grind, no cap 🔥") + "\n" +
		muted.Render("  run `pulse init` to get started\n")
}
