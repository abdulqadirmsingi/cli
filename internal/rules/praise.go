package rules

import (
	"math/rand"
	"strings"

	"github.com/devpulse-cli/devpulse/internal/git"
)

var commitPraises = []string{
	"pulse approves this commit message 🤌",
	"your teammates just smiled reading that commit 😊",
	"pulse has seen 10,000 commits — this one passes 🎯",
	"somewhere a senior dev shed a single tear of joy",
	"future-you will actually know what happened here — rare",
	"your code reviewer's prayers have been answered",
	"git log is going to look so clean because of this 📖",
	"clean, clear, conventional — the holy trinity ✅",
	"code archaeologists in 2040 will understand this. respect.",
	"pulse clocked the conventional format. keep going 📊",
	"no 'wip wip wip' detected — you're evolving 📈",
	"that commit message could honestly be in a textbook fr",
	"pulse is silently judging everyone else's commits rn",
	"the kind of message that makes PRs a joy to review 🚀",
	"your future self just sent you a thank you note",
}

var branchPraises = []string{
	"pulse can tell what this branch does without reading a single line of code 👀",
	"anyone opening this PR will immediately know what's going on",
	"clean branch, clean mind — pulse vibes with this 🧠",
	"that name tells a whole story 📖",
	"feat/fix/chore discipline — pulse respects it 🫡",
	"a branch name so good it barely needs a PR description",
	"your branch naming is so clean it's suspicious",
	"this is how PRs stay organized and reviewers stay sane",
	"future-you won't spend 10 minutes wondering what this was",
}

var pushPraises = []string{
	"keeping main clean — pulse clocked that 👀",
	"main is sacred and you already know it 🙏",
	"feature branch push, exactly how it's supposed to go ✅",
	"somewhere a CI pipeline is about to be very happy",
	"PR flow intact, team flow intact — pulse approves 🤝",
	"main stays untouched. pulse respects that discipline 🔥",
	"clean push to a feature branch — this is the way 🚀",
}

var goodBranchPrefixes = []string{
	"feat/", "fix/", "chore/", "docs/", "refactor/",
	"test/", "perf/", "ci/", "build/", "hotfix/",
}

// GoodCommitPraise fires when a commit message follows the conventional format.
type GoodCommitPraise struct{}

func (r *GoodCommitPraise) Name() string { return "good-commit" }

func (r *GoodCommitPraise) Evaluate(e *git.Event) *Praise {
	if e.Subcommand != "commit" || e.Message == "" {
		return nil
	}
	lower := strings.ToLower(strings.TrimSpace(e.Message))
	for _, p := range conventionalPrefixes {
		if strings.HasPrefix(lower, p) {
			return &Praise{
				Rule:    r.Name(),
				Message: commitPraises[rand.Intn(len(commitPraises))],
			}
		}
	}
	return nil
}

// GoodBranchPraise fires when a new branch follows the feat/fix/chore naming convention.
type GoodBranchPraise struct{}

func (r *GoodBranchPraise) Name() string { return "good-branch" }

func (r *GoodBranchPraise) Evaluate(e *git.Event) *Praise {
	if e.Subcommand != "checkout" && e.Subcommand != "switch" {
		return nil
	}
	if !hasCreateFlag(e.Args) {
		return nil
	}
	branch := strings.ToLower(newBranchName(e.Args))
	for _, p := range goodBranchPrefixes {
		if strings.HasPrefix(branch, p) {
			return &Praise{
				Rule:    r.Name(),
				Message: branchPraises[rand.Intn(len(branchPraises))],
			}
		}
	}
	return nil
}

// GoodPushPraise fires when pushing to a non-main branch (keeping main clean).
type GoodPushPraise struct{}

func (r *GoodPushPraise) Name() string { return "good-push" }

func (r *GoodPushPraise) Evaluate(e *git.Event) *Praise {
	if e.Subcommand != "push" || e.IsForce {
		return nil
	}
	target := e.PushTarget
	if target == "" {
		target = e.Branch
	}
	if mainBranches[target] {
		return nil
	}
	return &Praise{
		Rule:    r.Name(),
		Message: pushPraises[rand.Intn(len(pushPraises))],
	}
}
