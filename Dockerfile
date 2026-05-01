FROM golang:1.22

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o triple-s ./cmd/main.go

CMD ["./triple-s", "-port", "9000", "-dir", "./data"]