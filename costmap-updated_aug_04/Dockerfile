FROM golang:1.17-alpine as builder
RUN apk update
RUN apk add curl && apk add python3

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o ./


FROM alpine:3.11.3

WORKDIR /app
COPY --from=builder /app/costmap .
RUN mkdir /app/templates
COPY --from=builder /app/templates /app/templates

