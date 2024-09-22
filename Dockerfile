FROM golang:alpine as builder

RUN apk add make

WORKDIR /build

COPY . /build

RUN go mod tidy

RUN make build

FROM golang:alpine as runtime

WORKDIR /server

COPY --from=builder ./build/server /server/server

EXPOSE 8000

CMD ["./build/server"]
