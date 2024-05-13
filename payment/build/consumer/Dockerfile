FROM golang:1.22 AS builder
LABEL maintainer="taufikjanuar35@gmail.com"

WORKDIR /payment
COPY ./payment/go.mod ./
COPY ./payment/go.sum ./

RUN go mod download

COPY ./payment .

# Build Go App
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service ./cmd

FROM alpine:3.11.3

WORKDIR /app

RUN mkdir cmd

COPY --from=builder ./payment/cmd/ ./cmd
COPY --from=builder ./payment/payment-service .

# Command to run the executeable
ENTRYPOINT ["./payment-service","consumer"]
