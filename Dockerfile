FROM golang:1.19-alpine AS builder
ENV CGO_ENABLED=0

ADD . /build
WORKDIR /build

RUN go build -o ./out/block-balance

FROM alpine:3.15

COPY --from=builder /build/out/block-balance /app/block-balance
WORKDIR /app
ENTRYPOINT ["/app/block-balance"]
