# Release Process for ko

This document outlines the steps to create a new release of ko.

## Version Numbering

We follow semantic versioning (MAJOR.MINOR.PATCH):
- MAJOR: Incompatible API changes
- MINOR: New functionality in a backwards compatible manner
- PATCH: Backwards compatible bug fixes

## Automated Release (Recommended)

Use the automated release script to handle the entire process:

```bash
./release.sh VERSION
```

Example:
```bash
./release.sh 0.2.0
```

This script will:
1. ✓ Validate prerequisites (git status, gh CLI)
2. ✓ Update version in `ko.sh`
3. ✓ Commit the version bump
4. ✓ Create and push git tag
5. ✓ Create GitHub release
6. ✓ Calculate SHA256 of release tarball
7. ✓ Update `Formula/ko.rb` with new version and SHA256
8. ✓ Commit and push formula update

### Prerequisites for Automated Release

1. **GitHub CLI (gh)** must be installed and authenticated:
   ```bash
   brew install gh
   gh auth login
   ```

2. **Clean working directory** (no uncommitted changes)

3. **Up to date with remote** (all changes pushed)

### Update Homebrew Tap

After the release script completes, update the homebrew-ko tap:

```bash
./update-tap.sh
```

This script will:
1. Copy the updated formula to `../homebrew-ko`
2. Show you the changes
3. Commit and push to the tap repository

## Manual Release Process

If you prefer to release manually, follow these steps:

### 1. Update Version Number

Edit `ko.sh` and update the `VERSION` variable:

```bash
VERSION="x.y.z"
```

### 2. Update CHANGELOG (if exists)

Document all changes since the last release.

### 3. Commit Version Bump

```bash
git add ko.sh
git commit -m "Bump version to x.y.z"
```

### 4. Create and Push Tag

```bash
git tag -a vx.y.z -m "Release version x.y.z"
git push origin main
git push origin vx.y.z
```

### 5. Create GitHub Release

```bash
gh release create vx.y.z \
  --title "vx.y.z" \
  --notes "Release version x.y.z"
```

Or manually at: https://github.com/bshakr/ko/releases/new

### 6. Update Homebrew Formula

After creating the release, update the formula with the correct SHA256:

```bash
# Download the release tarball and calculate SHA256
SHA256=$(curl -sL https://github.com/bshakr/ko/archive/refs/tags/vx.y.z.tar.gz | shasum -a 256 | awk '{print $1}')

# Update Formula/ko.rb
# - Update the url line with the new version
# - Update the sha256 line with $SHA256
```

### 7. Push Formula Update

```bash
git add Formula/ko.rb
git commit -m "Update formula for version x.y.z"
git push origin main
```

### 8. Update Homebrew Tap

```bash
cd ../homebrew-ko
cp ../ko/Formula/ko.rb Formula/
git add Formula/ko.rb
git commit -m "Update ko to version x.y.z"
git push origin main
```

## First Release Checklist

For the first v0.1.0 release:

- [ ] Ensure VERSION is set to "0.1.0" in ko.sh
- [ ] Create tag: `git tag -a v0.1.0 -m "Initial release"`
- [ ] Push tag: `git push origin v0.1.0`
- [ ] Create GitHub release from the tag
- [ ] Calculate SHA256 of the release tarball
- [ ] Update Formula/ko.rb with the SHA256
- [ ] Create homebrew-ko tap repository
- [ ] Copy Formula/ko.rb to the tap repository
- [ ] Update README.md with Homebrew installation instructions
