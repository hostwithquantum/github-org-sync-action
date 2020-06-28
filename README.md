![ci](https://github.com/hostwithquantum/github-org-sync-action/workflows/ci/badge.svg) ![release](https://github.com/hostwithquantum/github-org-sync-action/workflows/release/badge.svg)

# github-org-sync(-action)

`github-org-sync(-action)` is a tool to sync files from a skeleton/template repository to other repositories in your organization. It works within a single organization, but of course it can be used multiple times.

It currently syncs `.github/workflows` — more is planned.

Inspiration from:
 - https://github.com/cloudalchemy/auto-maintenance
 - https://github.com/prometheus/prometheus/blob/master/scripts/sync_repo_files.sh

## GitHub Action

### Inputs

#### `github-user`

**Required** Name of the user to commit files and open PRs.

#### `github-email`

**Required** Email of the user to commit files and open PRs.

#### `github-access-token`

**Required** An access token (used for clone and push), scope repo.

#### `github-org`

**Required** Your Github organization. Default: `hostwithquantum`

#### `github-skeleton`

**Required** The base repository to sync files from. Default: `ansible-skeleton`.

#### `github-repos`

**Required** A (space-delimited) list of repositories to sync to. Default: `"ansible-weave"`

## Standalone

If GitHub Actions are not applicable to your environment, you may run the tool without (and e.g. via `cron.d`).

### Getting it

All releases are done using the amazing `goreleaser` in a release workflow in this repository. The workflow validates and builds binaries for different OS' and architectures and a Docker image (amd64) as well.

#### Installing

You may download a release for Mac, Linux or Windows here:
 - https://github.com/hostwithquantum/github-org-sync-action/releases

To use the Docker image, please follow one of these links:
 - https://quay.io/repository/hostwithquantum/github-org-sync?tab=tags
 - https://github.com/hostwithquantum/github-org-sync-action/packages

Otherwise: `make dev` to build a snapshot.

#### Configuration

 - `github-org-sync` uses environment variables. There is nothing else currently.
 - See [.envrc-dist](.envrc-dist) for necessary configuration.
 - Please note that, `GITHUB_ACCESS_TOKEN` **requires** full repo scope, in order to create branches, push them and open pull-requests.

# Author

[Planetary Quantum GmbH](https://www.planetary-quantum.com) — come check us out. :rocket: :)

## License

BSD-2-Clause ("Simplified BSD License")
