# Tag this image `qveen` for use in tests.
FROM docker.io/golang:1.22.5-bookworm AS build

WORKDIR /build
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN go build

FROM docker.io/debian:bookworm

COPY --from=build /build/qveen /usr/bin/qveen
ENTRYPOINT /usr/bin/qveen
