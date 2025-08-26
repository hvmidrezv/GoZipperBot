
FROM golang:1.23-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /zipper-bot ./main.go


FROM alpine:latest

WORKDIR /

COPY --from=builder /zipper-bot /zipper-bot


ENV TELEGRAM_BOT_TOKEN=""

CMD ["/zipper-bot"]