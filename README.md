# ko - Git Worktree tmux Automation

`ko` is a tmux automation script that creates git worktrees and sets up a complete development environment with a single command.

## What it does

When you run `ko <worktree-name>`, it will:

1. Create a new git worktree in `.ko/<worktree-name>`
2. Open a new tmux window with 4 panes arranged in a 2x2 grid:
   - **Top-left**: Opens vim
   - **Bottom-left**: Runs `./bin/setup`
   - **Top-right**: Waits for setup to complete, then runs `./bin/dev`
   - **Bottom-right**: Starts Claude Code CLI

## Prerequisites

- Git repository with worktree support
- tmux installed
- `./bin/setup` script in your repository
- `./bin/dev` script in your repository
- Claude Code CLI (optional - will just fail gracefully if not installed)

## Installation

1. Clone or download this repository
2. Add the script to your PATH or create an alias:

```bash
# Option 1: Add to PATH
export PATH="$HOME/code/ko:$PATH"

# Option 2: Create an alias
alias ko="$HOME/code/ko/ko.sh"
```

Add the above line to your `~/.bashrc`, `~/.zshrc`, or equivalent shell configuration file.

## Usage

Navigate to your git repository and run:

```bash
ko <worktree-name>
```

Example:
```bash
ko feature-auth
```

This will:
- Create a worktree at `.ko/feature-auth`
- Set up your complete dev environment in tmux
- Run setup first, then automatically start the dev server once setup completes

## How it works

The script uses a temporary marker file to coordinate between panes, ensuring `./bin/setup` completes before `./bin/dev` starts. Each invocation creates a unique marker file, allowing multiple worktrees to be created concurrently without conflicts.

## Worktree Management

All worktrees are created in a `.ko/` directory at the root of your repository. This keeps your repository organized and makes it easy to:

- See all active worktrees: `ls .ko/`
- Remove old worktrees: `git worktree remove .ko/<name>`
- List worktrees: `git worktree list`

You may want to add `.ko/` to your `.gitignore` file.

## Requirements Check

The script includes guard clauses that verify:
- You're in a git repository
- `./bin/setup` exists and is executable
- `./bin/dev` exists and is executable

If any requirement is missing, the script will exit with a helpful error message.

## License

MIT
