# Homebrew Setup - Implementation Complete

This document summarizes the Homebrew installation support that has been added to ko.

## What Was Done

### 1. Created Homebrew Formula
- **File**: `Formula/ko.rb`
- Formula for installing ko via Homebrew
- Includes metadata, installation instructions, and basic test
- Syntax validated successfully

### 2. Added Version Support
- **Updated**: `ko.sh`
- Added `VERSION` variable (currently 0.1.0)
- Added `ko version` command
- Updated help text to display version

### 3. Created Documentation
- **RELEASING.md**: Step-by-step release process guide
- **HOMEBREW_TAP_SETUP.md**: Instructions for creating the tap repository
- **test-formula.sh**: Testing script for the formula
- **.gitignore**: Excludes .ko directory and common temporary files

### 4. Updated README
- Added Homebrew installation as the recommended method
- Kept manual installation options
- Added `ko version` to commands list

## Next Steps to Complete Setup

### Step 1: Commit and Push Current Changes

```bash
git add .
git commit -m "Add Homebrew installation support"
git push origin add-homebrew
```

### Step 2: Merge to Main

Create a PR and merge the `add-homebrew` branch to main (or merge directly if you prefer):

```bash
git checkout main
git merge add-homebrew
git push origin main
```

### Step 3: Create v0.1.0 Release

```bash
# Create and push tag
git tag -a v0.1.0 -m "Initial release with Homebrew support"
git push origin v0.1.0

# Create GitHub release
# Go to: https://github.com/bshakr/ko/releases/new
# - Select tag: v0.1.0
# - Title: v0.1.0
# - Description: Initial release of ko with Homebrew installation support
# - Click "Publish release"
```

### Step 4: Calculate SHA256 for Formula

```bash
# Download the release tarball
curl -L https://github.com/bshakr/ko/archive/refs/tags/v0.1.0.tar.gz -o ko-0.1.0.tar.gz

# Calculate SHA256
shasum -a 256 ko-0.1.0.tar.gz

# Copy the hash and update Formula/ko.rb line 6:
# sha256 "paste_hash_here"
```

### Step 5: Create Homebrew Tap Repository

```bash
# Create a new repository at: https://github.com/new
# Repository name: homebrew-ko
# Make it public

# Clone and set up
git clone https://github.com/bshakr/homebrew-ko.git
cd homebrew-ko

# Create structure
mkdir Formula

# Copy the updated formula (with SHA256)
cp ../ko/Formula/ko.rb Formula/

# Create README
cat > README.md << 'EOF'
# Homebrew Tap for ko

Official Homebrew tap for [ko](https://github.com/bshakr/ko).

## Installation

```bash
brew install bshakr/ko/ko
```

## About

ko is a Git worktree + tmux automation tool for creating isolated development environments.

For more information, visit the [main repository](https://github.com/bshakr/ko).
EOF

# Commit and push
git add .
git commit -m "Initial tap setup with ko formula"
git push origin main
```

### Step 6: Test Installation

```bash
# Install from your tap
brew install bshakr/ko/ko

# Verify it works
ko --version
# Should output: ko version 0.1.0

# Test the help
ko help

# Uninstall (optional)
brew uninstall ko
```

## File Structure

```
ko/
├── ko.sh                          # Main script (updated with version)
├── Formula/
│   └── ko.rb                      # Homebrew formula
├── README.md                      # Updated with Homebrew installation
├── RELEASING.md                   # Release process documentation
├── HOMEBREW_TAP_SETUP.md          # Tap repository setup guide
├── HOMEBREW_SETUP_COMPLETE.md     # This file
├── test-formula.sh                # Formula testing script
├── .gitignore                     # Git ignore file
├── bin/
│   ├── setup                      # User setup template
│   └── dev                        # User dev template
└── .ko/                           # Worktree directory (gitignored)
```

## Quick Reference

### For End Users

```bash
# Install
brew install bshakr/ko/ko

# Update
brew upgrade ko

# Uninstall
brew uninstall ko
```

### For Maintainers

When releasing a new version:

1. Update `VERSION` in `ko.sh`
2. Commit: `git commit -am "Bump version to x.y.z"`
3. Tag: `git tag -a vx.y.z -m "Release version x.y.z"`
4. Push: `git push origin vx.y.z`
5. Create GitHub release
6. Calculate new SHA256
7. Update `Formula/ko.rb` in homebrew-ko repository
8. Commit and push to homebrew-ko

## Testing

Run the test script:

```bash
./test-formula.sh
```

This will audit the formula and provide testing instructions.

## Notes

- The formula is currently configured for v0.1.0
- SHA256 hash must be added after creating the GitHub release
- The tap repository must be created separately (not done automatically)
- Users will install via: `brew install bshakr/ko/ko`

## Support

For issues with:
- **ko functionality**: https://github.com/bshakr/ko/issues
- **Homebrew installation**: https://github.com/bshakr/homebrew-ko/issues (after creation)
