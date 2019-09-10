workflow "automerge pull requests on updates" {
  on = "pull_request"
  resolves = ["automerge"]
}

workflow "automerge pull requests on reviews" {
  on = "pull_request_review"
  resolves = ["automerge"]
}

workflow "automerge pull requests on status updates" {
  on = "status"
  resolves = ["automerge"]
}

workflow "automerge pull requests on check_run updates" {
  on = "check_run"
  resolves = ["automerge"]
}

action "automerge" {
  uses = "pascalgn/automerge-action@33f73f0a562129c7e96123e98af20d4378f1fa3b"
  secrets = ["GITHUB_TOKEN"]
  env = {
    LABELS = "!wip,!work in progress"
    AUTOMERGE = "ready-to-merge"
    MERGE_METHOD = "merge"
  }
}
