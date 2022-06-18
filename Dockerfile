# API build stage
FROM golang:1.18.3-alpine3.14 as go-builder
ARG GOPROXY=goproxy.cn

ENV GOPROXY=https://${GOPROXY},direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache make bash git tzdata

WORKDIR /data

COPY go.mod go.sum ./
RUN go mod download -x
COPY config config
COPY Makefile Makefile
COPY scripts scripts
RUN make build


# Fianl running stage
FROM alpine:3.14.3
LABEL maintainer="goproxy@gotomicro.com"

WORKDIR /data

COPY --from=go-builder /data/bin/goproxy ./bin/
COPY --from=go-builder /data/config ./config

RUN apk add --no-cache tzdata

CMD ["sh", "-c", "./bin/goproxy"]