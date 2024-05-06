FROM golang:1.20.3 AS builder
LABEL maintainer="taufikjanuar35@gmail.com"

RUN go version

WORKDIR /wallet
COPY ./wallet/go.mod ./
COPY ./wallet/go.sum ./

RUN go mod download

COPY ./wallet .

# Build Go App
RUN CGO_ENABLED=0 GOOS=linux go build -o wallet-service ./cmd

FROM alpine:3.11.3

WORKDIR /app

RUN mkdir cmd docs

COPY --from=builder ./wallet/cmd/ ./cmd
COPY --from=builder ./wallet/docs/ ./docs
COPY --from=builder ./wallet/wallet-service .

# Expose port
EXPOSE 8081

# Command to run the executeable
ENTRYPOINT ["./wallet-service","http"]
