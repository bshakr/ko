#!/bin/bash
# Script to test the Homebrew formula locally before releasing

set -e

echo "=== Testing ko Homebrew Formula ==="
echo ""

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "Error: Homebrew is not installed"
    echo "Please install Homebrew first: https://brew.sh"
    exit 1
fi

echo "1. Auditing formula for issues..."
brew audit --strict Formula/ko.rb || {
    echo "Warning: Formula audit found issues (this is expected before creating a release)"
}

echo ""
echo "=== Formula Testing Instructions ==="
echo ""
echo "To fully test the formula, you need to:"
echo ""
echo "1. Create a git tag:"
echo "   git tag -a v0.1.0 -m 'Release version 0.1.0'"
echo "   git push origin v0.1.0"
echo ""
echo "2. Create a GitHub release at:"
echo "   https://github.com/bshakr/ko/releases/new"
echo ""
echo "3. Calculate the SHA256 of the release tarball:"
echo "   curl -L https://github.com/bshakr/ko/archive/refs/tags/v0.1.0.tar.gz | shasum -a 256"
echo ""
echo "4. Update Formula/ko.rb with the SHA256 hash"
echo ""
echo "5. Test the installation:"
echo "   brew install --build-from-source Formula/ko.rb"
echo "   ko --version"
echo "   brew uninstall ko"
echo ""
echo "=== Quick Local Test (without release) ==="
echo ""
echo "For a quick test without creating a release, you can modify the formula"
echo "to install from the local directory. See Formula/ko.rb for instructions."
echo ""
