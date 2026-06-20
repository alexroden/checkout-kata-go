# syntax=docker/dockerfile:1.4

FROM golang:1.26.4-alpine

WORKDIR /api

RUN apk add --no-cache git openssh

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN mkdir -p /root/.ssh && ssh-keyscan bitbucket.org >> /root/.ssh/known_hosts

RUN --mount=type=ssh go mod download

COPY . .

WORKDIR /api

RUN go build -o main .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
