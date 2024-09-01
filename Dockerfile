FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./

RUN go build -o oauth-mock-server

CMD ["./oauth-mock-server"]