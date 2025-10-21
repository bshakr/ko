# Setting Up Homebrew Tap for ko

This guide explains how to create and maintain the Homebrew tap for ko.

## What is a Homebrew Tap?

A Homebrew tap is a third-party repository containing Homebrew formulae. Users can install software from taps using:

```bash
brew tap username/repo
brew install username/repo/formula
```

## Creating the Tap Repository

### 1. Create GitHub Repository

Create a new repository named `homebrew-ko` at:
https://github.com/new

**Important naming convention**: The repository MUST be named `homebrew-ko` where "ko" matches your formula name.

### 2. Initialize the Repository

```bash
# Clone the new repository
git clone https://github.com/bshakr/homebrew-ko.git
cd homebrew-ko

# Create Formula directory
mkdir Formula

# Copy the formula from the main ko repository
cp ../ko/Formula/ko.rb Formula/

# Create README
cat > README.md << 'EOF'
# Homebrew Tap for ko

This is the official Homebrew tap for [ko](https://github.com/bshakr/ko).

## Installation

```bash
brew install bshakr/ko/ko
```

## About ko

ko is a Git worktree + tmux automation tool for creating isolated development environments.

For more information, visit the [main repository](https://github.com/bshakr/ko).
EOF

# Add, commit, and push
git add .
git commit -m "Initial tap setup with ko formula"
git push origin main
```

### 3. Directory Structure

Your tap repository should look like this:

```
homebrew-ko/
├── Formula/
│   └── ko.rb
└── README.md
```

## Installing from the Tap

Once the tap is set up, users can install ko with:

```bash
# Install directly (taps automatically)
brew install bshakr/ko/ko

# Or tap first, then install
brew tap bshakr/ko
brew install ko
```

## Updating the Formula

When releasing a new version of ko:

1. Update the `url` and `sha256` in `Formula/ko.rb`
2. Commit and push the changes
3. Homebrew will automatically use the updated formula

```bash
cd homebrew-ko
# Edit Formula/ko.rb with new version and SHA256
git add Formula/ko.rb
git commit -m "Update ko to version x.y.z"
git push origin main
```

## Testing the Formula

Before pushing formula changes, test locally:

```bash
# Audit the formula for issues
brew audit --strict Formula/ko.rb

# Install from local formula
brew install --build-from-source Formula/ko.rb

# Test the installation
ko --version

# Uninstall
brew uninstall ko
```

## Publishing to Homebrew Core (Optional)

Once ko gains popularity and stability, you can submit it to the official homebrew-core:

1. Ensure the formula meets [Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae) requirements
2. Fork [Homebrew/homebrew-core](https://github.com/Homebrew/homebrew-core)
3. Add your formula to `Formula/k/ko.rb`
4. Submit a pull request

Benefits of homebrew-core:
- Users can install with just `brew install ko`
- Wider distribution
- Official support

For now, a personal tap is recommended for easier maintenance.
