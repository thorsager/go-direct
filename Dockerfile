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
RUN make

FROM gcr.io/distroless/static
USER nonroot
WORKDIR /

COPY --from=build /build/bin/godirectd /
COPY --from=build /build/web /web

EXPOSE 8080

VOLUME /data
ENTRYPOINT [ "/godirectd" ]