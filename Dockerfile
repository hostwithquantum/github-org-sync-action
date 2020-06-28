FROM alpine:3.10
COPY github-org-sync /github-org-sync
ENTRYPOINT ["/github-org-sync"]
