FROM golang:1.22.0-alpine

WORKDIR /app/

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bin/app ./cmd/app/main.go

EXPOSE 8080

CMD ["/app/bin/app"]