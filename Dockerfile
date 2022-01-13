FROM golang:alpine as builder

LABEL maintainer="Marcel Blijleven <marcelblijleven@gmail.com>"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM ubuntu:latest

WORKDIR /root/

COPY --from=builder /app/main .

ENTRYPOINT ["./main"]
