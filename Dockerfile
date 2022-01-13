FROM golang:alpine as builder

LABEL maintainer="Marcel Blijleven <marcelblijleven@gmail.com>"

RUN apk update && apk add --no-cache git
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM ubuntu:latest

WORKDIR /root/

RUN apt-get update
RUN apt-get install ca-certificates -y

COPY --from=builder /app/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /usr/local/share/ca-certificates

RUN update-ca-certificates

ENTRYPOINT ["./main"]
