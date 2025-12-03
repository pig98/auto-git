# typed: false
# frozen_string_literal: true

class AutoGit < Formula
  desc "Auto sync git repos on file changes with quiet period"
  homepage "https://github.com/pig98/auto-git"
  url "https://github.com/pig98/auto-git/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "MIT"
  version "0.1.0"

  depends_on "go" => :build

  def install
    system "go", "build",
           *std_go_args(output: bin/"auto-git", ldflags: "-s -w"),
           "./"
  end

  service do
    run [opt_bin/"auto-git"]
    keep_alive true
    log_path var/"log/auto-git.log"
    error_log_path var/"log/auto-git.err.log"
    working_dir HOMEBREW_PREFIX/"var/auto-git"

    environment_variables(
      "GIT_DIRS" => "",
      "QUIET_PERIOD_MINUTES" => "10",
      "LOG_LEVEL" => "INFO",
      "DISABLE_NOTIFICATIONS" => "0"
    )
  end

  test do
    # Test that the binary exists and is executable
    assert_predicate bin/"auto-git", :exist?
  end
end
