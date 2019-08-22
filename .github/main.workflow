workflow "on pull request pass, merge the branch" {
  resolves = ["Auto-merge pull requests"]
  on       = "check_run"
}

action "Auto-merge pull requests" {
  uses    = "./auto_merge_pull_requests"
  secrets = ["GITHUB_TOKEN"]
}
