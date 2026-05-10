FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

RUN GOSUMDB=off CGO_ENABLED=0 go mod tidy
RUN CGO_ENABLED=0 go build -o main .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config.*.env ./

EXPOSE 8080

CMD ["./main"]