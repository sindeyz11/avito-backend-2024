FROM golang:1.22.0-alpine

WORKDIR /app/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./bin/cmd ./cmd/app/main.go

EXPOSE 8080

CMD ["./bin/cmd"]