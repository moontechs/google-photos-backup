FROM cgr.dev/chainguard/go AS builder
RUN env GOBIN=/bin go install github.com/go-task/task/v3/cmd/task@latest

FROM surrealdb/surrealdb:v1.5.0
COPY --from=builder /bin/task /usr/bin/

ENTRYPOINT [ "/usr/bin/task" ]