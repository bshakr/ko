#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() { echo -e "${GREEN}✓${NC} $1"; }
print_error() { echo -e "${RED}✗${NC} $1"; }
print_info() { echo -e "${YELLOW}→${NC} $1"; }

# Function to validate semver
validate_version() {
    if [[ ! $1 =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        print_error "Invalid version format. Use semver: MAJOR.MINOR.PATCH (e.g., 1.0.0)"
        exit 1
    fi
}

# Function to check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."

    # Check if we're in a git repository
    if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
        print_error "Not in a git repository"
        exit 1
    fi

    # Check for uncommitted changes
    if [[ -n $(git status -s) ]]; then
        print_error "You have uncommitted changes. Please commit or stash them first."
        git status -s
        exit 1
    fi

    # Check if gh CLI is installed
    if ! command -v gh &> /dev/null; then
        print_error "GitHub CLI (gh) is not installed"
        echo "Install it with: brew install gh"
        echo "Or visit: https://cli.github.com/"
        exit 1
    fi

    # Check if authenticated with GitHub
    if ! gh auth status &> /dev/null; then
        print_error "Not authenticated with GitHub CLI"
        echo "Run: gh auth login"
        exit 1
    fi

    print_success "All prerequisites met"
}

# Function to update version in ko.sh
update_version() {
    local version=$1
    print_info "Updating version to $version in ko.sh..."

    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s/^VERSION=\".*\"/VERSION=\"$version\"/" ko.sh
    else
        # Linux
        sed -i "s/^VERSION=\".*\"/VERSION=\"$version\"/" ko.sh
    fi

    print_success "Version updated in ko.sh"
}

# Function to commit and push version bump
commit_version() {
    local version=$1
    print_info "Committing version bump..."

    git add ko.sh
    git commit -m "Bump version to $version"

    print_success "Version bump committed"
}

# Function to create and push tag
create_tag() {
    local version=$1
    print_info "Creating git tag v$version..."

    git tag -a "v$version" -m "Release version $version"

    print_success "Tag v$version created"
}

# Function to push changes
push_changes() {
    local version=$1
    print_info "Pushing changes and tags to origin..."

    git push origin $(git branch --show-current)
    git push origin "v$version"

    print_success "Changes and tag pushed to GitHub"
}

# Function to create GitHub release
create_github_release() {
    local version=$1
    print_info "Creating GitHub release..."

    # Create release notes
    local release_notes="Release version $version

## Installation

\`\`\`bash
brew install bshakr/ko/ko
\`\`\`

## Changes

See commit history for details.

---
Generated with ko release script"

    gh release create "v$version" \
        --title "v$version" \
        --notes "$release_notes"

    print_success "GitHub release created"
}

# Function to calculate SHA256
calculate_sha256() {
    local version=$1
    print_info "Downloading release tarball and calculating SHA256..."

    # Wait a bit for GitHub to generate the tarball
    sleep 3

    local url="https://github.com/bshakr/ko/archive/refs/tags/v${version}.tar.gz"
    local sha256=$(curl -sL "$url" | shasum -a 256 | awk '{print $1}')

    print_success "SHA256 calculated: $sha256"
    echo "$sha256"
}

# Function to update formula
update_formula() {
    local version=$1
    local sha256=$2
    print_info "Updating Formula/ko.rb..."

    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s|url \".*\"|url \"https://github.com/bshakr/ko/archive/refs/tags/v${version}.tar.gz\"|" Formula/ko.rb
        sed -i '' "s/sha256 \".*\"/sha256 \"${sha256}\"/" Formula/ko.rb
    else
        # Linux
        sed -i "s|url \".*\"|url \"https://github.com/bshakr/ko/archive/refs/tags/v${version}.tar.gz\"|" Formula/ko.rb
        sed -i "s/sha256 \".*\"/sha256 \"${sha256}\"/" Formula/ko.rb
    fi

    print_success "Formula updated with version $version and SHA256"
}

# Function to commit formula update
commit_formula() {
    local version=$1
    print_info "Committing formula update..."

    git add Formula/ko.rb
    git commit -m "Update formula for version $version"
    git push origin $(git branch --show-current)

    print_success "Formula update committed and pushed"
}

# Main release function
main() {
    echo ""
    echo "═══════════════════════════════════════"
    echo "  ko Release Script"
    echo "═══════════════════════════════════════"
    echo ""

    # Check if version argument is provided
    if [ -z "$1" ]; then
        print_error "Version number required"
        echo ""
        echo "Usage: ./release.sh VERSION"
        echo "Example: ./release.sh 0.1.0"
        exit 1
    fi

    local VERSION=$1
    validate_version "$VERSION"

    echo "Preparing to release version: $VERSION"
    echo ""

    # Confirm with user
    read -p "Continue with release? (y/n) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_error "Release cancelled"
        exit 1
    fi

    echo ""

    # Run release steps
    check_prerequisites
    update_version "$VERSION"
    commit_version "$VERSION"
    create_tag "$VERSION"
    push_changes "$VERSION"
    create_github_release "$VERSION"

    # Calculate SHA256
    SHA256=$(calculate_sha256 "$VERSION")

    # Update formula
    update_formula "$VERSION" "$SHA256"
    commit_formula "$VERSION"

    echo ""
    echo "═══════════════════════════════════════"
    print_success "Release $VERSION completed successfully!"
    echo "═══════════════════════════════════════"
    echo ""

    # Next steps
    echo "Next steps:"
    echo ""
    echo "1. Update homebrew-ko tap repository:"
    echo "   cd /path/to/homebrew-ko"
    echo "   cp ../ko/Formula/ko.rb Formula/"
    echo "   git add Formula/ko.rb"
    echo "   git commit -m \"Update ko to version $VERSION\""
    echo "   git push origin main"
    echo ""
    echo "2. Test the installation:"
    echo "   brew uninstall ko (if previously installed)"
    echo "   brew install bshakr/ko/ko"
    echo "   ko --version"
    echo ""
    echo "Release URL: https://github.com/bshakr/ko/releases/tag/v$VERSION"
    echo ""
}

main "$@"
