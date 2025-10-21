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

main() {
    echo ""
    echo "═══════════════════════════════════════"
    echo "  Update homebrew-ko Tap"
    echo "═══════════════════════════════════════"
    echo ""

    # Check if homebrew-ko directory exists
    if [ ! -d "../homebrew-ko" ]; then
        print_error "homebrew-ko repository not found at ../homebrew-ko"
        echo ""
        echo "Please clone the tap repository first:"
        echo "  cd .."
        echo "  git clone https://github.com/bshakr/homebrew-ko.git"
        echo ""
        echo "Or specify the correct path by editing this script."
        exit 1
    fi

    # Get current version from Formula/ko.rb
    VERSION=$(grep 'url ".*v' Formula/ko.rb | sed -n 's/.*v\([0-9.]*\)\.tar\.gz.*/\1/p')

    if [ -z "$VERSION" ]; then
        print_error "Could not determine version from Formula/ko.rb"
        exit 1
    fi

    print_info "Detected version: $VERSION"
    echo ""

    # Copy formula to tap repository
    print_info "Copying Formula/ko.rb to homebrew-ko..."
    cp Formula/ko.rb ../homebrew-ko/Formula/ko.rb
    print_success "Formula copied"

    # Navigate to tap repository
    cd ../homebrew-ko

    # Check if there are changes
    if [[ -z $(git status -s) ]]; then
        print_info "No changes to commit (formula is already up to date)"
        exit 0
    fi

    # Show the diff
    print_info "Changes to be committed:"
    echo ""
    git diff Formula/ko.rb

    echo ""
    read -p "Commit and push these changes? (y/n) " -n 1 -r
    echo ""

    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_error "Update cancelled"
        exit 1
    fi

    # Commit and push
    print_info "Committing changes..."
    git add Formula/ko.rb
    git commit -m "Update ko to version $VERSION"

    print_info "Pushing to GitHub..."
    git push origin main

    print_success "Tap updated successfully!"
    echo ""
    echo "Users can now install with:"
    echo "  brew upgrade ko"
    echo "  # or"
    echo "  brew install bshakr/ko/ko"
    echo ""
}

main "$@"
