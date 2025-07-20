FROM golang:1.24.5-alpine3.21 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o /app/tracker /app/cmd/tracker/main.go 

FROM alpine:latest
COPY --from=builder /app/tracker /tracker
ENTRYPOINT ["/tracker"]