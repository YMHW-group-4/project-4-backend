# syntax=docker/dockerfile:1
FROM golang:1.19 AS builder
WORKDIR /builder
COPY . .
RUN go mod download && go mod verify
RUN make build_arm64v8

FROM arm64v8/alpine:3.14
WORKDIR /app
COPY --from=builder /builder/build/crypto-linux-arm64v8 ./crypto
ENTRYPOINT ["./crypto"]
