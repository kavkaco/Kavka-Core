FROM golang:alpine as builder

RUN apk add make

WORKDIR /app

WORKDIR /build

RUN go mod tidy

RUN make build

FROM golang:alpine as runtime

COPY --from=builder /app/build/server /usr/bin/server

EXPOSE 8000

CMD ["/usr/bin/server"]
