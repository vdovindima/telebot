FROM golang:1.20.4-alpine3.18 as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /telebot main.go
FROM alpine:3.18
COPY --from=builder telebot /bin/telebot
ENTRYPOINT ["/bin/telebot"]
CMD [ "--port", "3000" ]