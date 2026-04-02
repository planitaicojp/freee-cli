# Skill Commands Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `freee skill install|update|remove` commands to manage Claude Code skills via git clone/pull/remove.

**Architecture:** Single file `cmd/skill/skill.go` with exported `runInstall/runUpdate/runRemove(baseDir)` functions. Each function takes a base directory for testability. Commands registered via `skill.Cmd` in `cmd/root.go`.

**Tech Stack:** Go, cobra, os/exec (git), internal/errors, internal/prompt

**Spec:** `docs/superpowers/specs/2026-04-02-skill-commands-design.md`

---

## File Structure

| Action | Path | Responsibility |
|--------|------|---------------|
| Create | `cmd/skill/skill.go` | Cmd, constants, runInstall/runUpdate/runRemove, cobra subcommands |
| Create | `cmd/skill/skill_test.go` | Tests for install/update/remove logic |
| Modify | `cmd/root.go` | Register `skill.Cmd` |

---

### Task 1: Write failing tests for install

**Files:**
- Create: `cmd/skill/skill_test.go`

- [ ] **Step 1: Create test file with install tests**

```go
package skill

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInstallCmd(t *testing.T) {
	t.Run("fails when git not found", func(t *testing.T) {
		t.Setenv("PATH", "/nonexistent")
		dir := t.TempDir()

		err := runInstall(dir)
		if err == nil {
			t.Fatal("expected error when git not found")
		}
		if err.Error() != "validation error: git is required to install skills" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("fails when already installed", func(t *testing.T) {
		dir := t.TempDir()
		skillDir := filepath.Join(dir, skillName)
		if err := os.MkdirAll(skillDir, 0o755); err != nil {
			t.Fatal(err)
		}

		err := runInstall(dir)
		if err == nil {
			t.Fatal("expected error when already installed")
		}
		if err.Error() != "validation error: already installed, use 'freee skill update'" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("clones successfully", func(t *testing.T) {
		if _, err := exec.LookPath("git"); err != nil {
			t.Skip("git not available")
		}
		dir := t.TempDir()

		err := runInstall(dir)
		if err != nil {
			t.Skipf("skipping: remote repo not accessible: %v", err)
		}

		skillDir := filepath.Join(dir, skillName)
		if _, err := os.Stat(filepath.Join(skillDir, ".git")); os.IsNotExist(err) {
			t.Error("expected .git directory after install")
		}
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/planitai-freee-cli && go test ./cmd/skill/ -v -run TestInstallCmd`
Expected: FAIL — package does not exist yet

---

### Task 2: Write failing tests for update and remove

**Files:**
- Modify: `cmd/skill/skill_test.go`

- [ ] **Step 1: Add update and remove tests**

Append to `skill_test.go`:

```go
func TestUpdateCmd(t *testing.T) {
	t.Run("fails when not installed", func(t *testing.T) {
		dir := t.TempDir()

		err := runUpdate(dir)
		if err == nil {
			t.Fatal("expected error when not installed")
		}
		if err.Error() != "validation error: not installed, use 'freee skill install'" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("fails when not a git repo", func(t *testing.T) {
		dir := t.TempDir()
		skillDir := filepath.Join(dir, skillName)
		if err := os.MkdirAll(skillDir, 0o755); err != nil {
			t.Fatal(err)
		}

		err := runUpdate(dir)
		if err == nil {
			t.Fatal("expected error when not a git repo")
		}
		if err.Error() != "validation error: not a git repository, remove and reinstall" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("pulls successfully", func(t *testing.T) {
		if _, err := exec.LookPath("git"); err != nil {
			t.Skip("git not available")
		}
		dir := t.TempDir()

		if err := runInstall(dir); err != nil {
			t.Skipf("install failed (remote repo may not exist): %v", err)
		}

		err := runUpdate(dir)
		if err != nil {
			t.Fatalf("update failed: %v", err)
		}
	})
}

func TestRemoveCmd(t *testing.T) {
	t.Run("fails when not installed", func(t *testing.T) {
		dir := t.TempDir()

		err := runRemove(dir)
		if err == nil {
			t.Fatal("expected error when not installed")
		}
		if err.Error() != "validation error: not installed" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("removes successfully with no-input bypassed", func(t *testing.T) {
		dir := t.TempDir()
		skillDir := filepath.Join(dir, skillName)
		if err := os.MkdirAll(skillDir, 0o755); err != nil {
			t.Fatal(err)
		}

		// Call runRemoveForce which skips confirmation
		err := runRemoveForce(dir)
		if err != nil {
			t.Fatalf("remove failed: %v", err)
		}

		if _, err := os.Stat(skillDir); !os.IsNotExist(err) {
			t.Error("expected skill directory to be removed")
		}
	})
}
```

- [ ] **Step 2: Run all tests to verify they fail**

Run: `cd /root/dev/planitai/planitai-freee-cli && go test ./cmd/skill/ -v`
Expected: FAIL — functions not defined

- [ ] **Step 3: Commit test files**

```bash
git add cmd/skill/skill_test.go
git commit -m "test: add failing tests for skill install/update/remove"
```

---

### Task 3: Implement skill.go

**Files:**
- Create: `cmd/skill/skill.go`

- [ ] **Step 1: Write the full implementation**

```go
package skill

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
	"github.com/planitaicojp/freee-cli/internal/prompt"
)

const (
	skillRepo = "https://github.com/planitaicojp/freee-cli-skill.git"
	skillName = "freee-cli-skill"
)

// Cmd is the parent command for skill management.
var Cmd = &cobra.Command{
	Use:   "skill",
	Short: "Manage Claude Code skills for freee-cli",
}

func init() {
	Cmd.AddCommand(installCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(removeCmd)
}

func defaultSkillBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".claude", "skills"), nil
}

func runInstall(baseDir string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return &cerrors.ValidationError{Message: "git is required to install skills"}
	}

	skillDir := filepath.Join(baseDir, skillName)
	if _, err := os.Stat(skillDir); err == nil {
		return &cerrors.ValidationError{Message: "already installed, use 'freee skill update'"}
	}

	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return fmt.Errorf("failed to create skills directory: %w", err)
	}

	cmd := exec.Command("git", "clone", skillRepo, skillDir)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return &cerrors.NetworkError{Err: fmt.Errorf("git clone failed: %w", err)}
	}

	fmt.Fprintln(os.Stderr, "Installed freee-cli-skill successfully.")
	return nil
}

func runUpdate(baseDir string) error {
	skillDir := filepath.Join(baseDir, skillName)
	if _, err := os.Stat(skillDir); os.IsNotExist(err) {
		return &cerrors.ValidationError{Message: "not installed, use 'freee skill install'"}
	}

	gitDir := filepath.Join(skillDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return &cerrors.ValidationError{Message: "not a git repository, remove and reinstall"}
	}

	cmd := exec.Command("git", "-C", skillDir, "pull")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return &cerrors.NetworkError{Err: fmt.Errorf("git pull failed: %w", err)}
	}

	fmt.Fprintln(os.Stderr, "Updated freee-cli-skill successfully.")
	return nil
}

func runRemove(baseDir string) error {
	skillDir := filepath.Join(baseDir, skillName)
	if _, err := os.Stat(skillDir); os.IsNotExist(err) {
		return &cerrors.ValidationError{Message: "not installed"}
	}

	ok, err := prompt.Confirm("Remove freee-cli-skill?")
	if err != nil {
		return err
	}
	if !ok {
		fmt.Fprintln(os.Stderr, "Cancelled.")
		return nil
	}

	return removeSkillDir(baseDir)
}

// runRemoveForce removes without confirmation (for testing).
func runRemoveForce(baseDir string) error {
	skillDir := filepath.Join(baseDir, skillName)
	if _, err := os.Stat(skillDir); os.IsNotExist(err) {
		return &cerrors.ValidationError{Message: "not installed"}
	}

	return removeSkillDir(baseDir)
}

func removeSkillDir(baseDir string) error {
	skillDir := filepath.Join(baseDir, skillName)
	if err := os.RemoveAll(skillDir); err != nil {
		return fmt.Errorf("failed to remove: %w", err)
	}

	fmt.Fprintln(os.Stderr, "Removed freee-cli-skill successfully.")
	return nil
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install freee-cli-skill for Claude Code",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, err := defaultSkillBase()
		if err != nil {
			return err
		}
		return runInstall(base)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update freee-cli-skill to latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, err := defaultSkillBase()
		if err != nil {
			return err
		}
		return runUpdate(base)
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove freee-cli-skill",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, err := defaultSkillBase()
		if err != nil {
			return err
		}
		return runRemove(base)
	},
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `cd /root/dev/planitai/planitai-freee-cli && go test ./cmd/skill/ -v`
Expected: All validation tests PASS, clone/pull tests may SKIP if repo not yet created

- [ ] **Step 3: Commit implementation**

```bash
git add cmd/skill/skill.go
git commit -m "feat: add skill install/update/remove commands"
```

---

### Task 4: Register skill command in root.go

**Files:**
- Modify: `cmd/root.go`

- [ ] **Step 1: Add import and registration**

Add to imports:
```go
"github.com/planitaicojp/freee-cli/cmd/skill"
```

Add to `init()` after `rootCmd.AddCommand(schema.NewCmd(rootCmd))`:
```go
rootCmd.AddCommand(skill.Cmd)
```

- [ ] **Step 2: Verify build succeeds**

Run: `cd /root/dev/planitai/planitai-freee-cli && go build ./...`
Expected: Success, no errors

- [ ] **Step 3: Verify help output shows skill command**

Run: `cd /root/dev/planitai/planitai-freee-cli && go run . skill --help`
Expected: Shows install, update, remove subcommands

- [ ] **Step 4: Run full test suite**

Run: `cd /root/dev/planitai/planitai-freee-cli && go test ./... -count=1`
Expected: All tests pass

- [ ] **Step 5: Commit**

```bash
git add cmd/root.go
git commit -m "feat: register skill command group in root"
```
