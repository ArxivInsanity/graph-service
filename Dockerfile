#Build stage
ARG GO_VERSION=1.20

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
COPY src ./src
RUN go build -o ./app ./main.go

#Run stage
FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api
COPY --from=builder /api/app .
COPY resources ./resources
ENV PORT 8081
ENV GIN_MODE release
EXPOSE 8081

ENTRYPOINT ["./app"]