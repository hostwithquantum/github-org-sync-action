FROM alpine:3.12
COPY github-org-sync /github-org-sync
ENTRYPOINT ["/github-org-sync"]
