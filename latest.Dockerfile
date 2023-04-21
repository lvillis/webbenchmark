ARG GOLANG_VERSION="1.18"
FROM golang:$GOLANG_VERSION-alpine as builder

WORKDIR /src/
COPY cmd/main.go go.* /src/
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/main .

FROM alpine:3.17 as runtime

COPY --from=builder /src/main /src/main
WORKDIR /src
ENV URL http://cachefly.cachefly.net/100mb.test

CMD chmod +x main && ./main