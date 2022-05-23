FROM golang:1.17-alpine as builder
RUN apk add --no-cache git make build-base

WORKDIR /go/src/backend

ADD go.mod ./
ADD go.sum ./
ENV GO111MODULE=on
ENV CGO_ENABLED=1

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .
RUN go build -o app .

FROM alpine:latest
ENV ZONEINFO=/var/app/zoneinfo.zip

WORKDIR /var/app/backend
COPY --from=builder /go/src/backend/app /var/app/backend/
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /var/app/

ENTRYPOINT  ["./app"]