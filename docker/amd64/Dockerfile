# syntax=docker/dockerfile:1
FROM golang:1.19 AS builder
WORKDIR /builder
COPY . .
RUN go mod download && go mod verify
RUN make build_amd64

FROM amd64/alpine:3.14
WORKDIR /app
COPY --from=builder /builder/build/crypto-linux-amd64 ./crypto
ENTRYPOINT ["./crypto"]
