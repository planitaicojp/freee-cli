# freee skill install|update|remove Design

## Overview

Add `freee skill` command group to manage Claude Code skills for freee-cli. Follows the same pattern as `conoha-cli`'s skill commands — single repository, git-based install/update/remove.

## Commands

```
freee skill install   — clone skill repo to ~/.claude/skills/freee-cli-skill/
freee skill update    — pull latest changes in skill directory
freee skill remove    — confirm and delete skill directory
```

## Constants

| Item | Value |
|------|-------|
| Repository URL | `https://github.com/planitaicojp/freee-cli-skill.git` |
| Skill name | `freee-cli-skill` |
| Install path | `~/.claude/skills/freee-cli-skill/` |

## Behavior

### install

1. Check `git` is on PATH → `ValidationError` if missing
2. Check skill directory does not exist → `ValidationError` if already installed (hint: use `freee skill update`)
3. `git clone <repo> <skillDir>` (stdout/stderr → os.Stderr)
4. Print success message to stderr

### update

1. Check skill directory exists → `ValidationError` if not installed (hint: use `freee skill install`)
2. Check `.git` subdirectory exists → `ValidationError` if not a git repo (hint: remove and reinstall)
3. `git -C <skillDir> pull` (stdout/stderr → os.Stderr)
4. Print success message to stderr

### remove

1. Check skill directory exists → `ValidationError` if not installed
2. `prompt.Confirm("Remove freee-cli-skill?")` — respects `--no-input` (returns error if set)
3. `os.RemoveAll(skillDir)`
4. Print success message to stderr

## Error Handling

| Condition | Error Type | Hint |
|-----------|-----------|------|
| git not installed | `ValidationError` | install git |
| Already installed | `ValidationError` | use `freee skill update` |
| Not installed | `ValidationError` | use `freee skill install` |
| Not a git repo | `ValidationError` | remove and reinstall |
| git clone/pull fails | `NetworkError` | check internet connection |

## File Structure

```
cmd/skill/
  skill.go        — Cmd, init(), constants, runInstall/runUpdate/runRemove, cobra commands
  skill_test.go   — temp directory based tests
```

Registration: add `skill.Cmd` to `cmd/root.go`.

## Testing Strategy

- Export `runInstall(baseDir)`, `runUpdate(baseDir)`, `runRemove(baseDir)` for testability
- Tests use `t.TempDir()` as baseDir
- install: verify directory created with `.git`
- update: create a git repo first, then verify pull succeeds
- remove: create directory, verify deletion (skip Confirm in tests via `--no-input` handling or direct function call)

## Reference

- conoha-cli skill implementation: `github.com/crowdy/conoha-cli` branch `feature/skill-commands`, `cmd/skill/skill.go`
