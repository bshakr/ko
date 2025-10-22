# ko - Git Worktree tmux Automation

`ko` is a CLI tool written in Go that creates git worktrees and sets up a complete development environment with a single command.

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

### Quick Install (recommended)

```bash
make install
```

This will build the binary and install it to `/usr/local/bin`.

### Manual Installation

```bash
# Build the binary
go build -o ko

# Move to a directory in your PATH
sudo mv ko /usr/local/bin/ko
sudo chmod +x /usr/local/bin/ko
```

### Building from source

```bash
# Clone the repository
git clone https://github.com/bshakr/ko.git
cd ko

# Build
make build

# Or use go directly
go build -o ko
```

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
ko list                     # List all ko worktrees
ko init                     # Interactive configuration setup
ko config                   # View current configuration
ko help                     # Show help message
```

### Interactive Configuration

Run `ko init` to configure:
- Default editor (vim, nvim, code, etc.)
- Setup and dev script paths
- Custom commands for each tmux pane

Configuration is saved to `~/.config/ko/config.json`.

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

## Why Go?

This tool is written in Go for several reasons:

**Pros:**
- Single binary distribution - no runtime dependencies
- Excellent CLI libraries (Cobra + Bubble Tea)
- Cross-platform support
- Strong error handling
- Easy to maintain and extend
- Interactive TUI capabilities for configuration
- Fast execution and startup time

**Libraries used:**
- **Cobra**: Command-line interface structure and flag parsing
- **Bubble Tea**: Interactive terminal UI for the `ko init` command
- **Bubbles**: Pre-built Bubble Tea components

The combination of Cobra for CLI structure and Bubble Tea for interactive prompts provides a professional, user-friendly experience while maintaining simplicity and performance.

## Contributing

Feel free to submit issues or pull requests!

## License

MIT
