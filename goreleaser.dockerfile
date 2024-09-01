# This dockerfile is only used to build the backend image for the application using goreleaser.
FROM alpine:3.12
COPY oauth-mock-server /app/oauth-mock-server
ENTRYPOINT ["/app/oauth-mock-server"]