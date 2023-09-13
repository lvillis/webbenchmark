ARG GOLANG_VERSION="1.21.0"
FROM golang:${GOLANG_VERSION}-alpine as builder

COPY . /src
WORKDIR /src

RUN <<EOF
    go mod download
    CGO_ENABLED=0 GOOS=linux go build -o ./output/webbenchmark ./cmd/main.go
EOF


FROM alpine:3.18.3 as runtime

COPY --from=builder /src/output/webbenchmark /webbenchmark

ENV URL http://cachefly.cachefly.net/100mb.test

RUN chmod +x /webbenchmark

ENTRYPOINT ["/webbenchmark"]