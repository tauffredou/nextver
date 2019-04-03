FROM golang:alpine as builder
ENV GO111MODULE=on
RUN apk add -U --no-cache ca-certificates gcc git musl-dev
WORKDIR /src/nextver
COPY . /src/nextver
RUN go test ./... && \
    go build

FROM alpine:3.8
RUN apk add -U --no-cache ca-certificates
COPY --from=builder /src/nextver/nextver /usr/local/bin/nextver

ENTRYPOINT [ "/usr/local/bin/nextver" ]
