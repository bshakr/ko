# ko - Git Worktree tmux Automation

`ko` is a CLI tool written in Go that creates git worktrees and sets up a configurable development environment with a single command.

## What it does

### `ko new <worktree-name>`

Creates a new development environment:

1. Creates a new git worktree in `.ko/<worktree-name>`
2. Opens a new tmux window with dynamically configured panes:
   - **First pane**: Runs your setup script (e.g., `./bin/setup`)
   - **Additional panes**: Runs any commands you configure (e.g., dev server, editor, etc.)

### `ko cleanup <worktree-name>`

Cleans up after you're done:

1. Closes the tmux window for the worktree
2. Removes the git worktree from `.ko/<worktree-name>`

## Prerequisites

- Git repository with worktree support
- tmux installed and running
- A setup script in your repository (optional, configurable via `ko init`)

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

### First time setup

Navigate to your git repository and run:

```bash
ko init
```

This will guide you through setting up your configuration (setup script path and pane commands).

### Creating a new worktree session

```bash
ko new <worktree-name>
```

Example:
```bash
ko new feature-auth
```

This will:
- Create a worktree at `.ko/feature-auth`
- Set up your configured tmux environment with panes running your specified commands

### Normal development workflow

Once your session is set up:
1. Work in your configured environment (editor, dev server, etc.)
2. Make commits as normal
3. Push your branch and create a PR when ready

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

## How it works

`ko` creates a new git worktree in the `.ko/` directory and opens a tmux window with panes configured based on your `.koconfig` file. The first pane runs your setup script, and additional panes run any commands you've configured (dev server, editor, etc.).

The cleanup command finds the tmux window by name and closes it, then removes the git worktree. If you have uncommitted changes, git will warn you and you'll need to either commit them or use `git worktree remove --force` manually.

## Worktree Management

All worktrees are created in a `.ko/` directory at the root of your repository. This keeps your repository organized and makes it easy to:

- See all active worktrees: `ls .ko/`
- Clean up a worktree: `ko cleanup <name>`
- Manually remove worktrees: `git worktree remove .ko/<name>`
- List all worktrees: `git worktree list`

You may want to add `.ko/` to your `.gitignore` file.

**Tip:** Use `ko cleanup` instead of manually removing worktrees - it will close the tmux window and clean up the worktree in one command!

## Configuration

Before creating your first worktree, run `ko init` to set up your configuration. The tool will prompt you for:
- Path to your setup script (if you have one)
- Additional commands to run in tmux panes

The configuration is stored in `.koconfig` at your repository root and can be updated anytime with `ko init`.


## Contributing

Feel free to submit issues or pull requests!

## License

MIT
