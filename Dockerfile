FROM golang:latest

WORKDIR /server

COPY . .

EXPOSE 8000

RUN make build

CMD ["./build/server"]