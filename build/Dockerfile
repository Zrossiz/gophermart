FROM golang:1.23.0-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /goapp cmd/engine/engine.go
CMD ["/goapp"]
