#!/bin/bash

# Display usage information
usage() {
    cat << EOF
ko - Git Worktree tmux Automation

Usage:
  ko new <worktree-name>      Create a new worktree and tmux session
  ko cleanup <worktree-name>  Close tmux session and remove worktree
  ko help                     Show this help message

Examples:
  ko new feature-auth
  ko cleanup feature-auth
EOF
    exit 0
}

# Create a new worktree and tmux session
cmd_new() {
    local WORKTREE_NAME="$1"
    local CURRENT_DIR=$(pwd)

    # Check if worktree name is provided
    if [ -z "$WORKTREE_NAME" ]; then
        echo "Error: Worktree name required"
        echo "Usage: ko new <worktree-name>"
        exit 1
    fi

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

    # Check if worktree already exists
    if [ -d ".ko/$WORKTREE_NAME" ]; then
        echo "Error: Worktree .ko/$WORKTREE_NAME already exists"
        exit 1
    fi

    # Create git worktree in .ko directory
    echo "Creating git worktree: .ko/$WORKTREE_NAME"
    git worktree add ".ko/$WORKTREE_NAME"

    if [ $? -ne 0 ]; then
        echo "Failed to create worktree"
        exit 1
    fi

    local WORKTREE_PATH="$CURRENT_DIR/.ko/$WORKTREE_NAME"

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
    local SETUP_MARKER="/tmp/ko-setup-done-${WORKTREE_NAME}-$$"

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
}

# Cleanup a worktree and its tmux session
cmd_cleanup() {
    local WORKTREE_NAME="$1"
    local CURRENT_DIR=$(pwd)

    # Check if worktree name is provided
    if [ -z "$WORKTREE_NAME" ]; then
        echo "Error: Worktree name required"
        echo "Usage: ko cleanup <worktree-name>"
        exit 1
    fi

    # Check if we're in a git repository
    if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
        echo "Error: Not in a git repository"
        echo "Please run this script from within a git repository"
        exit 1
    fi

    # Check if worktree exists
    if [ ! -d ".ko/$WORKTREE_NAME" ]; then
        echo "Warning: Worktree .ko/$WORKTREE_NAME not found"
        echo "Will attempt to clean up tmux window only"
    else
        # Remove the git worktree
        echo "Removing git worktree: .ko/$WORKTREE_NAME"
        git worktree remove ".ko/$WORKTREE_NAME"

        if [ $? -ne 0 ]; then
            echo "Warning: Failed to remove worktree automatically"
            echo "You may need to run: git worktree remove .ko/$WORKTREE_NAME --force"
            echo "Or manually delete uncommitted changes first"
        else
            echo "Worktree removed successfully"
        fi
    fi

    # Find and kill the tmux window
    if command -v tmux &> /dev/null; then
        # Check if we're in a tmux session
        if [ -n "$TMUX" ]; then
            # Find the window with the worktree name
            local WINDOW_INDEX=$(tmux list-windows -F "#{window_index}:#{window_name}" | grep ":${WORKTREE_NAME}$" | cut -d: -f1)

            if [ -n "$WINDOW_INDEX" ]; then
                echo "Closing tmux window: $WORKTREE_NAME"
                tmux kill-window -t "$WINDOW_INDEX"
                echo "Tmux window closed"
            else
                echo "No tmux window found with name: $WORKTREE_NAME"
            fi
        else
            echo "Not in a tmux session, skipping tmux cleanup"
        fi
    fi

    echo "Cleanup complete!"
}

# Main command dispatcher
main() {
    local COMMAND="$1"
    shift

    case "$COMMAND" in
        new)
            cmd_new "$@"
            ;;
        cleanup)
            cmd_cleanup "$@"
            ;;
        help|--help|-h|"")
            usage
            ;;
        *)
            echo "Error: Unknown command '$COMMAND'"
            echo ""
            usage
            ;;
    esac
}

main "$@"
