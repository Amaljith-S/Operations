FROM golang:1.17-alpine as builder
RUN apk update
RUN apk add curl && apk add python3
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN mkdir -p /home/amaljith/.kube/minikube
COPY config /home/amaljith/.kube
ADD minikube /home/amaljith/.kube/minikube
RUN go mod download
RUN go build main.go

FROM alpine:3.11.3

WORKDIR /app
RUN mkdir -p /home/amaljith/.kube/minikube
COPY --from=builder /app/config /home/amaljith/.kube
COPY --from=builder /app/minikube /home/amaljith/.kube/minikube
COPY --from=builder /app/main .