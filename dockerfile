
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -o /out/hris-api ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -o /out/db-migrate ./cmd/migrate

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app

COPY --from=builder --chown=app:app /out/hris-api ./hris-api
COPY --from=builder --chown=app:app /out/db-migrate ./db-migrate

USER app

EXPOSE 7100

CMD ["/bin/sh", "-c", "./db-migrate && exec ./hris-api"]
