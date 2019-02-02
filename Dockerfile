# Using docker 2 stage build
# https://docs.docker.com/develop/develop-images/multistage-build/
FROM golang:1 AS build
# Copy entire project and build it.
COPY . /go/src/github.com/kavehmz/math/
WORKDIR /go/src/github.com/kavehmz/math/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOROOT_FINAL=/ go build -o /bin/math cmd/main.go

FROM debian:stable-slim
COPY --from=build /bin/math /bin/math

ENTRYPOINT [ "/bin/math" ]
