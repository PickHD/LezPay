FROM golang:1.20.3 AS builder
LABEL maintainer="taufikjanuar35@gmail.com"

RUN go version

WORKDIR /merchant
COPY ./merchant/go.mod ./
COPY ./merchant/go.sum ./

RUN go mod download

COPY ./merchant .

# Build Go App
RUN CGO_ENABLED=0 GOOS=linux go build -o merchant-service ./cmd

FROM alpine:3.11.3

WORKDIR /app

RUN mkdir cmd docs

COPY --from=builder ./merchant/cmd/ ./cmd
COPY --from=builder ./merchant/docs/ ./docs
COPY --from=builder ./merchant/merchant-service .

# Command to run the executeable
ENTRYPOINT ["./merchant-service","consumer"]
