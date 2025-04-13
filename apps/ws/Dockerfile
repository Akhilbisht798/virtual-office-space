FROM golang:1.23.1-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o server cmd/main/main.go
EXPOSE 8080
CMD ["/app/server"]