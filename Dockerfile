FROM golang:1.22 as build

WORKDIR /app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -tags api_server -o main

FROM alpine

RUN apk add --no-cache bash
COPY --from=build /app/main ./
COPY ./scripts/bash/wait-for-it.sh ./
COPY .env ./
COPY configs/default.yaml ./configs/

ENTRYPOINT ["./wait-for-it.sh", "postgres:5432", "-t", "60", "--", "/main", "server"]