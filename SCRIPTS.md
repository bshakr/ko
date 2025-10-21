# Release and Build Scripts

This directory contains automation scripts for releasing and managing ko.

## Scripts Overview

### `release.sh` - Automated Release Script

Automates the entire release process including version bumping, tagging, GitHub release creation, and formula updates.

**Usage:**
```bash
./release.sh VERSION
```

**Example:**
```bash
./release.sh 0.2.0
```

**What it does:**
1. Validates prerequisites (git status, GitHub CLI)
2. Updates VERSION in ko.sh
3. Commits version bump
4. Creates and pushes git tag
5. Creates GitHub release
6. Downloads release tarball and calculates SHA256
7. Updates Formula/ko.rb with new version and SHA256
8. Commits and pushes formula update

**Prerequisites:**
- GitHub CLI (`gh`) installed and authenticated
- Clean git working directory
- All changes committed and pushed

### `update-tap.sh` - Homebrew Tap Update Script

Updates the homebrew-ko tap repository with the latest formula.

**Usage:**
```bash
./update-tap.sh
```

**What it does:**
1. Copies Formula/ko.rb to ../homebrew-ko
2. Shows diff of changes
3. Commits and pushes to homebrew-ko repository

**Prerequisites:**
- homebrew-ko repository cloned at ../homebrew-ko

### `test-formula.sh` - Formula Testing Script

Provides instructions and checks for testing the Homebrew formula.

**Usage:**
```bash
./test-formula.sh
```

**What it does:**
- Audits the formula for issues
- Provides instructions for local testing
- Shows commands for testing installation

## Typical Release Workflow

1. **Make your changes** to ko.sh and commit them

2. **Run the release script:**
   ```bash
   ./release.sh 0.2.0
   ```

3. **Update the Homebrew tap:**
   ```bash
   ./update-tap.sh
   ```

4. **Test the installation:**
   ```bash
   brew uninstall ko  # if previously installed
   brew install bshakr/ko/ko
   ko --version
   ```

## First Time Setup

Before using these scripts for the first time:

1. **Install GitHub CLI:**
   ```bash
   brew install gh
   ```

2. **Authenticate with GitHub:**
   ```bash
   gh auth login
   ```

3. **Create homebrew-ko repository:**
   - Follow instructions in HOMEBREW_TAP_SETUP.md
   - Clone it at ../homebrew-ko

## Troubleshooting

### "GitHub CLI (gh) is not installed"
```bash
brew install gh
gh auth login
```

### "You have uncommitted changes"
```bash
git status
git add .
git commit -m "Your commit message"
```

### "homebrew-ko repository not found"
```bash
cd ..
git clone https://github.com/bshakr/homebrew-ko.git
cd ko
```

### Testing locally without release
Edit Formula/ko.rb to use a local path temporarily:
```ruby
# url "https://..."
# Replace with:
def install
  bin.install "ko.sh" => "ko"
end
```

Then test:
```bash
brew install --build-from-source Formula/ko.rb
```

## See Also

- **RELEASING.md** - Detailed release process documentation
- **HOMEBREW_TAP_SETUP.md** - Setting up the Homebrew tap
- **HOMEBREW_SETUP_COMPLETE.md** - Complete setup summary
