## Merge from Main (Go GitHub Action)

Automates running commands on `main`, committing results, opening a PR, waiting for CI, and squash-merging when checks succeed. Designed to be re-entrant: runs triggered by its own commits are skipped via prefixes.

### How it works
- `ConfirmShouldRun`: skip if the last commit message starts with `Auto Merge`, `[Auto Merge]:`, the provided `commit_prefix`, or any extra prefixes in `PREFIXES_TO_IGNORE`.
- `RunCommands`: executes the supplied commands (newline or comma separated).
- `CommitAndOpenPR`: creates a branch, commits, pushes, and opens a PR using the prefix for the title and commit.
- `Wait`: pauses 30 seconds before checking CI.
- `WaitForCI`: polls combined status until success/failure (15m timeout).
- `Merge`: squash-merges the PR using the prefix.

### Inputs
- `github_access_token` (required): token with push and PR/merge rights.
- `commands` (required): newline-separated commands (e.g. `go run ./...`).
- `commit_prefix` (optional): commit/PR prefix, defaults to `[Auto Merge]`.

### Environment
- `PREFIXES_TO_IGNORE`: optional comma-delimited prefixes to skip reruns. Empty string is ignored.

### Usage
```yaml
name: Merge from Main
on:
  push:
    branches: [main]

jobs:
  merge-from-main:
    runs-on: ubuntu-latest
    steps:
      - uses: keithweaver/github-action-merge-from-main@v1
        with:
          github_access_token: ${{ secrets.GITHUB_TOKEN }}
          commands: |
            go run ./
          commit_prefix: "[Auto Merge]"
        env:
          PREFIXES_TO_IGNORE: "[Skip Me],[Do Not Run]"
```
