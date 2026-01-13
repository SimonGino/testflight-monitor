# Build frontend
FROM node:20-alpine AS frontend
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ .
RUN npm run build

# Build backend
FROM golang:1.21-alpine AS backend
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o tf-monitor ./cmd/server/main.go

# Final image
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend /app/tf-monitor .
COPY --from=backend /app/web/dist ./web/dist

EXPOSE 8080
ENV TZ=Asia/Shanghai

VOLUME ["/app/data"]

CMD ["./tf-monitor"]
