FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/service ./cmd/sso/main.go


FROM alpine:3.19

WORKDIR /app/bin

COPY --from=builder /bin/service /app/bin/service


EXPOSE 8080

CMD ["./service"]