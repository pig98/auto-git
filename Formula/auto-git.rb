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

  # Git commit ID (will be updated by GitHub Actions workflow)
  GIT_COMMIT = "PLACEHOLDER_COMMIT"

  def install
    # Use the commit ID from the constant (set by workflow)
    # Fallback to version tag if not set
    git_commit = GIT_COMMIT
    if git_commit == "PLACEHOLDER_COMMIT" || git_commit.empty?
      # Fallback: try to extract from URL or use version
      if url.to_s.include?("/v")
        tag_match = url.to_s.match(%r{/v([^/]+)\.tar\.gz})
        git_commit = "v#{tag_match[1]}" if tag_match
      end
      git_commit = "v#{version}" if git_commit.empty? || git_commit == "PLACEHOLDER_COMMIT"
    end

    # Build time
    build_time = Time.now.strftime("%Y-%m-%dT%H:%M:%S")

    # Build with version information
    ldflags = [
      "-s", "-w",
      "-X main.version=#{version}",
      "-X main.buildTime=#{build_time}",
      "-X main.gitCommit=#{git_commit}"
    ].join(" ")

    system "go", "build",
           *std_go_args(output: bin/"auto-git", ldflags: ldflags),
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
