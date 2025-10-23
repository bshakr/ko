class Ko < Formula
  desc "Git Worktree tmux Automation - Manage git worktrees with tmux sessions"
  homepage "https://github.com/bshakr/ko"
  url "https://github.com/bshakr/ko.git",
      using:    :git,
      tag:      "v0.1.0",
      revision: "REPLACE_WITH_COMMIT_SHA"
  license "MIT"
  head "https://github.com/bshakr/ko.git", branch: "main"

  depends_on "go" => :build

  def install
    # Build the binary
    system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"ko"
  end

  test do
    # Test that the binary runs and shows version/help
    assert_match "Git Worktree tmux Automation", shell_output("#{bin}/ko --help")
  end
end
