FROM golang:1.12.0-alpine3.9
COPY . /go/src/github.com/EG-easy/sample-cosmos-app
WORKDIR /go/src/github.com/EG-easy/sample-cosmos-app
# 必要なパッケージなどをインストールする
RUN apk add --no-cache \
        alpine-sdk \
        git \
    && go get -u github.com/golang/dep/cmd/dep \
    && dep ensure -v -vendor-only=true

EXPOSE 1317 26657
