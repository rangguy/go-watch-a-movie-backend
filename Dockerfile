FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./go-movies-backend ./cmd/api/main.go
 

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/go-movies-backend .
EXPOSE 8080
ENTRYPOINT ["./go-movies-backend"]