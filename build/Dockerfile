FROM golang:1.23.0-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /goapp cmd/gophermart/main.go
CMD ["/goapp"]
