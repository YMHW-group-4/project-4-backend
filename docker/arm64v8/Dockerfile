# syntax=docker/dockerfile:1
FROM golang:1.19 AS builder

WORKDIR /crypto
COPY . .

RUN go mod download && go mod verify
RUN make build_arm64v8

FROM arm64v8/alpine:3.14

WORKDIR /node
COPY --from=builder /crypto/build/crypto-linux-arm64v8 ./crypto
ENTRYPOINT ["./crypto"]
