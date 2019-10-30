FROM golang:1.13-buster as build

RUN apt-get update && apt-get install -y --no-install-recommends \
        git \
        && rm -rf /var/lib/apt/lists/*

RUN groupadd --non-unique --gid 1001 buid-group \
    && useradd --non-unique -m --uid 1001 --gid 1001 build-user
RUN mkdir /build && chown build-user /build
USER build-user

WORKDIR /build

COPY go.* /build/
RUN go mod download

ADD . /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -tags netgo -ldflags "-X main.version=$(git describe --tags --dirty --always) -w -extldflags -static" \
        -o /build/go-direct .

FROM gcr.io/distroless/static
USER nonroot
WORKDIR /

COPY --from=build /build/go-direct /

EXPOSE 8080

ENTRYPOINT ["/go-direct"]