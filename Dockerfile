FROM golang:1.22-alpine AS build

# Install dependencies
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh make build-base

WORKDIR /build

ADD . /build/code

RUN  cd /build/code && make && cp build/bin/bfbackend /bfbackend

FROM alpine

WORKDIR /root

COPY  --from=build /bfbackend /usr/bin/bfbackend

ENTRYPOINT [ "bfbackend" ]