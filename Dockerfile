# ---------- BUILD ----------

FROM golang:1.26.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \

    go build -o app ./cmd/app

# ---------- RUNTIME ----------

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/app .

CMD ["./app"]