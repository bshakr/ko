#!/bin/bash

# Check if worktree name is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <worktree-name>"
    exit 1
fi

WORKTREE_NAME="$1"
CURRENT_DIR=$(pwd)

# Check if we're in a git repository
if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    echo "Error: Not in a git repository"
    echo "Please run this script from within a git repository"
    exit 1
fi

# Check if bin/setup exists
if [ ! -f "./bin/setup" ]; then
    echo "Error: ./bin/setup not found"
    echo "Please create a setup script at ./bin/setup"
    exit 1
fi

# Check if bin/dev exists
if [ ! -f "./bin/dev" ]; then
    echo "Error: ./bin/dev not found"
    echo "Please create a dev script at ./bin/dev"
    exit 1
fi

# Create .ko directory if it doesn't exist
if [ ! -d ".ko" ]; then
    echo "Creating .ko directory..."
    mkdir -p .ko
fi

# Create git worktree in .ko directory
echo "Creating git worktree: .ko/$WORKTREE_NAME"
git worktree add ".ko/$WORKTREE_NAME"

if [ $? -ne 0 ]; then
    echo "Failed to create worktree"
    exit 1
fi

WORKTREE_PATH="$CURRENT_DIR/.ko/$WORKTREE_NAME"

# Create new tmux window with the worktree name
tmux new-window -n "$WORKTREE_NAME" -c "$WORKTREE_PATH"

# Split window into 4 panes
# First split vertically (left and right)
tmux split-window -h -c "$WORKTREE_PATH"
# Split left pane horizontally
tmux select-pane -t 0
tmux split-window -v -c "$WORKTREE_PATH"
# Split right pane horizontally
tmux select-pane -t 2
tmux split-window -v -c "$WORKTREE_PATH"

# Create a marker file path for setup completion (unique per worktree)
SETUP_MARKER="/tmp/ko-setup-done-${WORKTREE_NAME}-$$"

# Run commands in each pane
# Pane 0 (top-left): vim
tmux select-pane -t 0
tmux send-keys -t 0 "vim" C-m

# Pane 1 (bottom-left): setup script (creates marker when done)
tmux send-keys -t 1 "./bin/setup && touch $SETUP_MARKER" C-m

# Pane 2 (top-right): wait for setup, then run dev script
tmux send-keys -t 2 "echo 'Waiting for setup to complete...'; while [ ! -f $SETUP_MARKER ]; do sleep 1; done; rm $SETUP_MARKER; ./bin/dev" C-m

# Pane 3 (bottom-right): claude
tmux send-keys -t 3 "claude" C-m

# Focus on the vim pane
tmux select-pane -t 0

echo "Worktree setup complete!"
