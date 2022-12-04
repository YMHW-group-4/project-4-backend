# syntax=docker/dockerfile:1
FROM golang:1.19 AS builder
WORKDIR /crypto
COPY . .
RUN go mod download && go mod verify
RUN make build_amd64

FROM amd64/alpine:3.14
WORKDIR /node
COPY --from=builder /crypto/build/crypto-linux-amd64 ./crypto
ENTRYPOINT ["./crypto"]
