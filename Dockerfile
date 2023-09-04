ARG GOLANG_VERSION="1.21.0"
FROM golang:${GOLANG_VERSION}-alpine as builder

COPY . /src/
WORKDIR /src

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go mod download
RUN go build -o ./bin/webbenchmark ./cmd/main.go

FROM alpine:3.18.3 as runtime

COPY --from=builder /src/bin/webbenchmark /bin/webbenchmark

ENV URL http://cachefly.cachefly.net/100mb.test

RUN chmod +x /bin/webbenchmark

ENTRYPOINT ["/bin/webbenchmark"]