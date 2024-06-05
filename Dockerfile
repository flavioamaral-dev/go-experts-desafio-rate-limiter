FROM golang:1.22.3-alpine as builder
WORKDIR /app

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod tidy
RUN go mod download
COPY . .
COPY .env ./
RUN go build -o rate ./cmd/server


EXPOSE 8080

CMD ["/app/rate"]

