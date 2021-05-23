FROM golang:1.16-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/main .

FROM alpine:latest
COPY --from=build /src/main /src/main
WORKDIR /src

ENV URL http://cachefly.cachefly.net/100mb.test

CMD chmod +x main && ./main