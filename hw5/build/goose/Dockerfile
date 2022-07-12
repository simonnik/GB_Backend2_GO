FROM golang:1.18 as builder

RUN useradd -u 10001 goose
RUN mkdir -p /src

RUN go install github.com/pressly/goose/cmd/goose@latest
RUN cp $(which goose) /src

FROM busybox

COPY --from=builder /etc/passwd /etc/passwd

USER goose

WORKDIR /app

COPY --from=builder /src/goose /app/goose

CMD ["/app/goose"]