FROM golang:1.22 AS builder
LABEL maintainer="taufikjanuar35@gmail.com"

WORKDIR /utility
COPY ./utility/go.mod ./
COPY ./utility/go.sum ./

RUN go mod download

COPY ./utility .

# Build Go App
RUN CGO_ENABLED=0 GOOS=linux go build -o utility-service ./cmd

FROM alpine:3.11.3

WORKDIR /app

RUN mkdir cmd docs

COPY --from=builder ./utility/cmd/ ./cmd
COPY --from=builder ./utility/docs/ ./docs
COPY --from=builder ./utility/utility-service .

# Expose port
EXPOSE 8085

# Command to run the executeable
ENTRYPOINT ["./utility-service","http"]
