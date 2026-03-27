FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/antblog-server ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /out/antblog-server /app/antblog-server
COPY config /app/config
COPY docs/sqls /app/docs/sqls
RUN mkdir -p /app/uploads
EXPOSE 8080
ENTRYPOINT ["/app/antblog-server"]
