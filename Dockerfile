# 1. Builder
FROM golang:alpine AS builder

RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates && \
    apk add --update-cache tzdata && \
    update-ca-certificates

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod tidy

COPY . .

RUN go build -v -o main .

ENV GOOS=js \
    GOARCH=wasm

RUN go build -v -o web/app.wasm .

WORKDIR /dist

RUN cp /build/main .
RUN cp /build/setting.ini .
RUN cp -r /build/web ./web

### 2. Make executable image
FROM scratch

COPY --from=builder /dist/ .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Asia/Seoul \
    ZONEINFO=/zoneinfo.zip

ENTRYPOINT ["/main"]
