# ko - Git Worktree tmux Automation

`ko` is a tmux automation script that creates git worktrees and sets up a complete development environment with a single command.

## What it does

`ko` provides two commands for managing git worktrees with tmux:

### `ko new <worktree-name>`

Creates a new development environment:

1. Creates a new git worktree in `.ko/<worktree-name>`
2. Opens a new tmux window with 4 panes arranged in a 2x2 grid:
   - **Top-left**: Opens vim
   - **Bottom-left**: Runs `./bin/setup`
   - **Top-right**: Waits for setup to complete, then runs `./bin/dev`
   - **Bottom-right**: Starts Claude Code CLI

### `ko cleanup <worktree-name>`

Cleans up after you're done:

1. Closes the tmux window for the worktree
2. Removes the git worktree from `.ko/<worktree-name>`

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

### Creating a new worktree session

Navigate to your git repository and run:

```bash
ko new <worktree-name>
```

Example:
```bash
ko new feature-auth
```

This will:
- Create a worktree at `.ko/feature-auth`
- Set up your complete dev environment in tmux
- Run setup first, then automatically start the dev server once setup completes

### Normal development workflow

Once your session is set up:
1. Write code in vim (or exit vim and use your preferred editor)
2. Make commits as normal
3. Push your branch and create a PR when ready
4. Merge to main using your standard git workflow

### Cleaning up after you're done

When your work is merged and you want to clean up:

```bash
ko cleanup feature-auth
```

This will:
- Close the tmux window with all its panes
- Remove the git worktree

**Note:** Make sure you've pushed or merged your changes before running cleanup!

## Commands

```bash
ko new <worktree-name>      # Create a new worktree and tmux session
ko cleanup <worktree-name>  # Close tmux session and remove worktree
ko help                     # Show help message
```

## How it works

### Creating worktrees
The script uses a temporary marker file to coordinate between panes, ensuring `./bin/setup` completes before `./bin/dev` starts. Each invocation creates a unique marker file, allowing multiple worktrees to be created concurrently without conflicts.

### Cleanup
The cleanup command finds the tmux window by name and closes it, then removes the git worktree. If you have uncommitted changes, git will warn you and you'll need to either commit them or use `git worktree remove --force` manually.

## Worktree Management

All worktrees are created in a `.ko/` directory at the root of your repository. This keeps your repository organized and makes it easy to:

- See all active worktrees: `ls .ko/`
- Clean up a worktree: `ko cleanup <name>`
- Manually remove worktrees: `git worktree remove .ko/<name>`
- List all worktrees: `git worktree list`

You may want to add `.ko/` to your `.gitignore` file.

**Tip:** Use `ko cleanup` instead of manually removing worktrees - it will close the tmux window and clean up the worktree in one command!

## Requirements Check

The script includes guard clauses that verify:
- You're in a git repository
- `./bin/setup` exists and is executable
- `./bin/dev` exists and is executable

If any requirement is missing, the script will exit with a helpful error message.

## Why Bash?

This script is written in Bash for several reasons:

**Pros:**
- Universally available on Unix-like systems (macOS, Linux)
- Perfect for orchestrating shell commands (git, tmux)
- Zero dependencies - works out of the box
- Easy to read and modify
- Fast execution for this use case

**Alternatives considered:**
- **Python**: Better error handling and data structures, but requires Python installation and adds complexity for simple shell orchestration
- **Go/Rust**: Excellent for complex tools, but overkill for this use case and requires compilation/distribution
- **Node.js**: Good for cross-platform, but requires Node installation and npm dependencies
- **Fish/Zsh**: Better scripting features but less portable

For this specific use case (automating git and tmux commands), Bash is the sweet spot of portability, simplicity, and functionality.

## Contributing

Feel free to submit issues or pull requests!

## License

MIT
