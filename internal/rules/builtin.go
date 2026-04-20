package rules

import (
	"strings"
	"unicode/utf8"

	"github.com/devpulse-cli/devpulse/internal/git"
)

var mainBranches = map[string]bool{"main": true, "master": true}

// ForceMainRule blocks force-pushing to main or master.
type ForceMainRule struct{}

func (r *ForceMainRule) Name() string { return "force-push-main" }

func (r *ForceMainRule) Evaluate(e *git.Event) *Violation {
	if e.Subcommand != "push" || !e.IsForce {
		return nil
	}
	if !mainBranches[e.Branch] {
		return nil
	}
	return &Violation{
		Severity: SeverityBlock,
		Rule:     r.Name(),
		Message:  "force push to " + e.Branch + " — this rewrites shared history",
		Fix:      "use --force-with-lease if you really must, or open a PR instead",
	}
}

// DirectMainRule warns when committing directly to main or master.
type DirectMainRule struct{}

func (r *DirectMainRule) Name() string { return "direct-main-commit" }

func (r *DirectMainRule) Evaluate(e *git.Event) *Violation {
	if e.Subcommand != "commit" {
		return nil
	}
	if !mainBranches[e.Branch] {
		return nil
	}
	return &Violation{
		Severity: SeverityWarn,
		Rule:     r.Name(),
		Message:  "direct commit to " + e.Branch + " — consider a feature branch",
		Fix:      "try: git checkout -b feat/your-change",
	}
}

// vagueNames are branch names that carry no information about the work.
var vagueNames = map[string]bool{
	"fix": true, "test": true, "update": true, "dev": true,
	"temp": true, "wip": true, "work": true, "stuff": true,
	"changes": true, "patch": true, "hotfix": true, "misc": true,
}

// BranchNameRule warns on vague or unstructured branch names.
type BranchNameRule struct{}

func (r *BranchNameRule) Name() string { return "branch-name" }

func (r *BranchNameRule) Evaluate(e *git.Event) *Violation {
	if e.Subcommand != "checkout" && e.Subcommand != "switch" {
		return nil
	}
	// only applies when creating a branch (-b or -B flag)
	if !hasCreateFlag(e.Args) {
		return nil
	}
	branch := newBranchName(e.Args)
	if branch == "" || mainBranches[branch] {
		return nil
	}
	if vagueNames[strings.ToLower(branch)] {
		return &Violation{
			Severity: SeverityWarn,
			Rule:     r.Name(),
			Message:  "branch name \"" + branch + "\" is too vague",
			Fix:      "try: feat/your-feature, fix/the-bug, chore/what-you-did",
		}
	}
	return nil
}

// vagueMessages are single-word commit messages that mean nothing.
var vagueMessages = map[string]bool{
	"fix": true, "update": true, "wip": true, "changes": true,
	"stuff": true, "misc": true, "test": true, "temp": true,
	"patch": true, "work": true, "commit": true, "save": true,
	"done": true, "edit": true, "refactor": true, "cleanup": true,
}

// VagueCommitRule warns on commit messages that carry no information.
type VagueCommitRule struct{}

func (r *VagueCommitRule) Name() string { return "vague-commit" }

func (r *VagueCommitRule) Evaluate(e *git.Event) *Violation {
	if e.Subcommand != "commit" || e.Message == "" {
		return nil
	}
	msg := strings.TrimSpace(e.Message)
	// too short — under 8 visible characters
	if utf8.RuneCountInString(msg) < 8 {
		return &Violation{
			Severity: SeverityWarn,
			Rule:     r.Name(),
			Message:  "commit message \"" + msg + "\" is too short to be useful",
			Fix:      "describe the why: \"fix: prevent nil panic in auth handler\"",
		}
	}
	// exact match on known vague words (case-insensitive, ignores punctuation)
	clean := strings.ToLower(strings.Trim(msg, " .,!?"))
	if vagueMessages[clean] {
		return &Violation{
			Severity: SeverityWarn,
			Rule:     r.Name(),
			Message:  "commit message \"" + msg + "\" tells future-you nothing",
			Fix:      "describe the why: \"fix: prevent nil panic in auth handler\"",
		}
	}
	return nil
}

func hasCreateFlag(args []string) bool {
	for _, a := range args {
		if a == "-b" || a == "-B" || a == "--orphan" {
			return true
		}
	}
	return false
}

// newBranchName returns the new branch name from checkout/switch args.
func newBranchName(args []string) string {
	skip := false
	for _, a := range args {
		if skip {
			skip = false
			continue
		}
		if a == "-b" || a == "-B" {
			skip = true
			continue
		}
		if strings.HasPrefix(a, "-") {
			continue
		}
		return a
	}
	return ""
}
