FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .

RUN <<EOF
go mod tidy 
go build
EOF

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
WORKDIR /app
COPY --from=builder /app/mcp-dungeon .
ENTRYPOINT ["./mcp-dungeon"]
