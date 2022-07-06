FROM golang:1.17 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

FROM golang:1.17 as builder

COPY --from=modules /go/pkg /go/pkg

RUN mkdir -p /src

ADD . /src

WORKDIR /src

RUN useradd -u 10001 shortener

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/api ./cmd/api

FROM busybox

ENV APP_HOST http://localhost:8083
ENV APP_DB_HOST host.docker.internal
ENV APP_DB_NAME shortener
ENV APP_DB_PASSWORD shortener
ENV APP_DB_PORT 54322
ENV APP_DB_USER shortener
ENV APP_DB_SSL_MODE disable
ENV APP_DB_MIGRATIONS_DIR /app/migrations
ENV APP_LOG_LEVEL debug
ENV APP_PORT 8083
ENV APP_JWT_SECRET shazam
ENV APP_HASH_SALT hahash
ENV APP_HASH_MIN_LENGTH 7

COPY --from=builder /etc/passwd /etc/passwd

USER shortener

WORKDIR /app

COPY --from=builder /app/api /app/api
COPY --from=builder /src/web /app/web
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

CMD ["/app/api"]