class Koh < Formula
  desc "Git Worktree tmux Automation - Manage git worktrees with tmux sessions"
  homepage "https://github.com/bshakr/koh"
  url "https://github.com/bshakr/koh.git",
      using:    :git,
      tag:      "v0.1.1",
      revision: "21fbb1a30e0b775a02b75b2c1e73fea1c10bf3d1"
  license "MIT"
  head "https://github.com/bshakr/koh.git", branch: "main"

  depends_on "go" => :build

  def install
    # Build the binary with version injection
    ldflags = "-s -w -X github.com/bshakr/koh/cmd.Version=#{version}"
    system "go", "build", *std_go_args(ldflags: ldflags), "-o", bin/"koh"
  end

  test do
    # Test that the binary runs and shows version/help
    assert_match "Git Worktree tmux Automation", shell_output("#{bin}/koh --help")
  end
end
