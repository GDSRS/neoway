FROM golang:1.22.0-bullseye as builder
WORKDIR /src/app
COPY . .
RUN go mod download && CGO_ENABLED=0 go build -o main main.go
CMD ./main
