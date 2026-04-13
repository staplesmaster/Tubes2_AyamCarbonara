# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS base
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY go.mod ./
RUN go mod download

FROM base AS dev
COPY . .
EXPOSE 8080
CMD ["go", "run", "./src"]

FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./src

FROM alpine:3.20 AS prod
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /out/server /app/server
EXPOSE 8080
CMD ["/app/server"]