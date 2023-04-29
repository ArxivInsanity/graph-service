#Build stage
ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /graph-service
WORKDIR /graph-service

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./app ./src/main.go

#Run stage
FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /graph-service
WORKDIR /graph-service
COPY --from=builder /graph-service/graph-service .
ENV PORT 8080
ENV GIN_MODE release
EXPOSE 8080

ENTRYPOINT ["./graph-service"]
