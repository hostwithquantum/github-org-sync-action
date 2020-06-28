![ci](https://github.com/hostwithquantum/github-org-sync-action/workflows/ci/badge.svg) ![release](https://github.com/hostwithquantum/github-org-sync-action/workflows/release/badge.svg)

# github-org-sync(-action)

`github-org-sync` is a little tool to sync files from a skeleton/template repository to other repositories in your organization. It works within a single organization, but of course it can be used multiple times.

Inspiration from:
 - https://github.com/cloudalchemy/auto-maintenance
 - https://github.com/prometheus/prometheus/blob/master/scripts/sync_repo_files.sh

### Status

We currently provide a tool (written in Golang) to do the syncing. We sync files from `.github/workflows` only, more files/directories are planned and will be added when time permits.

The goal is to offer a GitHub Action as well.

## Getting it

All releases are done using the amazing `goreleaser` in a release workflow in this repository. The workflow validates and builds binaries for different OS' and architectures and a Docker image (amd64) as well.

### Installing

You may download a release for Mac, Linux or Windows here:
 - https://github.com/hostwithquantum/github-org-sync-action/releases

To use the Docker image, please follow this link:
 - https://github.com/hostwithquantum/github-org-sync-action/packages

Otherwise: `make dev` to build a snapshot.

### Configuration

 - `github-org-sync` uses environment variables. There is nothing else currently.
 - See [.envrc-dist](.envrc-dist) for necessary configuration.
 - Please note that, `GITHUB_ACCESS_TOKEN` **requires** full repo scope, in order to create branches, push them and open pull-requests.

# Author

[Planetary Quantum GmbH](https://www.planetary-quantum.com) — come check us out. :rocket: :)

## License

BSD-2-Clause ("Simplified BSD License")
