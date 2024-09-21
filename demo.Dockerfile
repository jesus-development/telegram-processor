FROM golang:1.22

COPY . .

RUN go mod download

RUN go build -tags demo -o telegram-processor