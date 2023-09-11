FROM golang:1.20.3 AS builder
LABEL maintainer="taufikjanuar35@gmail.com"

RUN go version

WORKDIR /customer
COPY ./customer/go.mod ./
COPY ./customer/go.sum ./

RUN go mod download

COPY ./customer .

# Build Go App
RUN CGO_ENABLED=0 GOOS=linux go build -o customer-service ./cmd/v1

FROM alpine:3.11.3

WORKDIR /app

RUN mkdir cmd docs

COPY --from=builder ./customer/cmd/ ./cmd
COPY --from=builder ./customer/docs/ ./docs
COPY --from=builder ./customer/customer-service .

# Expose port
EXPOSE 8091

# Command to run the executeable
ENTRYPOINT ["./customer-service","grpc"]
