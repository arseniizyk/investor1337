FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /app/main ./cmd/main.go

FROM alpine
COPY --from=builder /app/main /main
CMD ["/main"]